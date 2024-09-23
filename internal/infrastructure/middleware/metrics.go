package middleware

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
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

			err := next(ctx)

			logger := context.LogFromCtx(ctx.Request().Context())

			method := ctx.Request().Method
			path := MatchedRoute(ctx)
			status := ctx.Response().Status

			// ctx.Response().Status might not be set yet for errors
			httpErr := new(echo.HTTPError)
			if errors.As(err, &httpErr) {
				status = httpErr.Code
			}
			statusStr := strconv.Itoa(status)
			headerBuf := strings.Builder{}
			headerSize := 0.0
			if errHeaderBuf := ctx.Request().Header.Write(&headerBuf); errHeaderBuf == nil {
				headerSize = float64(len(headerBuf.String()))
				config.Metrics.HTTPRequestHeaderSize.WithLabelValues(statusStr, method, path).Observe(headerSize)
			} else {
				logger.Warn("writing headers in string buffer",
					slog.String("err", errHeaderBuf.Error()),
				)
			}
			bodySize := float64(ctx.Request().ContentLength)
			config.Metrics.HTTPRequestBodySize.WithLabelValues(statusStr, method, path).Observe(bodySize)
			logger.Debug("measured request size",
				slog.String("status", statusStr),
				slog.String("method", method),
				slog.String("path", path),
				slog.Float64("header_size", headerSize),
				slog.Float64("body_size", bodySize),
			)

			config.Metrics.HTTPRequestDuration.WithLabelValues(statusStr, method, path).Observe(time.Since(start).Seconds())

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
