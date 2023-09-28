package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
)

func (a *application) GetSigningKeys(ctx echo.Context, params public.GetSigningKeysParams) error {
	// TODO: Not Implemented
	output := public.SigningKeysResponse{
		Keys: []string{},
	}
	return ctx.JSON(http.StatusOK, output)
}
