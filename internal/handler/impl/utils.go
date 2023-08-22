package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

func getXRHID(ctx echo.Context) (*identity.XRHID, error) {
	domainCtx := ctx.(middleware.DomainContextInterface)
	xrhid := domainCtx.XRHID()
	if xrhid == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "'xrhid' is nil")
	} else {
		return xrhid, nil
	}
}
