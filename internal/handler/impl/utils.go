package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

func getXRHID(ctx echo.Context) (*identity.XRHID, error) {
	domainCtx, ok := ctx.(middleware.DomainContextInterface)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "'ctx' is not a DomainContextInterface")
	}
	xrhid := domainCtx.XRHID()
	if xrhid == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "'xrhid' is nil")
	} else {
		return xrhid, nil
	}
}
