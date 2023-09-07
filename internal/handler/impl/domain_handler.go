package impl

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
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
	if xrhid, err = getXRHID(ctx); err != nil {
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
		"/api/idmsvc/v1",
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
	UUID uuid.UUID,
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
	if xrhid, err = getXRHID(ctx); err != nil {
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
		UUID,
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

// Delete a Domain resource
// (DELETE /domains/{uuid})
func (a *application) DeleteDomain(
	ctx echo.Context,
	UUID uuid.UUID,
	params public.DeleteDomainParams,
) error {
	var (
		err         error
		tx          *gorm.DB
		orgId       string
		domain_uuid uuid.UUID
		xrhid       *identity.XRHID
	)
	if xrhid, err = getXRHID(ctx); err != nil {
		return err
	}

	if orgId, domain_uuid, err = a.domain.interactor.Delete(
		xrhid,
		UUID,
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

// RegisterIpa (PUT /domains) initialize the
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
	params public.RegisterDomainParams,
) error {
	var (
		err           error
		input         public.Domain
		data          *model.Domain
		orgId         string
		tx            *gorm.DB
		output        *public.RegisterDomainResponse
		clientVersion *header.XRHIDMVersion
		xrhid         *identity.XRHID
	)
	if xrhid, err = getXRHID(ctx); err != nil {
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		return err
	}

	if orgId, clientVersion, data, err = a.domain.interactor.Register(
		a.secrets.domainRegKey,
		xrhid,
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

	if err = a.domain.repository.Register(tx, orgId, data); err != nil {
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

// UpdateDomain (PUT /domains/{uuid}) update the
// IPA domain information into the database. Only update
// information that belong to the current organization stored
// into the X-Rh-Identity header, and the host associated to the
// CN is checked against the host inventory, and the list
// of servers into the IPA domain.
// ctx the echo context for the request.
// UUID the domain uuid that identify
// params contains the x-rh-identity, x-rh-insights-request-id
// and x-rh-idm-version header contents.
func (a *application) UpdateDomainAgent(ctx echo.Context, domain_id uuid.UUID, params public.UpdateDomainAgentParams) error {
	var (
		err         error
		input       public.Domain
		data        *model.Domain
		currentData *model.Domain
		// host          client.InventoryHost
		orgID         string
		tx            *gorm.DB
		output        *public.UpdateDomainAgentResponse
		clientVersion *header.XRHIDMVersion
		xrhid         *identity.XRHID
	)
	if xrhid, err = getXRHID(ctx); err != nil {
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	if orgID, clientVersion, data, err = a.domain.interactor.UpdateAgent(
		xrhid,
		domain_id,
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
	if currentData, err = a.findIpaById(tx, orgID, domain_id); err != nil {
		// FIXME It is not found it should return a 404 Status
		return err
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
	if err = tx.Commit().Error; err != nil {
		return tx.Error
	}

	if output, err = a.domain.presenter.UpdateAgent(currentData); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, *output)
}

// UpdateDomain (PATCH /domains/{uuid}) update the
// IPA domain information into the database. Only update
// information that belong to the current organization stored
// into the X-Rh-Identity header, and the host associated to the
// CN is checked against the host inventory, and the list
// of servers into the IPA domain.
// ctx the echo context for the request.
// UUID the domain uuid that identify
// params contains the x-rh-identity and x-rh-insights-request-id
// header contents.
func (a *application) UpdateDomainUser(ctx echo.Context, domain_id uuid.UUID, params public.UpdateDomainUserParams) error {
	var (
		err         error
		input       public.Domain
		data        *model.Domain
		currentData *model.Domain
		orgID       string
		tx          *gorm.DB
		output      *public.UpdateDomainUserResponse
		xrhid       *identity.XRHID
	)
	if xrhid, err = getXRHID(ctx); err != nil {
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	if orgID, data, err = a.domain.interactor.UpdateUser(
		xrhid,
		domain_id,
		&params,
		&input,
	); err != nil {
		return err
	}
	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	// Load Domain data
	if currentData, err = a.findIpaById(tx, orgID, domain_id); err != nil {
		// FIXME It is not found it should return a 404 Status
		return err
	}

	if err = a.fillDomain(currentData, data); err != nil {
		return err
	}

	if err = a.domain.repository.Update(tx, orgID, currentData); err != nil {
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return tx.Error
	}

	if output, err = a.domain.presenter.UpdateUser(currentData); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, *output)
}

// domains/token route
// Create a domain token for given orgID, domainType, and current time stamp
func (a *application) CreateDomainToken(ctx echo.Context, params public.CreateDomainTokenParams) error {
	var (
		err        error
		input      public.DomainRegTokenRequest
		token      *repository.DomainRegToken
		domainType public.DomainType
		orgID      string
		output     *public.DomainRegToken
		xrhid      *identity.XRHID
	)
	if xrhid, err = getXRHID(ctx); err != nil {
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		return err
	}
	if orgID, domainType, err = a.domain.interactor.CreateDomainToken(
		xrhid,
		&params,
		&input,
	); err != nil {
		return err
	}

	validity := time.Duration(a.config.Application.ExpirationTimeSeconds) * time.Second
	if token, err = a.domain.repository.CreateDomainToken(
		a.secrets.domainRegKey,
		validity,
		orgID,
		domainType,
	); err != nil {
		return err
	}

	// TODO: logging
	// ctx.Logger().Info()

	if output, err = a.domain.presenter.CreateDomainToken(token); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, *output)
}
