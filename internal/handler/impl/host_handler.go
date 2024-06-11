package impl

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"gorm.io/gorm"
)

func (a *application) HostConf(
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
	if xrhid, err = getXRHID(ctx); err != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return err
	}
	if options, err = a.host.interactor.HostConf(xrhid, inventoryId, fqdn, &params, &input); err != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return err
	}

	if tx = a.db.Begin(); tx.Error != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return tx.Error
	}
	defer tx.Rollback()

	if domain, err = a.host.repository.MatchDomain(
		tx,
		options,
	); err != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return err
	}

	if keys, err = a.hostconfjwk.repository.GetPrivateSigningKeys(tx); err != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return err
	}
	if len(keys) == 0 {
		err = echo.NewHTTPError(http.StatusInternalServerError, "no keys available")
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return err
	}

	if hctoken, err = a.host.repository.SignHostConfToken(
		keys,
		options,
		domain,
	); err != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return err
	}

	if tx.Commit(); tx.Error != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return tx.Error
	}

	if output, err = a.host.presenter.HostConf(
		domain, hctoken,
	); err != nil {
		slog.ErrorContext(ctx.Request().Context(), err.Error())
		return err
	}
	return ctx.JSON(http.StatusOK, *output)
}
