package impl

import (
	"net/http"

	"github.com/hmsidm/internal/api/header"
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/infrastructure/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

//
// This file implements the interface public.

// List Domains
// (GET /domains)
func (a *application) ListDomains(
	ctx echo.Context,
	params public.ListDomainsParams,
) error {
	var (
		err    error
		data   []model.Domain
		output *public.ListDomainsResponse
		orgId  string
		offset int
		limit  int
		tx     *gorm.DB
	)
	// TODO A call to an internal validator could be here to check public.ListTodosParams
	orgId, offset, limit, err = a.domain.interactor.List(&params)
	if err != nil {
		return err
	}
	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	if data, err = a.domain.repository.FindAll(
		tx,
		orgId,
		int64(offset),
		int32(limit),
	); err != nil {
		tx.Rollback()
		return err
	}
	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}
	// TODO Read prefix from configuration
	output, err = a.domain.presenter.List(
		"/api/hmsidm/v1",
		int64(offset),
		int32(limit),
		data,
	)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *output)
}

// Return a Domain resource
// (GET /domains/{id})
func (a *application) ReadDomain(
	ctx echo.Context,
	uuid string,
	params public.ReadDomainParams,
) error {
	var (
		err    error
		data   model.Domain
		output *public.ReadDomainResponse
		orgId  string
		itemId string
		tx     *gorm.DB
	)

	if orgId, itemId, err = a.domain.interactor.GetById(uuid, &params); err != nil {
		return err
	}
	tx = a.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if data, err = a.domain.repository.FindById(tx, orgId, itemId); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return err
	}
	if output, err = a.domain.presenter.Get(&data); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *output)
}

// // Modify an existing Domain
// // (PATCH /domains/{id})
// func (a *application) PartialUpdateTodo(ctx echo.Context, id public.Id, params public.PartialUpdateTodoParams) error {
// 	var (
// 		err    error
// 		data   *model.Domain
// 		output public.Todo
// 		input  public.Todo
// 		tx     *gorm.DB
// 	)

// 	if err = ctx.Bind(&input); err != nil {
// 		return err
// 	}
// 	data = &model.Todo{}
// 	if err = a.todo.interactor.PartialUpdate(id, &params, &input, data); err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	tx = a.db.Begin()
// 	if *data, err = a.todo.repository.PartialUpdate(tx, data); err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	tx.Commit()
// 	if err = a.todo.presenter.PartialUpdate(data, &output); err != nil {
// 		return err
// 	}
// 	return ctx.JSON(http.StatusOK, output)
// }

// // Replace an existing Todo
// // (PUT /todo/{id})
// func (a *application) UpdateTodo(ctx echo.Context, id public.Id, params public.UpdateTodoParams) error {
// 	var (
// 		err    error
// 		data   *model.Todo
// 		output public.Todo
// 		input  public.Todo
// 		tx     *gorm.DB
// 	)

// 	if err = ctx.Bind(&input); err != nil {
// 		return err
// 	}
// 	data = &model.Todo{}
// 	if err = a.todo.interactor.FullUpdate(id, &params, &input, data); err != nil {
// 		return err
// 	}
// 	tx = a.db.Begin()
// 	if *data, err = a.todo.repository.Update(tx, data); err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	tx.Commit()
// 	if err = a.todo.presenter.FullUpdate(data, &output); err != nil {
// 		return err
// 	}
// 	return ctx.JSON(http.StatusOK, output)
// }

// Create a Todo resource
// (POST /todo)
func (a *application) CreateDomain(
	ctx echo.Context,
	params public.CreateDomainParams,
) error {
	var (
		err    error
		input  public.CreateDomain
		orgId  string
		data   *model.Domain
		output *public.CreateDomainResponse
		tx     *gorm.DB
	)

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	if orgId, data, err = a.domain.interactor.Create(&params, &input); err != nil {
		return err
	}

	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	if err = a.domain.repository.Create(tx, orgId, data); err != nil {
		tx.Rollback()
		return err
	}
	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}
	if output, err = a.domain.presenter.Create(data); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, *output)
}

// Delete a Domain resource
// (DELETE /domains/{uuid})
func (a *application) DeleteDomain(
	ctx echo.Context,
	uuid string,
	params public.DeleteDomainParams,
) error {
	var (
		err         error
		tx          *gorm.DB
		orgId       string
		domain_uuid string
	)
	if orgId, domain_uuid, err = a.domain.interactor.Delete(uuid, &params); err != nil {
		return err
	}
	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	if err = a.domain.repository.DeleteById(tx, orgId, domain_uuid); err != nil {
		tx.Rollback()
		return err
	}
	if tx.Commit(); tx.Error != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

// TODO Document this method
func (a *application) HostConf(
	ctx echo.Context,
	fqdn string,
	params public.HostConfParams,
) error {
	// TODO Implement this endpoint
	return http.ErrNotSupported
}

// TODO Document this method
func (a *application) CheckHost(
	ctx echo.Context,
	subscriptionManagerId string,
	fqdn string,
	params public.CheckHostParams,
) error {
	// TODO Implement this endpoint
	return http.ErrNotSupported
}

// RegisterIpaDomain (PUT /domains/{uuid}/ipa) initialize the
// IPA domain information into the database. This requires
// a valid X-Rh-IDM-Token. The token is removed when the
// operation is success. Only update information that
// belong to the current organization stored into the
// X-Rh-Identity header, and the host associated to the
// CN is checked against the host inventory, and the list
// of servers into the IPA domain.
// ctx the echo context for the request.
// uuid the domain uuid that identify
// params contains the x-rh-identity, x-rh-insights-request-id
// and x-rh-idm-token header contents.
func (a *application) RegisterDomain(
	ctx echo.Context,
	uuid string,
	params public.RegisterDomainParams,
) error {
	var (
		err   error
		input public.RegisterDomain
		data  *model.Domain
		// host          client.InventoryHost
		domain        *model.Domain
		orgId         string
		tx            *gorm.DB
		output        *public.DomainResponse
		domainCtx     middleware.DomainContextInterface
		clientVersion *header.XRHIDMVersion
	)
	domainCtx = ctx.(middleware.DomainContextInterface)
	if err = ctx.Bind(&input); err != nil {
		return err
	}
	orgId, clientVersion, domain, err = a.domain.interactor.Register(
		domainCtx.XRHID(),
		&params,
		&input,
	)
	if err != nil {
		return err
	}
	ctx.Logger().Info(
		"ipa-hcc",
		clientVersion.IPAHCCVersion,
		"ipa",
		clientVersion.IPAVersion,
	)
	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	// Load Domain data
	if data, err = a.findIpaById(tx, orgId, uuid); err != nil {
		return err
	}

	// Check token
	if err = a.checkToken(
		params.XRhIDMRegistrationToken,
		data.IpaDomain,
	); err != nil {
		return err
	}

	// TODO Check the source host exists in HBI
	// FIXME Set the value from the unencoded and unmarshalled Identity
	xrhid := domainCtx.XRHID()
	if xrhid == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "'xrhid' is nil")
	}
	// subscription_manager_id := xrhid.Identity.System.CommonName
	// if host, err = a.inventory.GetHostByCN(
	// 	params.XRhIdentity,
	// 	params.XRhInsightsRequestId,
	// 	subscription_manager_id,
	// ); err != nil {
	// 	return err
	// }

	// if err = a.existsHostInServers(
	// 	host.FQDN,
	// 	domain.IpaDomain.Servers,
	// ); err != nil {
	// 	return err
	// }

	if err = a.fillIpaDomain(data.IpaDomain, domain.IpaDomain); err != nil {
		return err
	}
	data.IpaDomain.Token = nil
	data.IpaDomain.TokenExpiration = nil

	if *data, err = a.domain.repository.Update(tx, orgId, data); err != nil {
		return err
	}

	if err = a.domain.repository.RhelIdmClearToken(tx, orgId, uuid); err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return tx.Error
	}

	if output, err = a.domain.presenter.Register(data); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, *output)
}
