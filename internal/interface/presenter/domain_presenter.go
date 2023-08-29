package presenter

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
)

type DomainPresenter interface {
	Create(domain *model.Domain) (*public.Domain, error)
	List(prefix string, count int64, offset int, limit int, data []model.Domain) (*public.ListDomainsResponse, error)
	Get(domain *model.Domain) (*public.Domain, error)
	// PartialUpdate(domain *model.Todo) (*public.UpdateDomainResponse, error)
	// FullUpdate(domain *model.Todo) (*public.UpdateDomainResponse, error)
	Register(domain *model.Domain) (*public.Domain, error)
	Update(domain *model.Domain) (*public.Domain, error)
	CreateDomainToken(token *repository.DomainRegToken) (*public.DomainRegToken, error)
}
