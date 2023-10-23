package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/redhatinsights/platform-go-middlewares/identity"
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
	)
	if xrhid, err = getXRHID(ctx); err != nil {
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	if options, err = a.host.interactor.HostConf(xrhid, inventoryId, fqdn, &params, &input); err != nil {
		return err
	}

	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if domain, err = a.host.repository.MatchDomain(
		tx,
		options,
	); err != nil {
		return err
	}
	if hctoken, err = a.host.repository.SignHostConfToken(
		a.secrets.signingKeys,
		options,
		domain,
	); err != nil {
		return err
	}

	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}

	if output, err = a.host.presenter.HostConf(
		domain, hctoken,
	); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *output)
}
