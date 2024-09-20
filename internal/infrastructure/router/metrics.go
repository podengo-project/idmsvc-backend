package router

import (
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/config"
)

func newGroupMetrics(e *echo.Echo, cfg *config.Config, app metrics.ServerInterface) *echo.Echo {
	metrics.RegisterHandlersWithBaseURL(e, app, cfg.Metrics.Path)
	return e
}
