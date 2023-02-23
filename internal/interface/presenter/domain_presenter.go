package presenter

import (
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
)

type DomainPresenter interface {
	Create(domain *model.Domain) (*public.CreateDomainResponse, error)
	List(prefix string, offset int64, count int32, data []model.Domain) (*public.ListDomainsResponse, error)
	Get(domain *model.Domain) (*public.ReadDomainResponse, error)
	// PartialUpdate(domain *model.Todo) (*public.UpdateDomainResponse, error)
	// FullUpdate(domain *model.Todo) (*public.UpdateDomainResponse, error)
}
