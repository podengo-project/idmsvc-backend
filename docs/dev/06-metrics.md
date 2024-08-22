# About metrics

- The application metrics are registered in a instance
  stored into the application handler.
- The metrics are defined at `internal/metrics/metrics.go` file.
- The router for the metrics is found at
  `internal/infrastructure/router/metrics.go` file.
- The middleware that collect the status metric is defined
  at `internal/metrics/middleware/metrics.go` file.

## Prometheus locally

- You can try your metrics with a local prometheus
  instance by running `make prometheus-up`.
- You can stop the prometheus instance by `make prometheus-stop` or clean it up by
  `make prometheus-clean`.
- You can open the prometheus console by `make prometheus-ui`.
- The local prometheus instance is configured at `configs/prometheus.yaml` file.

```sh
# First run:
$ cp configs/prometheus.example.yaml configs/prometheus.yaml
$ make run  # the example config scrapes locally running idmsvc-backend
$ make prometheus-up
```

Into the ephemeral environment, no prometheus instance is deployed into the
reserved namespace by default, as it consume resources that does not use to be
used. But you can enable it if you need. More information here:
https://consoledot.pages.redhat.com/docs/dev/getting-started/ephemeral/prometheus.html

## Grafana locally

You can see the metrics scraped by Prometheus visualized in Grafana dashboard.
No configuration is needed.

```sh
$ make grafana-up  # to start
$ make grafana-down # to stop
```

### Updating dashboard

Dashboard can be updated by doing changes in Grafana UI and then exporting
it into JSON and saving in `scripts/grafana/idmsvc_dashboard.json`.

See: https://grafana.com/docs/grafana/latest/dashboards/share-dashboards-panels/#export-a-dashboard-as-json

## Additional information

Additional information about metrics can be found
into the HMSCONTENT ticket at:
https://issues.redhat.com/browse/HMS-599
You will find several useful references on that
ticket.
