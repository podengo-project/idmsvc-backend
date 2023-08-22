package router

import (
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/metrics"
)

func newGroupMetrics(e *echo.Group, c RouterConfig) *echo.Group {
	metrics.RegisterHandlers(e, c.Handlers)
	return e
}
