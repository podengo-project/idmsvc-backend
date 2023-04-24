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
	if domain == nil {
		return nil, fmt.Errorf("'domain' is nil")
	}
	output := &public.Domain{}
	// TODO Maybe some nil values should be considered as a no valid response?
	// TODO Important to be consistent, whatever is the response
	output.DomainUuid = domain.DomainUuid.String()

	if domain.AutoEnrollmentEnabled == nil {
		return nil, fmt.Errorf("'AutoEnrollmentEnabled' is nil")
	}
	output.AutoEnrollmentEnabled = *domain.AutoEnrollmentEnabled

	if domain.DomainName == nil {
		output.DomainName = ""
	} else {
		output.DomainName = *domain.DomainName
	}

	if domain.Type == nil {
		return nil, fmt.Errorf("'Type' is nil")
	}
	output.Type = public.DomainType(model.DomainTypeString(*domain.Type))

	switch *domain.Type {
	case model.DomainTypeIpa:
		if domain.IpaDomain == nil {
			return output, nil
		}
		if err := p.createRhelIdm(output, domain); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("'Type' is invalid")
	}

	return output, nil
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
	if err := p.getChecks(domain); err != nil {
		return nil, err
	}
	output := &public.Domain{}
	// TODO Maybe some nil values should be considered as a no valid response?
	// TODO Important to be consistent, whatever is the response
	output.DomainUuid = domain.DomainUuid.String()
	output.AutoEnrollmentEnabled = *domain.AutoEnrollmentEnabled
	if domain.DomainName != nil {
		output.DomainName = *domain.DomainName
	}
	if domain.Title != nil {
		output.Title = *domain.Title
	}
	if domain.Description != nil {
		output.Description = *domain.Description
	}
	output.Type = public.DomainType(model.DomainTypeString(*domain.Type))
	switch *domain.Type {
	case model.DomainTypeIpa:
		output.RhelIdm = &public.DomainIpa{}
		if err := p.fillRhelIdmCerts(output, domain); err != nil {
			return nil, err
		}
		output.RhelIdm.RealmName = *domain.IpaDomain.RealmName
		if err := p.fillRhelIdmServers(output, domain); err != nil {
			return nil, err
		}
		if domain.IpaDomain.RealmDomains == nil {
			output.RhelIdm.RealmDomains = []string{}
		} else {
			output.RhelIdm.RealmDomains = domain.IpaDomain.RealmDomains
		}
	default:
		return nil, fmt.Errorf("'Type' is invalid")
	}

	return output, nil
}

// Register translate model.Ipa instance to DomainResponseIpa
// representation for the API response.
// domain Not nil reference to the model.Ipa to represent.
// Return a reference to a DomainResponseIpa and nil error for
// a success translation, else nil and an error with the details.
func (p domainPresenter) Register(
	domain *model.Domain,
) (output *public.Domain, err error) {
	if domain == nil {
		return nil, fmt.Errorf("'domain' is nil")
	}
	if domain.Type == nil || *domain.Type == model.DomainTypeUndefined {
		return nil, fmt.Errorf("'domain.Type' is invalid")
	}
	output = &public.Domain{}
	p.registerFillDomainData(domain, output)
	switch *domain.Type {
	case model.DomainTypeIpa:
		output.Type = model.DomainTypeIpaString
		output.RhelIdm = &public.DomainIpa{}
		err = p.registerIpa(domain, output)
	default:
		err = fmt.Errorf("'domain.Type=%d' is unsupported", *domain.Type)
	}
	if err != nil {
		return nil, err
	}
	return output, nil
}

// func (p todoPresenter) PartialUpdate(in *model.Todo, out *public.Todo) error {
// 	if in == nil {
// 		return fmt.Errorf("'in' cannot be nil")
// 	}
// 	if out == nil {
// 		return fmt.Errorf("'out' cannot be nil")
// 	}
// 	if in.ID == 0 {
// 		out.Id = nil
// 	} else {
// 		out.Id = pointy.Uint(in.ID)
// 	}
// 	if in.Title == nil {
// 		out.Title = nil
// 	} else {
// 		out.Title = pointy.String(*in.Title)
// 	}
// 	if in.Description == nil {
// 		out.Body = nil
// 	} else {
// 		out.Body = pointy.String(*in.Description)
// 	}
// 	return nil
// }

// func (p todoPresenter) FullUpdate(in *model.Todo, out *public.Todo) error {
// 	return p.PartialUpdate(in, out)
// }
