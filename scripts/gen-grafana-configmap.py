#!/usr/bin/env python3
import yaml
import json
import logging

# TODO
# - Keep the header of the template as it is, to keep
#   the data section at the end.
# - Write the content as string with \n representation.

CONFIGMAP_TEMPLATE = """
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-dashboard-idmsvc
  labels:
    grafana_dashboard: "true"
  annotations:
    grafana-folder: /grafana-dashboard-definitions/Insights
data: {}
"""

def generate_configmap(sources: dict = {}):
    config_map = yaml.safe_load(CONFIGMAP_TEMPLATE)
    config_map["data"].update(sources)
    return yaml.dump(config_map, default_flow_style=False)

def load_dashboard_from_json(path: str):
    try:
        with open(path, "r") as grafana_file:
            json_data = json.load(grafana_file)
            return json.dumps(json_data, indent=2)
    except FileNotFoundError:
        logging.error("File {path} not found.")
    except json.JSONDecodeError as e:
        logging.error("Error parsing JSON file: {e}")
    except Exception as anyerr:
        logging.error("Unexpected error loading JSON file {path}: {anyerr}")
    exit(1)

if __name__ == "__main__":
    try:
        content = generate_configmap({
            "grafana.js": load_dashboard_from_json("./dashboards/grafana-dashboard-1.json")
        })
        with open("./dashboards/grafana-dashboards-idmsvc-configmap.yaml", "w") as configmap_file:
            configmap_file.write(content)
    except Exception as e:
        logging.error("error: {error}", error=e)
        exit(1)
