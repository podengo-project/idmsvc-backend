package impl

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
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

func sendPendoTrackEvent(ctx echo.Context, client pendo.Pendo, event string) {
	c := ctx.Request().Context()
	logger := app_context.LogFromCtx(c)

	xrhid, err := getXRHID(ctx)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	track := pendo.TrackRequest{
		AccountID: xrhid.Identity.OrgID,
		Type:      "track",
		Event:     event,
		VisitorID: header.GetPrincipal(xrhid),
		Timestamp: time.Now().UTC().UnixMilli(),
	}

	err = client.SendTrackEvent(c, &track)
	if err != nil {
		logger.Warn(err.Error())
	}
}
