package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// NameSpace is the namespace for the Prometheus metrics
	NameSpace = "idmsvc"
)

// Metrics holds all the Prometheus metrics for the application
type Metrics struct {
	// HTTPRequestDuration is a histogram that measures the duration of HTTP requests
	HTTPRequestDuration *prometheus.HistogramVec

	reg *prometheus.Registry
}

// See: https://prometheus.io/docs/instrumenting/writing_exporters/#naming
// See: https://consoledot.pages.redhat.com/docs/dev/platform-documentation/understanding-slo.html
// See: https://prometheus.io/docs/tutorials/understanding_metric_types/#types-of-metrics

// NewMetrics creates a new Metrics instance
func NewMetrics(reg *prometheus.Registry) *Metrics {
	if reg == nil {
		panic("reg cannot be nil")
	}
	metrics := &Metrics{
		reg: reg,

		HTTPRequestDuration: promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
			Namespace: NameSpace,
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests",
			Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 20),
		}, []string{"status", "method", "path"}),
	}

	reg.MustRegister(collectors.NewBuildInfoCollector())

	return metrics
}

func (m Metrics) Registry() *prometheus.Registry {
	return m.reg
}
