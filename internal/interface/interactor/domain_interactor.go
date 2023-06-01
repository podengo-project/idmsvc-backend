package interactor

import (
	"github.com/hmsidm/internal/api/header"
	"github.com/hmsidm/internal/api/public"
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type DomainInteractor interface {
	Create(params *api_public.CreateDomainParams, body *api_public.CreateDomain) (string, *model.Domain, error)
	Delete(uuid string, params *api_public.DeleteDomainParams) (string, string, error)
	List(xrhid *identity.XRHID, params *api_public.ListDomainsParams) (orgID string, offset int, limit int, err error)
	GetByID(xrhid *identity.XRHID, params *public.ReadDomainParams) (orgID string, err error)
	Register(xrhid *identity.XRHID, UUID string, params *api_public.RegisterDomainParams, body *api_public.Domain) (string, *header.XRHIDMVersion, *model.Domain, error)
	Update(xrhid *identity.XRHID, UUID string, params *api_public.UpdateDomainParams, body *api_public.Domain) (string, *header.XRHIDMVersion, *model.Domain, error)
}
