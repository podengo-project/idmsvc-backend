package impl

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"gorm.io/gorm"
)

const (
	pendoHostConfSuccess = "idmsvc-host-conf-success"
	pendoHostConfFailure = "idmsvc-host-conf-failure"
)

func (a *application) HostConf(
	ctx echo.Context,
	inventoryId public.HostId,
	fqdn string,
	params public.HostConfParams,
) error {
	if err := a.hostConf(ctx, inventoryId, fqdn, params); err != nil {
		sendPendoTrackEvent(ctx, a.pendo, pendoHostConfFailure)
		return err
	}
	sendPendoTrackEvent(ctx, a.pendo, pendoHostConfSuccess)
	return nil
}

func (a *application) hostConf(
	ctx echo.Context,
	inventoryId public.HostId,
	fqdn string,
	params public.HostConfParams,
) error {
	var (
		err     error
		input   public.HostConf
		domain  *model.Domain
		output  *public.HostConfResponse
		options *interactor.HostConfOptions
		hctoken public.HostToken
		tx      *gorm.DB
		xrhid   *identity.XRHID
		keys    []jwk.Key
	)
	handlerName := "HostConf"
	logger := app_context.LogFromCtx(ctx.Request().Context())
	logger = logger.With(slog.String("handler", handlerName))
	c := ctx.Request().Context()
	if xrhid, err = getXRHID(ctx); err != nil {
		logger.Error(errXRHIDIsNil)
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		logger.Error(errUnserializing)
		return err
	}
	if options, err = a.host.interactor.HostConf(xrhid, inventoryId, fqdn, &params, &input); err != nil {
		logger.Error(errInputAdapter)
		return err
	}
	logger = logger.With(
		slog.String("inventory_id", inventoryId.String()),
		slog.String("fqdn", fqdn),
	)

	if tx = a.db.Begin(); tx.Error != nil {
		logger.Error(errDBTXBegin)
		return tx.Error
	}
	defer tx.Rollback()

	c = app_context.CtxWithDB(c, tx)
	if domain, err = a.host.repository.MatchDomain(
		c,
		options,
	); err != nil {
		logger.Error("failed to match domain on requesting host-conf")
		return err
	}

	if keys, err = a.hostconfjwk.repository.GetPrivateSigningKeys(c); err != nil {
		logger.Error("failed to read private signing keys")
		return err
	}
	if len(keys) == 0 {
		logger.Error("failed because no keys available")
		err = echo.NewHTTPError(http.StatusInternalServerError, "no keys available")
		return err
	}

	if hctoken, err = a.host.repository.SignHostConfToken(
		c,
		keys,
		options,
		domain,
	); err != nil {
		logger.Error("failed to sign host-conf token")
		return err
	}

	if tx.Commit(); tx.Error != nil {
		logger.Error(errDBTXCommit)
		return tx.Error
	}

	if output, err = a.host.presenter.HostConf(
		domain, hctoken,
	); err != nil {
		logger.Error(errOutputAdapter)
		return err
	}
	return ctx.JSON(http.StatusOK, *output)
}
