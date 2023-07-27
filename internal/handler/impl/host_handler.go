package impl

import (
	"net/http"

	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/labstack/echo/v4"
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
		tx      *gorm.DB
		xrhid   *identity.XRHID
	)
	xrhid, err = getXRHID(ctx)
	if err != nil {
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

	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}

	if output, err = a.host.presenter.HostConf(
		domain,
	); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *output)
}
