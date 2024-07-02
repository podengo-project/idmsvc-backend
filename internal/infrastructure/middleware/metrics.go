package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsConfig struct {
	Skipper echo_middleware.Skipper
	Metrics *metrics.Metrics
}

var defaultConfig MetricsConfig = MetricsConfig{
	Skipper: echo_middleware.DefaultSkipper,
	Metrics: metrics.NewMetrics(prometheus.NewRegistry()),
}

func mapStatus(status int) string {
	switch {
	case status >= 100 && status < 200:
		return "1xx"
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	case status >= 500 && status < 600:
		return "5xx"
	default:
		return ""
	}
}

func MetricsMiddlewareWithConfig(config *MetricsConfig) echo.MiddlewareFunc {
	if config == nil {
		config = &defaultConfig
	}
	if config.Skipper == nil {
		config.Skipper = echo_middleware.DefaultSkipper
	}
	if config.Metrics == nil {
		panic("config.Metrics can not be nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()
			if config.Skipper(ctx) {
				return next(ctx)
			}
			method := ctx.Request().Method
			path := MatchedRoute(ctx)
			err := next(ctx)
			status := mapStatus(ctx.Response().Status)
			defer func() {
				config.Metrics.HttpStatusHistogram.WithLabelValues(status, method, path).Observe(time.Since(start).Seconds())
			}()
			return err
		}
	}
}

func CreateMetricsMiddleware(metrics *metrics.Metrics) echo.MiddlewareFunc {
	return MetricsMiddlewareWithConfig(
		&MetricsConfig{
			Metrics: metrics,
		})
}
