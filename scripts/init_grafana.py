#!.venv/bin/python3

"""
This script initializes Grafana with a Prometheus datasource and a dashboard

It is meant to be used via the Makefile in this project.:

```
$ make prometheus-up
$ make grafana-up
```
"""

import json
import os
import time

from grafana_client import GrafanaApi

GF_USER = "admin"
GF_PASSWORD = "admin"
GF_URL = "http://localhost:3000"

PROMETHEUS_DATASOURCE_NAME = "prometheus"
PROMETHEUS_URL = "http://prometheus:9090"

DASHBOARD_PATH = "grafana/idmsvc_dashboard.json"


def get_client():
    grafana = GrafanaApi.from_url(GF_URL, (GF_USER, GF_PASSWORD))
    return grafana


def create_prometheus_datasource(grafana):
    print(f"Creating datasource \"{PROMETHEUS_DATASOURCE_NAME}\"")
    data_source = {
        "name": PROMETHEUS_DATASOURCE_NAME,
        "type": "prometheus",
        "url": PROMETHEUS_URL,
        "access": "proxy",
    }

    response = grafana.datasource.create_datasource(data_source)
    return response


def get_datasource(grafana, name):
    datasources = grafana.datasource.list_datasources()
    for ds in datasources:
        if ds["name"] == name:
            return ds
    return None


def _replace_datasource_uid(parent, datasource):
    if "datasource" in parent and parent["datasource"]["type"] == "prometheus":
        parent["datasource"]["uid"] = datasource["uid"]


def create_dashboard(grafana, dashboard_path, prometheus_datasource):
    with open(dashboard_path, "r") as f:
        dashboard = f.read()
    dashboard = json.loads(dashboard)

    # update datasources in dashboard to use UID with the correct datasource
    for panel in dashboard["panels"]:
        _replace_datasource_uid(panel, prometheus_datasource)

        if "targets" in panel:
            for target in panel["targets"]:
                _replace_datasource_uid(target, prometheus_datasource)

    dashboard["id"] = None  # to create a new dashboard

    response = grafana.dashboard.update_dashboard({
        "dashboard": dashboard,
        "overwrite": True,
        "message": "Initial import, created by automation.",
    })
    return response


def main():
    print("Initializing Grafana")

    grafana = get_client()
    print("Waiting for Grafana to start")
    # Wait for Grafana to start
    max_retries = 20
    for _ in range(max_retries):
        try:
            grafana.connect()
            grafana.datasource.list_datasources()
            print("Grafana is up")
            break
        except Exception as e:
            print(f"Failed to connect to Grafana: {e}")
            time.sleep(1)
    else:
        print("Failed to connect to Grafana")
        return

    # Ensure Prometheus datasource exists
    prometheus_datasource = get_datasource(grafana, PROMETHEUS_DATASOURCE_NAME)
    if not prometheus_datasource:
        response = create_prometheus_datasource(grafana)
        prometheus_datasource = response["datasource"] if response else None
        print(f"Datasource \"{PROMETHEUS_DATASOURCE_NAME}\" created:")
    else:
        print(f"Datasource \"{PROMETHEUS_DATASOURCE_NAME}\" already exists")
    print(prometheus_datasource)

    # Create dashboard with Prometheus datasource from file
    script_path = os.path.dirname(os.path.realpath(__file__))
    dashboard_path = os.path.join(script_path, DASHBOARD_PATH)
    dashboard = create_dashboard(
        grafana, dashboard_path, prometheus_datasource
    )
    print("Dashboard created:")
    print(dashboard)


if __name__ == "__main__":
    main()
