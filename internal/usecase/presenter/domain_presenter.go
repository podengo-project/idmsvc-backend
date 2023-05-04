package presenter

// TODO Too much code duplication
// TODO Investigate if some "inheritence" mechanism in
//      opanapi specification generate common structures
//      letting to reduce the boilerplate when transforming
//      internal types <--> api types

import (
	"fmt"

	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/presenter"
	"github.com/openlyinc/pointy"
)

type domainPresenter struct{}

// NewDomainPresenter create a new DomainPresenter instance
// Return a new presenter.DomainPresenter instance
func NewDomainPresenter() presenter.DomainPresenter {
	return domainPresenter{}
}

// Create translate from internal domain to the API response.
// Return a new response domain representation and nil error on success,
// or a nil response with an error on failure.
func (p domainPresenter) Create(domain *model.Domain) (*public.Domain, error) {
	return p.sharedDomain(domain)
}

// TODO Document the method
func (p domainPresenter) List(prefix string, offset int64, count int32, data []model.Domain) (*public.ListDomainsResponse, error) {
	output := &public.ListDomainsResponse{}
	output.Meta.Count = pointy.Int32(count)
	if offset > 0 {
		output.Links.First = pointy.String(fmt.Sprintf(prefix+"/todo?offset=%d&count=%d", 0, count))
	}
	if offset-int64(count) < 0 {
		output.Links.Previous = pointy.String(fmt.Sprintf(prefix+"/todo?offset=%d&count=%d", 0, count))
	} else {
		output.Links.Previous = pointy.String(fmt.Sprintf(prefix+"/todo?offset=%d&count=%d", offset-int64(count), count))
	}
	output.Links.Next = pointy.String(fmt.Sprintf(prefix+"/todo?offset=%d&count=%d", offset+int64(count), count))
	// TODO Add Links.Last
	// FIXME this is weird and I am not happy with this.
	//       I would need to find to modify the openapi spec to
	//       generate a structure more accessible.
	output.Data = make([]public.ListDomainsData, len(data))
	for idx, item := range data {
		output.Data[idx].AutoEnrollmentEnabled = pointy.Bool(*item.AutoEnrollmentEnabled)
		if item.DomainName != nil {
			output.Data[idx].DomainName = pointy.String(*item.DomainName)
		}
		output.Data[idx].DomainType = pointy.String(model.DomainTypeString(*item.Type))
		output.Data[idx].DomainUuid = pointy.String(item.DomainUuid.String())
	}
	return output, nil
}

// TODO Document the method
func (p domainPresenter) Get(domain *model.Domain) (*public.Domain, error) {
	return p.sharedDomain(domain)
}

// Register translate model.Domain instance to Domain output
// representation for the API response.
// domain Not nil reference to the domain model.
// Return a reference to a pubic.Domain and nil error for
// a success translation, else nil and an error with the details.
func (p domainPresenter) Register(
	domain *model.Domain,
) (output *public.Domain, err error) {
	return p.sharedDomain(domain)
}

// Update translate model.Domain instance to Domain output
// representation for the API response.
// domain Not nil reference to the domain model.
// Return a reference to a pubic.Domain and nil error for
// a success translation, else nil and an error with the details.
func (p domainPresenter) Update(
	domain *model.Domain,
) (output *public.Domain, err error) {
	return p.sharedDomain(domain)
}

// func (p domainPresenter) PartialUpdate(
// 	domain *model.Domain,
// ) (output *public.Domain, err error) {
// 	return p.sharedDomain(domain)
// }

// func (p domainPresenter) FullUpdate(
// 	domain *model.Domain,
// ) (output *public.Domain, err error) {
// 	return p.sharedDomain(domain)
// }
