package impl

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hmsidm/internal/api/header"
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/infrastructure/middleware"
	"github.com/labstack/echo/v4"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"gorm.io/gorm"
)

// About defer Rollback
//
// In the case of a commit, the rollback operation does not have any effect into
// the transaction because it was already committed; but it returns ErrTxDone
// error; no side effects because the session created by the transaction is not
// used anymore when the handler returns. It was double checked debugging the code.

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
		orgID  string
		offset int
		limit  int
		count  int64
		tx     *gorm.DB
		xrhid  *identity.XRHID
	)
	xrhid, err = getXRHID(ctx)
	if err != nil {
		return err
	}
	// TODO A call to an internal validator could be here to check public.ListTodosParams
	if orgID, offset, limit, err = a.domain.interactor.List(xrhid, &params); err != nil {
		return err
	}
	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	// https://stackoverflow.com/a/46421989
	defer tx.Rollback()
	if data, count, err = a.domain.repository.List(
		tx,
		orgID,
		offset,
		limit,
	); err != nil {
		return err
	}
	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}
	// TODO Read prefix from configuration
	if output, err = a.domain.presenter.List(
		"/api/hmsidm/v1",
		count,
		offset,
		limit,
		data,
	); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *output)
}

// ReadDomain retrieve a domain resource identified by the uuid for
// the GET /domains/:id endpoint.
// ctx is the echo.Context for this request.
// uuid is the identifier for the domain to be retrieved.
// params represent the parameters for the /domains/:uuid endpoint.
// Return nil if the handler execute successfully, else an error
// interface providing the error details.
func (a *application) ReadDomain(
	ctx echo.Context,
	uuid string,
	params public.ReadDomainParams,
) error {
	var (
		err    error
		data   *model.Domain
		output *public.Domain
		orgID  string
		tx     *gorm.DB
		xrhid  *identity.XRHID
	)
	xrhid, err = getXRHID(ctx)
	if err != nil {
		return err
	}

	if orgID, err = a.domain.interactor.GetByID(
		xrhid,
		&params,
	); err != nil {
		return err
	}
	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()
	if data, err = a.domain.repository.FindByID(
		tx,
		orgID,
		uuid,
	); err != nil {
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return err
	}
	if output, err = a.domain.presenter.Get(data); err != nil {
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
		err      error
		input    public.CreateDomain
		orgId    string
		data     *model.Domain
		output   *public.CreateDomainResponse
		tx       *gorm.DB
		tokenStr string
		xrhid    *identity.XRHID
	)
	xrhid, err = getXRHID(ctx)
	if err != nil {
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	if orgId, data, err = a.domain.interactor.Create(
		xrhid,
		&params,
		&input,
	); err != nil {
		return err
	}

	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	// https://stackoverflow.com/a/46421989
	if err = a.domain.repository.Create(tx, orgId, data); err != nil {
		return err
	}
	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}
	if output, err = a.domain.presenter.Create(data); err != nil {
		return err
	}

	// Add X-Rh-Idm-RhelIdm-Register-Token
	if tokenStr, err = header.EncodeRhelIdmToken(
		&header.RhelIdmToken{
			Secret:     data.IpaDomain.Token,
			Expiration: data.IpaDomain.TokenExpiration,
		},
	); err != nil {
		return err
	}
	ctx.Response().Header().Add(
		header.XRHIDMRHELIDMRegisterToken,
		tokenStr,
	)

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
		xrhid *identity.XRHID
	)
	xrhid, err = getXRHID(ctx)
	if err != nil {
		return err
	}
	if orgId, domain_uuid, err = a.domain.interactor.Delete(
		xrhid,
		uuid,
		&params,
	); err != nil {
		return err
	}
	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()
	if err = a.domain.repository.DeleteById(
		tx,
		orgId,
		domain_uuid,
	); err != nil {
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

// RegisterIpaDomain (PUT /domains/{uuid}/register) initialize the
// IPA domain information into the database. This requires
// a valid X-Rh-IDM-Token. The token is removed when the
// operation is success. Only update information that
// belong to the current organization stored into the
// X-Rh-Identity header.
// ctx the echo context for the request.
// UUID the domain uuid that identify
// params contains the x-rh-identity, x-rh-insights-request-id
// and x-rh-idm-token header contents.
func (a *application) RegisterDomain(
	ctx echo.Context,
	UUID string,
	params public.RegisterDomainParams,
) error {
	var (
		err     error
		input   public.Domain
		data    *model.Domain
		oldData *model.Domain
		// host          client.InventoryHost
		orgId         string
		tx            *gorm.DB
		output        *public.Domain
		clientVersion *header.XRHIDMVersion
		xrhid         *identity.XRHID
	)
	xrhid, err = getXRHID(ctx)
	if err != nil {
		return err
	}
	if err = ctx.Bind(&input); err != nil {
		return err
	}
	if orgId, clientVersion, data, err = a.domain.interactor.Register(
		xrhid,
		UUID,
		&params,
		&input,
	); err != nil {
		return err
	}
	ctx.Logger().Info(
		"ipa-hcc",
		clientVersion.IPAHCCVersion,
		"ipa",
		clientVersion.IPAVersion,
		"os-release-id",
		clientVersion.OSReleaseID,
		"os-release-version-id",
		clientVersion.OSReleaseVersionID,
	)
	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	// Load Domain data
	if oldData, err = a.findIpaById(tx, orgId, UUID); err != nil {
		// FIXME It is not found it should return a 404 Status
		return err
	}

	// Check token
	if err = a.checkToken(
		params.XRhIdmRegistrationToken,
		oldData.IpaDomain,
	); err != nil {
		return err
	}

	data.IpaDomain.Token = nil
	data.IpaDomain.TokenExpiration = nil
	data.DomainUuid = uuid.MustParse(UUID)
	data.ID = oldData.ID

	if err = a.domain.repository.Update(tx, orgId, data); err != nil {
		return err
	}

	if err = a.domain.repository.RhelIdmClearToken(tx, orgId, UUID); err != nil {
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

// UpdateDomain (PUT /domains/{uuid}/update) update the
// IPA domain information into the database. Only update
// information that belong to the current organization stored
// into the X-Rh-Identity header, and the host associated to the
// CN is checked against the host inventory, and the list
// of servers into the IPA domain.
// ctx the echo context for the request.
// UUID the domain uuid that identify
// params contains the x-rh-identity, x-rh-insights-request-id
// and x-rh-idm-token header contents.
func (a *application) UpdateDomain(ctx echo.Context, UUID string, params public.UpdateDomainParams) error {
	var (
		err         error
		input       public.Domain
		data        *model.Domain
		currentData *model.Domain
		// host          client.InventoryHost
		orgID         string
		tx            *gorm.DB
		output        *public.Domain
		clientVersion *header.XRHIDMVersion
		xrhid         *identity.XRHID
	)
	xrhid, err = getXRHID(ctx)
	if err != nil {
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	if orgID, clientVersion, data, err = a.domain.interactor.Update(
		xrhid,
		UUID,
		&params,
		&input,
	); err != nil {
		return err
	}
	ctx.Logger().Info(
		"{",
		"\"ipa-hcc\":\"", clientVersion.IPAHCCVersion, "\",",
		"\"ipa\": \"", clientVersion.IPAVersion, "\",",
		"\"os-release-id\": \"", clientVersion.OSReleaseID, "\",",
		"\"os-release-version-id\": \"", clientVersion.OSReleaseVersionID, "\"",
		"}",
	)
	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	// Load Domain data
	if currentData, err = a.findIpaById(tx, orgID, UUID); err != nil {
		// FIXME It is not found it should return a 404 Status
		return err
	}

	if currentData.IpaDomain.Token != nil || currentData.IpaDomain.TokenExpiration != nil {
		return fmt.Errorf("Bad Request")
	}

	subscriptionManagerID := xrhid.Identity.System.CommonName
	if err = a.isSubscriptionManagerIDAuthorizedToUpdate(
		subscriptionManagerID,
		currentData.IpaDomain.Servers,
	); err != nil {
		return err
	}

	if err = a.fillDomain(currentData, data); err != nil {
		return err
	}

	if err = a.domain.repository.Update(tx, orgID, currentData); err != nil {
		return err
	}

	if err = a.domain.repository.RhelIdmClearToken(tx, orgID, currentData.DomainUuid.String()); err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return tx.Error
	}

	if output, err = a.domain.presenter.Update(currentData); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, *output)
}

func getXRHID(ctx echo.Context) (*identity.XRHID, error) {
	domainCtx := ctx.(middleware.DomainContextInterface)
	xrhid := domainCtx.XRHID()
	if xrhid == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "'xrhid' is nil")
	} else {
		return xrhid, nil
	}
}
