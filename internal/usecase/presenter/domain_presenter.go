package presenter

// TODO Too much code duplication
// TODO Investigate if some "inheritence" mechanism in
//      opanapi specification generate common structures
//      letting to reduce the boilerplate when transforming
//      internal types <--> api types

import (
	"fmt"
	"strings"

	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/presenter"
	"github.com/openlyinc/pointy"
)

type domainPresenter struct{}

// TODO Document the function
func NewDomainPresenter() presenter.DomainPresenter {
	return domainPresenter{}
}

// TODO Document the method
func (p domainPresenter) Create(domain *model.Domain) (*public.CreateDomainResponse, error) {
	if domain == nil {
		return nil, fmt.Errorf("domain cannot be nil")
	}
	output := &public.CreateDomainResponse{}
	// TODO Maybe some nil values should be considered as a no valid response?
	// TODO Important to be consistent, whatever is the response
	output.DomainUuid = domain.DomainUuid.String()

	if domain.AutoEnrollmentEnabled == nil {
		return nil, fmt.Errorf("AutoenrollmentEnabled cannot be nil")
	}
	output.AutoEnrollmentEnabled = *domain.AutoEnrollmentEnabled

	if domain.DomainName == nil {
		return nil, fmt.Errorf("DomainName cannot be nil")
	}
	output.DomainName = *domain.DomainName

	if domain.DomainType == nil {
		return nil, fmt.Errorf("DomainType cannot be nil")
	}
	output.DomainType = public.CreateDomainResponseDomainType(model.DomainTypeString(*domain.DomainType))

	if domain.IpaDomain == nil {
		return nil, fmt.Errorf("IpaDomain cannot be nil")
	}
	if domain.IpaDomain.CaList == nil {
		return nil, fmt.Errorf("CaList cannot be nil")
	}
	output.Ipa.CaList = *domain.IpaDomain.CaList

	if domain.IpaDomain.RealmName == nil {
		return nil, fmt.Errorf("RealmName cannot be nil")
	}
	output.Ipa.RealmName = pointy.String(*domain.IpaDomain.RealmName)

	if domain.IpaDomain.ServerList == nil {
		return nil, fmt.Errorf("ServerList cannot be nil")
	}
	output.Ipa.ServerList = &[]string{}
	*output.Ipa.ServerList = strings.Split(*domain.IpaDomain.ServerList, ",")

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
		output.Data[idx].DomainName = pointy.String(*item.DomainName)
		output.Data[idx].DomainType = pointy.String(model.DomainTypeString(*item.DomainType))
		output.Data[idx].DomainUuid = pointy.String(item.DomainUuid.String())
	}
	return output, nil
}

// TODO Document the method
func (p domainPresenter) Get(domain *model.Domain) (*public.ReadDomainResponse, error) {
	if domain == nil {
		return nil, fmt.Errorf("'domain' cannot be nil")
	}
	output := &public.ReadDomainResponse{}
	// TODO Maybe some nil values should be considered as a no valid response?
	// TODO Important to be consistent, whatever is the response
	output.DomainUuid = pointy.String(domain.DomainUuid.String())

	if domain.AutoEnrollmentEnabled == nil {
		return nil, fmt.Errorf("AutoenrollmentEnabled cannot be nil")
	}
	output.AutoEnrollmentEnabled = pointy.Bool(*domain.AutoEnrollmentEnabled)

	if domain.DomainName == nil {
		return nil, fmt.Errorf("DomainName cannot be nil")
	}
	output.DomainName = pointy.String(*domain.DomainName)

	if domain.DomainType == nil {
		return nil, fmt.Errorf("DomainType cannot be nil")
	}
	output.DomainType = pointy.String(model.DomainTypeString(*domain.DomainType))

	if domain.IpaDomain == nil {
		return nil, fmt.Errorf("IpaDomain cannot be nil")
	}
	if domain.IpaDomain.CaList == nil {
		return nil, fmt.Errorf("CaList cannot be nil")
	}
	output.Ipa = &public.ReadDomainIpa{}
	output.Ipa.CaList = *domain.IpaDomain.CaList

	if domain.IpaDomain.RealmName == nil {
		return nil, fmt.Errorf("RealmName cannot be nil")
	}
	output.Ipa.RealmName = pointy.String(*domain.IpaDomain.RealmName)

	if domain.IpaDomain.ServerList == nil {
		return nil, fmt.Errorf("ServerList cannot be nil")
	}
	output.Ipa.ServerList = &[]string{}
	*output.Ipa.ServerList = strings.Split(*domain.IpaDomain.ServerList, ",")

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
