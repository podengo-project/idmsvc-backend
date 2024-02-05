package router

import (
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/metrics"
)

func newGroupMetrics(e *echo.Echo, c RouterConfig) *echo.Echo {
	metrics.RegisterHandlersWithBaseURL(e, c.Handlers, c.MetricsPath)
	return e
}
