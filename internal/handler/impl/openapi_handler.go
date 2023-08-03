package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
)

var openapiSpec = public.PathToRawSpec("/openapi.json")

// GetOpenapi return the openapi specification as a json content
// from the boilerplate generated.
// ctx is the echo context.
// Return nil for success execution or an error object.
func (a application) GetOpenapi(ctx echo.Context) error {
	resp, err := openapiSpec["/openapi.json"]()
	if err != nil {
		return err
	}
	return ctx.Blob(http.StatusOK, echo.MIMEApplicationJSON, resp)
}
