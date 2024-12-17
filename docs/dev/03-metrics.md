# Metrics

Console.Dot uses [prometheus](https://prometheus.io/docs/introduction/overview/)
for gathering the metrics, and [grafana](https://grafana.com/docs/grafana/v9.3/introduction/)
for visualizing the metrics.

## Design

Define which metrics want to be exposed:

- Read the naming conventions at: https://prometheus.io/docs/practices/naming/
- Add a namespace for your metrics (eg `idmsvc`), this will prefix
  your metric names.
- Common metric could be the http_status_count with `status`, `method`, `path`
  labels; this can raise state information of our API.
- Other common is a histogram for the service latency.

## Concerns

- How to monitor kafka events? It will need some mechanism.
- What to measure here?
  - event_status_count with `status`, `topic`.

## Developer actions

- Create a new Service that expose the metrics endpoint.
- The configuration about port and path should be read from
  `clowder.Configuration`.
- For the service, some metrics are generated at 3scale API gateway, but it only
  provides the external vision of the service.
- To gather internal metrics when our endpoints are called, define a middleware
  that collect the metrics.
- For custom metrics, it could require a go routine that gather the necessary
  information at some intervals.
- The metrics primitives use to be enough with the defaults provided by the
  prometheus client library.

## Checking locally

- Start local infrastructure by `make compose-up`
- Run your backend service locally by `make run`
- Launch some actions on your service.
- Checkout the metrics: `curl http://localhost:9000/metrics`
- Checkout that prometheus is pulling the metrics: `make prometheus-ui`

## About the metrics endpoint

- Prometheus launch a http request to an internal endpoint exposed at some port
  eg: http://idmsvc:9000/metrics

## Client libraries

- https://prometheus.io/docs/instrumenting/clientlibs/
- For the golang client, this helper functions will be useful for unit testing: 
  https://github.com/prometheus/client_golang/blob/main/prometheus/testutil/testutil.go

## References

- https://prometheus.io/docs/introduction/overview/
- https://grafana.com/docs/grafana/v9.3/introduction/
