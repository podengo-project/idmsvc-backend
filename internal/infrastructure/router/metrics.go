package router

import (
	"github.com/hmsidm/internal/api/metrics"
	"github.com/labstack/echo/v4"
)

func newGroupMetrics(e *echo.Group, c RouterConfig) *echo.Group {
	metrics.RegisterHandlers(e, c.Handlers)
	return e
}
