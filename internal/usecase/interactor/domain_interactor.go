package interactor

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/domain_token"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type domainInteractor struct{}

// NewDomainInteractor Create an interactor for the /domain endpoint handler
// Return an initialized instance of interactor.DomainInteractor
func NewDomainInteractor() interactor.DomainInteractor {
	return domainInteractor{}
}

// helperDomainTypeToUint transform public.DomainType to an uint const
// Return the uint representation or model.DomainTypeUndefined if it does not match
// the current types.
func helperDomainTypeToUint(domainType public.DomainType) uint {
	switch domainType {
	case api_public.RhelIdm:
		return model.DomainTypeIpa
	default:
		return model.DomainTypeUndefined
	}
}

// func (i domainInteractor) PartialUpdate(id public.Id, params *api_public.PartialUpdateTodoParams, in *api_public.Todo, out *model.Todo) error {
// 	if id <= 0 {
// 		return fmt.Errorf("'id' should be a positive int64")
// 	}
// 	if params == nil {
// 		return fmt.Errorf("'params' cannot be nil")
// 	}
// 	if in == nil {
// 		return fmt.Errorf("'in' cannot be nil")
// 	}
// 	if out == nil {
// 		return fmt.Errorf("'out' cannot be nil")
// 	}
// 	out.Model.ID = id
// 	if in.Title != nil {
// 		out.Title = pointy.String(*in.Title)
// 	}
// 	if in.Body != nil {
// 		out.Description = pointy.String(*in.Body)
// 	}
// 	return nil
// }

// func (i domainInteractor) FullUpdate(id public.Id, params *api_public.UpdateTodoParams, in *api_public.Todo, out *model.Todo) error {
// 	if id <= 0 {
// 		return fmt.Errorf("'id' should be a positive int64")
// 	}
// 	if params == nil {
// 		return fmt.Errorf("'params' cannot be nil")
// 	}
// 	if in == nil {
// 		return fmt.Errorf("'in' cannot be nil")
// 	}
// 	if out == nil {
// 		return fmt.Errorf("'out' cannot be nil")
// 	}
// 	out.ID = id
// 	out.Title = pointy.String(*in.Title)
// 	out.Description = pointy.String(*in.Body)
// 	return nil
// }

// TODO Document method
func (i domainInteractor) Delete(xrhid *identity.XRHID, UUID uuid.UUID, params *api_public.DeleteDomainParams) (string, uuid.UUID, error) {
	if xrhid == nil {
		return "", uuid.Nil, internal_errors.NilArgError("xrhid")
	}
	if params == nil {
		return "", uuid.Nil, fmt.Errorf("'params' cannot be nil")
	}
	return xrhid.Identity.OrgID, UUID, nil
}

// List is the input adapter to list the domains that belongs to
// the current organization by using pagination.
// params is the pagination parameters.
// Return the organization id, offset of the page, number of items
// to retrieve and nil error for a success scenario, else empty or
// zero values and an error interface filled on error.
func (i domainInteractor) List(xrhid *identity.XRHID, params *api_public.ListDomainsParams) (orgID string, offset int, limit int, err error) {
	if xrhid == nil {
		return "", -1, -1, internal_errors.NilArgError("xrhid")
	}
	if params == nil {
		return "", -1, -1, internal_errors.NilArgError("params")
	}
	if params.Offset == nil {
		offset = 0
	} else {
		offset = *params.Offset
	}
	if params.Limit == nil {
		limit = 10
	} else {
		limit = *params.Limit
	}
	return xrhid.Identity.OrgID, offset, limit, nil
}

// GetByID translate from input api to model information.
// xrhid is the unserialized identity structure stored into the request
// context.
// params is the GET /domains/:uuid endpoint parameters.
// Return the organization id and nil error for success invokation, else
// an empty organizaion id and a filled error with the situation details.
func (i domainInteractor) GetByID(xrhid *identity.XRHID, params *public.ReadDomainParams) (orgID string, err error) {
	if xrhid == nil {
		return "", internal_errors.NilArgError("xrhid")
	}
	if params == nil {
		return "", internal_errors.NilArgError("params")
	}

	return xrhid.Identity.OrgID, nil
}

// Register translates the API input format into the business
// data models for the PUT /domains/{uuid}/register endpoint.
// params contains the header parameters.
// body contains the input payload.
// Return the orgId and the business model for Ipa information,
// when success translation, else it returns empty string for orgId,
// nil for the Ipa data, and an error filled.
func (i domainInteractor) Register(domainRegKey []byte, xrhid *identity.XRHID, params *api_public.RegisterDomainParams, body *public.Domain) (string, *header.XRHIDMVersion, *model.Domain, error) {
	var (
		domainID uuid.UUID
		domain   *model.Domain
		err      error
	)
	if err = i.guardRegister(xrhid, params, body); err != nil {
		return "", nil, nil, err
	}
	orgID := xrhid.Identity.OrgID

	// Retrieve the ipa-hcc version information
	clientVersion := header.NewXRHIDMVersionWithHeader(params.XRhIdmVersion)
	if clientVersion == nil {
		return "", nil, nil, fmt.Errorf("'X-Rh-Idm-Version' is invalid")
	}

	// verify token
	if domainID, err = domain_token.VerifyDomainRegistrationToken(
		domainRegKey,
		string(body.DomainType),
		orgID,
		domain_token.DomainRegistrationToken(params.XRhIdmRegistrationToken),
	); err != nil {
		msg := fmt.Sprintf("Domain registration token is invalid: %s", err)
		return "", nil, nil, echo.NewHTTPError(http.StatusUnauthorized, msg)
	}

	// Read the body payload
	if domain, err = i.commonRegisterUpdate(orgID, domainID, body); err != nil {
		return "", nil, nil, err
	}

	// Agent request body has no/nil title, description, and auto enrollment flag
	// default title is domain name
	domain.Title = pointy.String(body.DomainName)
	domain.Description = pointy.String("")
	// new domains are disabled by default
	domain.AutoEnrollmentEnabled = pointy.Bool(false)

	return orgID, clientVersion, domain, nil
}

// Update translates the API input format into the business
// data models for the PUT /domains/{uuid} endpoint.
// params contains the header parameters.
// body contains the input payload.
// Return the orgId and the business model for Ipa information,
// when success translation, else it returns empty string for orgId,
// nil for the Ipa data, and an error filled.
func (i domainInteractor) UpdateAgent(xrhid *identity.XRHID, UUID uuid.UUID, params *api_public.UpdateDomainAgentParams, body *public.Domain) (string, *header.XRHIDMVersion, *model.Domain, error) {
	var (
		domain *model.Domain
		err    error
	)
	if err = i.guardUpdate(xrhid, UUID, body); err != nil {
		return "", nil, nil, err
	}
	if params == nil {
		return "", nil, nil, internal_errors.NilArgError("params")
	}
	orgID := xrhid.Identity.Internal.OrgID

	// Retrieve the ipa-hcc version information
	clientVersion := header.NewXRHIDMVersionWithHeader(params.XRhIdmVersion)
	if clientVersion == nil {
		return "", nil, nil, fmt.Errorf("'X-Rh-Idm-Version' is invalid")
	}

	// Read the body payload
	if domain, err = i.commonRegisterUpdate(orgID, UUID, body); err != nil {
		return "", nil, nil, err
	}
	return orgID, clientVersion, domain, nil
}

// Update translates the API input format into the business
// data models for the PATCH /domains/{uuid} endpoint.
// params contains the header parameters.
// body contains the input payload.
// Return the orgId and the business model for Ipa information,
// when success translation, else it returns empty string for orgId,
// nil for the Ipa data, and an error filled.
func (i domainInteractor) UpdateUser(xrhid *identity.XRHID, UUID uuid.UUID, params *api_public.UpdateDomainUserParams, body *public.Domain) (string, *model.Domain, error) {
	var (
		domain *model.Domain
		err    error
	)
	if err = i.guardUpdate(xrhid, UUID, body); err != nil {
		return "", nil, err
	}
	if params == nil {
		return "", nil, internal_errors.NilArgError("params")
	}
	orgID := xrhid.Identity.Internal.OrgID

	// Read the body payload
	domain = i.commonRegisterUpdateUser(orgID, UUID, body)
	return orgID, domain, nil
}

// Create domain registration token /domains/token
//
// Verify input parameters and check for supported domain types.
func (i domainInteractor) CreateDomainToken(
	xrhid *identity.XRHID,
	params *public.CreateDomainTokenParams,
	body *public.DomainRegTokenRequest,
) (orgID string, domainType public.DomainType, err error) {
	if xrhid == nil {
		return "", "", internal_errors.NilArgError("xrhid")
	}
	if xrhid.Identity.Type != "User" {
		return "", "", fmt.Errorf("invalid identity type '%s'", xrhid.Identity.Type)
	}
	if params == nil {
		return "", "", internal_errors.NilArgError("params")
	}
	if body == nil {
		return "", "", internal_errors.NilArgError("body")
	}

	orgID = xrhid.Identity.OrgID

	switch body.DomainType {
	case api_public.DomainType(api_public.RhelIdm):
		domainType = body.DomainType
	default:
		return "", "", fmt.Errorf("Unsupported domain_type='%s'", body.DomainType)
	}

	return orgID, domainType, nil
}

// --------- Private methods -----------

func (i domainInteractor) registerOrUpdateRhelIdm(body *public.Domain, domainIpa *model.Ipa) error {
	domainIpa.RealmName = pointy.String(body.RhelIdm.RealmName)

	// Translate realm domains
	i.registerOrUpdateRhelIdmRealmDomains(body, domainIpa)

	// Certificate list
	i.registerOrUpdateRhelIdmCaCerts(body, domainIpa)

	// Server list
	i.registerOrUpdateRhelIdmServers(body, domainIpa)

	// Location list
	i.registerOrUpdateRhelIdmLocations(body, domainIpa)

	return nil
}

func (i domainInteractor) registerOrUpdateRhelIdmRealmDomains(body *public.Domain, domainIpa *model.Ipa) {
	if body.RhelIdm.RealmDomains == nil {
		domainIpa.RealmDomains = pq.StringArray{}
		return
	}
	domainIpa.RealmDomains = make(pq.StringArray, 0)
	domainIpa.RealmDomains = append(
		domainIpa.RealmDomains,
		body.RhelIdm.RealmDomains...,
	)
}

func (i domainInteractor) registerOrUpdateRhelIdmCaCerts(body *public.Domain, domainIpa *model.Ipa) {
	if body.RhelIdm.CaCerts == nil {
		domainIpa.CaCerts = []model.IpaCert{}
		return
	}
	domainIpa.CaCerts = make([]model.IpaCert, len(body.RhelIdm.CaCerts))
	for idx := range body.RhelIdm.CaCerts {
		i.registerOrUpdateRhelIdmCaCertOne(&domainIpa.CaCerts[idx], &body.RhelIdm.CaCerts[idx])
	}
}

func (i domainInteractor) registerOrUpdateRhelIdmCaCertOne(caCert *model.IpaCert, cert *api_public.Certificate) {
	caCert.Nickname = cert.Nickname
	caCert.Issuer = cert.Issuer
	caCert.Subject = cert.Subject
	caCert.SerialNumber = cert.SerialNumber
	caCert.NotBefore = cert.NotBefore
	caCert.NotAfter = cert.NotAfter
	caCert.Pem = cert.Pem
}

func (i domainInteractor) registerOrUpdateRhelIdmServers(body *public.Domain, domainIpa *model.Ipa) {
	if body.RhelIdm.Servers == nil {
		domainIpa.Servers = []model.IpaServer{}
		return
	}
	domainIpa.Servers = make([]model.IpaServer, len(body.RhelIdm.Servers))
	for idx, server := range body.RhelIdm.Servers {
		domainIpa.Servers[idx].FQDN = server.Fqdn
		if server.SubscriptionManagerId != nil {
			domainIpa.Servers[idx].RHSMId = pointy.String(server.SubscriptionManagerId.String())
		}
		domainIpa.Servers[idx].Location = server.Location
		domainIpa.Servers[idx].PKInitServer = server.PkinitServer
		domainIpa.Servers[idx].CaServer = server.CaServer
		domainIpa.Servers[idx].HCCEnrollmentServer = server.HccEnrollmentServer
		domainIpa.Servers[idx].HCCUpdateServer = server.HccUpdateServer
	}
}

func (i domainInteractor) registerOrUpdateRhelIdmLocations(body *public.Domain, domainIpa *model.Ipa) {
	if body.RhelIdm.Locations == nil {
		domainIpa.Locations = []model.IpaLocation{}
		return
	}
	domainIpa.Locations = make([]model.IpaLocation, len(body.RhelIdm.Locations))
	for idx, location := range body.RhelIdm.Locations {
		domainIpa.Locations[idx].Name = location.Name
		domainIpa.Locations[idx].Description = location.Description
	}
}

func (i domainInteractor) guardRegister(xrhid *identity.XRHID, params *api_public.RegisterDomainParams, body *public.Domain) (err error) {
	if xrhid == nil {
		return internal_errors.NilArgError("xrhid")
	}
	if params == nil {
		return internal_errors.NilArgError("params")
	}
	if body == nil {
		return internal_errors.NilArgError("body")
	}

	return nil
}

func (i domainInteractor) guardUpdate(xrhid *identity.XRHID, UUID uuid.UUID, body *public.Domain) (err error) {
	if xrhid == nil {
		return internal_errors.NilArgError("xrhid")
	}
	if UUID == uuid.Nil {
		return fmt.Errorf("'UUID' is invalid")
	}
	if body == nil {
		return internal_errors.NilArgError("body")
	}

	return nil
}

func (i domainInteractor) commonRegisterUpdate(orgID string, UUID uuid.UUID, body *public.Domain) (domain *model.Domain, err error) {
	domain = &model.Domain{}
	domain.OrgId = orgID
	domain.DomainUuid = UUID
	domain.Title = body.Title
	domain.Description = body.Description
	domain.AutoEnrollmentEnabled = body.AutoEnrollmentEnabled
	domain.DomainName = pointy.String(body.DomainName)
	switch body.DomainType {
	case api_public.DomainType(api_public.RhelIdm):
		domain.Type = pointy.Uint(model.DomainTypeIpa)
		domain.IpaDomain = &model.Ipa{}
		err = i.registerOrUpdateRhelIdm(body, domain.IpaDomain)
	default:
		err = fmt.Errorf("Unsupported domain_type='%s'", body.DomainType)
	}
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (i domainInteractor) commonRegisterUpdateUser(orgID string, UUID uuid.UUID, body *public.Domain) (domain *model.Domain) {
	// Only copy the necessary information for the associated handler
	domain = &model.Domain{}
	domain.OrgId = orgID
	domain.DomainUuid = UUID
	domain.Title = body.Title
	domain.Description = body.Description
	domain.AutoEnrollmentEnabled = body.AutoEnrollmentEnabled
	return domain
}
