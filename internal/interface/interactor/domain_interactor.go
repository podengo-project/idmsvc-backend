package interactor

import (
	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

type DomainInteractor interface {
	Delete(xrhid *identity.XRHID, UUID uuid.UUID, params *api_public.DeleteDomainParams) (string, uuid.UUID, error)
	List(xrhid *identity.XRHID, params *api_public.ListDomainsParams) (orgID string, offset, limit int, err error)
	GetByID(xrhid *identity.XRHID, params *public.ReadDomainParams) (orgID string, err error)
	Register(domainRegKey []byte, xrhid *identity.XRHID, params *api_public.RegisterDomainParams, body *api_public.Domain) (string, *header.XRHIDMVersion, *model.Domain, error)
	UpdateAgent(xrhid *identity.XRHID, UUID uuid.UUID, params *api_public.UpdateDomainAgentParams, body *api_public.UpdateDomainAgentRequest) (string, *header.XRHIDMVersion, *model.Domain, error)
	UpdateUser(xrhid *identity.XRHID, UUID uuid.UUID, params *api_public.UpdateDomainUserParams, body *api_public.UpdateDomainUserRequest) (string, *model.Domain, error)
	CreateDomainToken(xrhid *identity.XRHID, params *api_public.CreateDomainTokenParams, body *api_public.DomainRegTokenRequest) (orgID string, domainType public.DomainType, err error)
}
