package interactor

import (
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
)

type DomainInteractor interface {
	Create(params *api_public.CreateDomainParams, body *api_public.CreateDomain, out *model.Domain) error
	Delete(id string, params *api_public.DeleteDomainParams, out *string) error
	List(params *api_public.ListDomainsParams, offset *int64, limit *int32) error
	GetById(params string, out *string) error
}
