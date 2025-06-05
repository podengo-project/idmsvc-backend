package impl

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
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
	handlerName := "ListDomains"
	logger := app_context.LogFromCtx(ctx.Request().Context())
	logger = logger.With(slog.String("handler", handlerName))
	if xrhid, err = getXRHID(ctx); err != nil {
		logger.Error(errXRHIDIsNil)
		return err
	}
	// TODO A call to an internal validator could be here to check public.ListTodosParams
	if orgID, offset, limit, err = a.domain.interactor.List(xrhid, &params); err != nil {
		logger.Error(errInputAdapter)
		return err
	}
	logger = logger.With(
		slog.Int("offset", offset),
		slog.Int("limit", limit),
	)
	if tx = a.db.Begin(); tx.Error != nil {
		logger.Error(errDBTXBegin)
		return tx.Error
	}
	// https://stackoverflow.com/a/46421989
	defer tx.Rollback()
	c := app_context.CtxWithDB(ctx.Request().Context(), tx)
	if data, count, err = a.domain.repository.List(
		c,
		orgID,
		offset,
		limit,
	); err != nil {
		logger.Error("failed to list domains from the database")
		return err
	}
	if tx.Commit(); tx.Error != nil {
		logger.Error(errDBTXCommit)
		return tx.Error
	}
	// TODO Read prefix from configuration
	if output, err = a.domain.presenter.List(
		count,
		offset,
		limit,
		data,
	); err != nil {
		logger.Error(errOutputAdapter)
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
	handlerName := "ReadDomain"
	logger := app_context.LogFromCtx(ctx.Request().Context())
	logger = logger.With(
		slog.String("handler", handlerName),
		slog.String("uuid", UUID.String()),
	)
	if xrhid, err = getXRHID(ctx); err != nil {
		logger.Error(errXRHIDIsNil)
		return err
	}

	if orgID, err = a.domain.interactor.GetByID(
		xrhid,
		&params,
	); err != nil {
		logger.Error(errInputAdapter)
		return err
	}
	if tx = a.db.Begin(); tx.Error != nil {
		logger.Error(errDBTXBegin)
		return tx.Error
	}
	defer tx.Rollback()
	c := app_context.CtxWithDB(ctx.Request().Context(), tx)
	if data, err = a.domain.repository.FindByID(
		c,
		orgID,
		UUID,
	); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(errDBNotFound)
			return internal_errors.NewHTTPErrorF(
				http.StatusNotFound,
				"cannot read unknown domain '%s'.",
				UUID.String(),
			)
		}
		logger.Error("failed to find a domain by ID on the database")
		return err
	}
	if err = tx.Commit().Error; err != nil {
		logger.Error(errDBTXCommit)
		return err
	}
	if output, err = a.domain.presenter.Get(data); err != nil {
		logger.Error(errOutputAdapter)
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
		err        error
		tx         *gorm.DB
		orgId      string
		domainUUID uuid.UUID
		xrhid      *identity.XRHID
	)
	handlerName := "DeleteDomain"
	logger := app_context.LogFromCtx(ctx.Request().Context())
	logger = logger.With(slog.String("handler", handlerName))
	if xrhid, err = getXRHID(ctx); err != nil {
		logger.Error(errXRHIDIsNil)
		return err
	}

	if orgId, domainUUID, err = a.domain.interactor.Delete(
		xrhid,
		UUID,
		&params,
	); err != nil {
		logger.Error(errInputAdapter)
		return err
	}
	logger = logger.With(slog.String("uuid", UUID.String()))
	if tx = a.db.Begin(); tx.Error != nil {
		logger.Error(errDBTXBegin)
		return tx.Error
	}
	defer tx.Rollback()
	c := app_context.CtxWithDB(ctx.Request().Context(), tx)
	if err = a.domain.repository.DeleteById(
		c,
		orgId,
		domainUUID,
	); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(errDBNotFound)
			return internal_errors.NewHTTPErrorF(
				http.StatusNotFound,
				"cannot delete unknown domain '%s'.",
				UUID.String(),
			)
		}
		logger.Error("failed to delete domain by ID on the database")
		return err
	}
	if tx.Commit(); tx.Error != nil {
		logger.Error(errDBTXCommit)
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

// RegisterDomain (PUT /domains) initialize the
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
	handlerName := "RegisterDomain"
	logger := app_context.LogFromCtx(ctx.Request().Context())
	logger = logger.With(slog.String("handler", handlerName))
	if xrhid, err = getXRHID(ctx); err != nil {
		logger.Error(errXRHIDIsNil)
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		logger.Error(errUnserializing)
		return err
	}

	if orgId, clientVersion, data, err = a.domain.interactor.Register(
		a.config.Secrets.DomainRegKey,
		xrhid,
		&params,
		&input,
	); err != nil {
		logger.Error(errInputAdapter)
		return err
	}
	logger.Info("ipa-hcc client version",
		slog.Group("client-version",
			slog.String("ipa-hcc", clientVersion.IPAHCCVersion),
			slog.String("ipa", clientVersion.IPAVersion),
			slog.String("os-release-id", clientVersion.OSReleaseID),
			slog.String("os-release-version", clientVersion.OSReleaseVersionID),
		),
	)

	updateServerRSHMId := xrhid.Identity.System.CommonName
	if err = ensureUpdateServerEnabledForUpdates(
		ctx.Request().Context(),
		updateServerRSHMId,
		data.IpaDomain.Servers,
	); err != nil {
		logger.Error("failed to ensure that the requesting server is authorized for updating the data, on register domain process")
		return err
	}

	if tx = a.db.Begin(); tx.Error != nil {
		logger.Error(errDBTXCommit)
		return tx.Error
	}
	defer tx.Rollback()

	c := app_context.CtxWithDB(ctx.Request().Context(), tx)
	if err = a.domain.repository.Register(c, orgId, data); err != nil {
		logger.Error("failed to register domain on the database")
		return err
	}

	if err = tx.Commit().Error; err != nil {
		logger.Error(errDBTXCommit)
		return tx.Error
	}

	if output, err = a.domain.presenter.Register(data); err != nil {
		logger.Error(errOutputAdapter)
		return err
	}

	return ctx.JSON(http.StatusCreated, *output)
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
		err           error
		input         public.UpdateDomainAgentRequest
		data          *model.Domain
		currentData   *model.Domain
		orgID         string
		tx            *gorm.DB
		output        *public.UpdateDomainAgentResponse
		clientVersion *header.XRHIDMVersion
		xrhid         *identity.XRHID
	)
	handlerName := "UpdateDomainAgent"
	logger := app_context.LogFromCtx(ctx.Request().Context())
	logger = logger.With(slog.String("handler", handlerName))
	if xrhid, err = getXRHID(ctx); err != nil {
		logger.Error(errXRHIDIsNil)
		return err
	}
	logger = logger.With(slog.String("uuid", domain_id.String()))

	if err = ctx.Bind(&input); err != nil {
		logger.Error(errUnserializing)
		return err
	}
	if orgID, clientVersion, data, err = a.domain.interactor.UpdateAgent(
		xrhid,
		domain_id,
		&params,
		&input,
	); err != nil {
		logger.Error(errInputAdapter)
		return err
	}
	logger.Info("ipa-hcc client version",
		slog.Group("client-version",
			slog.String("ipa-hcc", clientVersion.IPAHCCVersion),
			slog.String("ipa", clientVersion.IPAVersion),
			slog.String("os-release-id", clientVersion.OSReleaseID),
			slog.String("os-release-version", clientVersion.OSReleaseVersionID),
		),
	)
	if tx = a.db.Begin(); tx.Error != nil {
		logger.Error(errDBTXBegin)
		return tx.Error
	}
	defer tx.Rollback()

	// Check that the update server is included in the request
	updateServerRSHMId := xrhid.Identity.System.CommonName
	logger = logger.With(slog.String("server", updateServerRSHMId))
	if err = ensureUpdateServerEnabledForUpdates(
		ctx.Request().Context(),
		updateServerRSHMId,
		data.IpaDomain.Servers,
	); err != nil {
		logger.Error("failed to ensure that the requesting server is authorized for updating the data, on updating domain process from a system")
		return err
	}

	// Load Domain data
	c := app_context.CtxWithDB(ctx.Request().Context(), tx)
	if currentData, err = a.domain.repository.FindByID(c, orgID, domain_id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(errDBNotFound)
			return internal_errors.NewHTTPErrorF(
				http.StatusNotFound,
				"%s",
				err.Error(),
			)
		}
		logger.Error(errDBGeneralError)
		return err
	}

	if err = ensureSubscriptionManagerIDAuthorizedToUpdate(
		c,
		updateServerRSHMId,
		currentData.IpaDomain.Servers,
	); err != nil {
		logger.Error("failed because the requesting server is not authorized to update the domain")
		return err
	}

	if data.DomainName != nil &&
		currentData.DomainName != nil &&
		*data.DomainName != *currentData.DomainName {
		logger.Error("failed because domain_name is immutable and cannot be modified, on updating domain data from a server",
			slog.String("domain_name", *currentData.DomainName),
		)
		return internal_errors.NewHTTPErrorF(
			http.StatusBadRequest,
			"'domain_name' may not be changed",
		)
	}

	if data.IpaDomain != nil && currentData.IpaDomain != nil &&
		data.IpaDomain.RealmName != nil && currentData.IpaDomain.RealmName != nil &&
		*data.IpaDomain.RealmName != *currentData.IpaDomain.RealmName {
		logger.Error("failed because realm_name is immutable and cannot be modified, on updating domain data from a server",
			slog.String("realm_name", *currentData.IpaDomain.RealmName),
		)
		return internal_errors.NewHTTPErrorF(
			http.StatusBadRequest,
			"'realm_name' may not be changed",
		)
	}

	if err = a.fillDomain(currentData, data); err != nil {
		logger.Error("failed to fill the new domain information for an agent update")
		return err
	}

	if err = a.domain.repository.UpdateAgent(c, orgID, currentData); err != nil {
		logger.Error("failed to update the new data in the database")
		return err
	}
	if err = tx.Commit().Error; err != nil {
		logger.Error(errDBTXCommit)
		return tx.Error
	}

	if output, err = a.domain.presenter.UpdateAgent(currentData); err != nil {
		logger.Error(errOutputAdapter)
		return err
	}

	return ctx.JSON(http.StatusOK, *output)
}

// UpdateDomainUser (PATCH /domains/{uuid}) update the
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
		input       public.UpdateDomainUserRequest
		data        *model.Domain
		currentData *model.Domain
		orgID       string
		tx          *gorm.DB
		output      *public.UpdateDomainUserResponse
		xrhid       *identity.XRHID
	)
	handlerName := "UpdateDomainUser"
	logger := app_context.LogFromCtx(ctx.Request().Context())
	logger = logger.With(
		slog.String("handler", handlerName),
		slog.String("uuid", domain_id.String()),
	)
	if xrhid, err = getXRHID(ctx); err != nil {
		logger.Error(errXRHIDIsNil)
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		logger.Error(errUnserializing)
		return err
	}
	if orgID, data, err = a.domain.interactor.UpdateUser(
		xrhid,
		domain_id,
		&params,
		&input,
	); err != nil {
		logger.Error(errInputAdapter)
		return err
	}
	if tx = a.db.Begin(); tx.Error != nil {
		logger.Error(errDBTXBegin)
		return tx.Error
	}
	defer tx.Rollback()

	// Load Domain data
	c := app_context.CtxWithDB(ctx.Request().Context(), tx)
	if currentData, err = a.domain.repository.FindByID(c, orgID, domain_id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(errDBNotFound)
			return internal_errors.NewHTTPErrorF(
				http.StatusNotFound,
				"%s",
				err.Error(),
			)
		}
		logger.Error(errDBGeneralError)
		return err
	}

	if err = a.fillDomainUser(currentData, data); err != nil {
		logger.Error("failed to fill the domain information for a user update")
		return err
	}

	if err = a.domain.repository.UpdateUser(c, orgID, data); err != nil {
		logger.Error("failed to update domain information in the database for a user update")
		return err
	}
	if err = tx.Commit().Error; err != nil {
		logger.Error(errDBTXCommit)
		return tx.Error
	}

	if output, err = a.domain.presenter.UpdateUser(currentData); err != nil {
		logger.Error(errOutputAdapter)
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
	handlerName := "CreateDomainToken"
	logger := app_context.LogFromCtx(ctx.Request().Context())
	logger = logger.With(slog.String("handler", handlerName))
	if xrhid, err = getXRHID(ctx); err != nil {
		logger.Error(errXRHIDIsNil, slog.String("handler", handlerName))
		return err
	}

	if err = ctx.Bind(&input); err != nil {
		logger.Error(errUnserializing)
		return err
	}
	if orgID, domainType, err = a.domain.interactor.CreateDomainToken(
		xrhid,
		&params,
		&input,
	); err != nil {
		logger.Error(errInputAdapter)
		return err
	}
	logger = logger.With(slog.String("domain_type", string(domainType)))

	validity := time.Duration(a.config.Application.TokenExpirationTimeSeconds) * time.Second
	if token, err = a.domain.repository.CreateDomainToken(
		ctx.Request().Context(),
		a.config.Secrets.DomainRegKey,
		validity,
		orgID,
		domainType,
	); err != nil {
		logger.Error("failed to create a registration token")
		return err
	}

	if output, err = a.domain.presenter.CreateDomainToken(token); err != nil {
		logger.Error(errOutputAdapter)
		return err
	}

	return ctx.JSON(http.StatusOK, *output)
}
