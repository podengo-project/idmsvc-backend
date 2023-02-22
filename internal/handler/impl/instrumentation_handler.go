package impl

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *application) GetMetrics(ctx echo.Context) error {
	return echo.WrapHandler(promhttp.HandlerFor(
		app.metrics.Registry(),
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
			// Pass custom registry
			Registry: app.metrics.Registry(),
		},
	))(ctx)
}
