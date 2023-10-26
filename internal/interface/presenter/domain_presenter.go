package presenter

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
)

type DomainPresenter interface {
	List(count int64, offset int, limit int, data []model.Domain) (*public.ListDomainsResponse, error)
	Get(domain *model.Domain) (*public.Domain, error)
	// PartialUpdate(domain *model.Todo) (*public.UpdateDomainResponse, error)
	// FullUpdate(domain *model.Todo) (*public.UpdateDomainResponse, error)
	Register(domain *model.Domain) (*public.RegisterDomainResponse, error)
	UpdateAgent(domain *model.Domain) (*public.UpdateDomainAgentResponse, error)
	UpdateUser(domain *model.Domain) (*public.UpdateDomainUserResponse, error)
	CreateDomainToken(token *repository.DomainRegToken) (*public.DomainRegToken, error)
}
