package impl

import (
	"net/http"

	api_private "github.com/hmsidm/internal/api/private"
	"github.com/labstack/echo/v4"
)

// Liveness kubernetes probe endpoint
// (GET /livez)
func (a application) GetLivez(ctx echo.Context) error {
	// TODO Add probes to the added services
	return ctx.JSON(http.StatusOK, api_private.HealthySuccess(api_private.Ok))
}

// Readiness kubernetes probe endpoint
// (GET /readyz)
func (a application) GetReadyz(ctx echo.Context) error {
	// TODO Add probes to the added services
	return ctx.JSON(http.StatusOK, api_private.HealthySuccess(api_private.Ok))
}
