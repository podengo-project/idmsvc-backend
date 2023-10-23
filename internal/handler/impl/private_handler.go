package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	api_private "github.com/podengo-project/idmsvc-backend/internal/api/private"
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
