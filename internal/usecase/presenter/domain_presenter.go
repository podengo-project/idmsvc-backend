package presenter

// TODO Too much code duplication
// TODO Investigate if some "inheritence" mechanism in
//      opanapi specification generate common structures
//      letting to reduce the boilerplate when transforming
//      internal types <--> api types

import (
	"fmt"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/presenter"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
)

type domainPresenter struct {
	cfg *config.Config
}

// NewDomainPresenter create a new DomainPresenter instance
// Return a new presenter.DomainPresenter instance
func NewDomainPresenter(cfg *config.Config) presenter.DomainPresenter {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	return &domainPresenter{cfg}
}

// Create translate from internal domain to the API response.
// Return a new response domain representation and nil error on success,
// or a nil response with an error on failure.
func (p *domainPresenter) Create(domain *model.Domain) (*public.Domain, error) {
	return p.sharedDomain(domain)
}

// List is the output adapter to list the domains with pagination.
// prefix is the prefix that is used to compose the pagination links.
// offset is the starting point of the page for a given ordered list of domains.
// count is the number of items on the current page.
// data is the slice with the model.Domain
func (p *domainPresenter) List(prefix string, count int64, offset int, limit int, data []model.Domain) (*public.ListDomainsResponse, error) {
	// https://consoledot.pages.redhat.com/docs/dev/developer-references/rest/pagination.html
	if offset < 0 {
		return nil, fmt.Errorf("'offset' is lower than 0")
	}
	if limit < 0 {
		return nil, fmt.Errorf("'limit' is lower than 0")
	}
	if limit == 0 {
		limit = p.cfg.Application.PaginationDefaultLimit
	}
	if limit > p.cfg.Application.PaginationMaxLimit {
		limit = p.cfg.Application.PaginationMaxLimit
	}
	output := &public.ListDomainsResponse{}
	p.listFillMeta(output, count, offset, limit)
	p.listFillLinks(output, prefix, count, offset, limit)

	sizeData := limit
	if len(data) < limit {
		sizeData = len(data)
	}
	output.Data = make([]public.ListDomainsData, sizeData)
	for idx := range data {
		if idx == sizeData {
			break
		}

		p.listFillItem(&output.Data[idx], &data[idx])
	}
	return output, nil
}

// TODO Document the method
func (p *domainPresenter) Get(domain *model.Domain) (*public.Domain, error) {
	return p.sharedDomain(domain)
}

// Register translate model.Domain instance to Domain output
// representation for the API response.
// domain Not nil reference to the domain model.
// Return a reference to a pubic.Domain and nil error for
// a success translation, else nil and an error with the details.
func (p *domainPresenter) Register(
	domain *model.Domain,
) (output *public.Domain, err error) {
	return p.sharedDomain(domain)
}

// Update translate model.Domain instance to Domain output
// representation for the API response.
// domain Not nil reference to the domain model.
// Return a reference to a pubic.Domain and nil error for
// a success translation, else nil and an error with the details.
func (p *domainPresenter) Update(
	domain *model.Domain,
) (output *public.Domain, err error) {
	return p.sharedDomain(domain)
}

// func (p *domainPresenter) PartialUpdate(
// 	domain *model.Domain,
// ) (output *public.Domain, err error) {
// 	return p.sharedDomain(domain)
// }

// func (p *domainPresenter) FullUpdate(
// 	domain *model.Domain,
// ) (output *public.Domain, err error) {
// 	return p.sharedDomain(domain)
// }

// Create domain registration token
// Translate the internal token represenatation to public API
func (p *domainPresenter) CreateDomainToken(token *repository.DomainRegToken) (*public.DomainRegToken, error) {
	drt := &public.DomainRegToken{
		DomainId:    token.DomainId,
		DomainToken: token.DomainToken,
		DomainType:  token.DomainType,
		Expiration:  int(token.ExpirationNS / 1_000_000_000),
	}
	return drt, nil
}
