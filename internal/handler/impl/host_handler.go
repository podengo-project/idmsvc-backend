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
	c := ctx.Request().Context()
	if xrhid, err = getXRHID(ctx); err != nil {
		logger.Error(errXRHIDIsNil, slog.String("handler", handlerName))
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		logger.Error(errUnserializing, slog.String("handler", handlerName))
		return err
	}
	if options, err = a.host.interactor.HostConf(xrhid, inventoryId, fqdn, &params, &input); err != nil {
		logger.Error(errInputAdapter, slog.String("handler", handlerName))
		return err
	}

	if tx = a.db.Begin(); tx.Error != nil {
		logger.Error(errDBTXBegin, slog.String("handler", handlerName))
		return tx.Error
	}
	defer tx.Rollback()

	c = app_context.CtxWithDB(c, tx)
	if domain, err = a.host.repository.MatchDomain(
		c,
		options,
	); err != nil {
		logger.Error("error matching domain",
			slog.String("handler", handlerName),
			slog.String("detail", err.Error()),
		)
		return err
	}

	if keys, err = a.hostconfjwk.repository.GetPrivateSigningKeys(c); err != nil {
		logger.Error("error reading private signing keys", slog.String("handler", handlerName))
		return err
	}
	if len(keys) == 0 {
		logger.Error("no keys available", slog.String("handler", handlerName))
		err = echo.NewHTTPError(http.StatusInternalServerError, "no keys available")
		return err
	}

	if hctoken, err = a.host.repository.SignHostConfToken(
		c,
		keys,
		options,
		domain,
	); err != nil {
		logger.Error("error signing host-conf token", slog.String("handler", handlerName))
		return err
	}

	if tx.Commit(); tx.Error != nil {
		logger.Error(errDBTXCommit, slog.String("handler", handlerName))
		return tx.Error
	}

	if output, err = a.host.presenter.HostConf(
		domain, hctoken,
	); err != nil {
		logger.Error(errOutputAdapter, slog.String("handler", handlerName))
		return err
	}
	return ctx.JSON(http.StatusOK, *output)
}
