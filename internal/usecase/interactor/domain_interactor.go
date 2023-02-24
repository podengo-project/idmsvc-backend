package interactor

import (
	"fmt"
	"strings"

	"github.com/hmsidm/internal/api/public"
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/openlyinc/pointy"
)

type domainInteractor struct{}

func NewDomainInteractor() interactor.DomainInteractor {
	return domainInteractor{}
}

func helperDomainTypeToUint(domainType public.CreateDomainDomainType) uint {
	switch domainType {
	case public.CreateDomainDomainTypeIpa:
		return model.DomainTypeIpa
	default:
		return model.DomainTypeUndefined
	}
}

// TODO Document method
func (i domainInteractor) Create(params *api_public.CreateDomainParams, body *api_public.CreateDomain) (string, *model.Domain, error) {
	if params == nil {
		return "", nil, fmt.Errorf("'params' cannot be nil")
	}
	if body == nil {
		return "", nil, fmt.Errorf("'body' cannot be nil")
	}
	// FIXME Add a middleware that decode the X-Rh-Identity and store
	//       the structure into the request context, so we can use directly
	//       into the specific service handler and pass the information to
	//       the interactors

	domain := &model.Domain{}
	identity, err := DecodeIdentity(string(params.XRhIdentity))
	if err != nil {
		return "", nil, err
	}
	domain.OrgId = identity.OrgID
	domain.AutoEnrollmentEnabled = pointy.Bool(body.AutoEnrollmentEnabled)
	domain.DomainName = pointy.String(body.DomainName)
	domain.DomainType = pointy.Uint(helperDomainTypeToUint(body.DomainType))

	domain.IpaDomain = &model.Ipa{}
	if body.Ipa.RealmName != nil {
		domain.IpaDomain.RealmName = pointy.String(*body.Ipa.RealmName)
	}
	if body.Ipa.ServerList != nil {
		domain.IpaDomain.ServerList = pointy.String(strings.Join(*body.Ipa.ServerList, ","))
	}
	if body.Ipa.CaList != "" {
		domain.IpaDomain.CaList = pointy.String(body.Ipa.CaList)
	} else {
		domain.IpaDomain.CaList = pointy.String("")
	}
	return identity.OrgID, domain, nil
}

// func (i domainInteractor) PartialUpdate(id public.Id, params *api_public.PartialUpdateTodoParams, in *api_public.Todo, out *model.Todo) error {
// 	if id <= 0 {
// 		return fmt.Errorf("'id' should be a positive int64")
// 	}
// 	if params == nil {
// 		return fmt.Errorf("'params' cannot be nil")
// 	}
// 	if in == nil {
// 		return fmt.Errorf("'in' cannot be nil")
// 	}
// 	if out == nil {
// 		return fmt.Errorf("'out' cannot be nil")
// 	}
// 	out.Model.ID = id
// 	if in.Title != nil {
// 		out.Title = pointy.String(*in.Title)
// 	}
// 	if in.Body != nil {
// 		out.Description = pointy.String(*in.Body)
// 	}
// 	return nil
// }

// func (i domainInteractor) FullUpdate(id public.Id, params *api_public.UpdateTodoParams, in *api_public.Todo, out *model.Todo) error {
// 	if id <= 0 {
// 		return fmt.Errorf("'id' should be a positive int64")
// 	}
// 	if params == nil {
// 		return fmt.Errorf("'params' cannot be nil")
// 	}
// 	if in == nil {
// 		return fmt.Errorf("'in' cannot be nil")
// 	}
// 	if out == nil {
// 		return fmt.Errorf("'out' cannot be nil")
// 	}
// 	out.ID = id
// 	out.Title = pointy.String(*in.Title)
// 	out.Description = pointy.String(*in.Body)
// 	return nil
// }

// TODO Document method
func (i domainInteractor) Delete(uuid string, params *api_public.DeleteDomainParams) (string, string, error) {
	if params == nil {
		return "", "", fmt.Errorf("'params' cannot be nil")
	}
	identity, err := DecodeIdentity(string(params.XRhIdentity))
	if err != nil {
		return "", "", err
	}
	return identity.OrgID, uuid, nil
}

// TODO Document method
func (i domainInteractor) List(params *api_public.ListDomainsParams) (orgId string, offset int, limit int, err error) {
	if params == nil {
		return "", -1, -1, fmt.Errorf("'params' cannot be nil")
	}
	if params.Offset == nil {
		offset = 0
	} else {
		offset = *params.Offset
	}
	if params.Limit == nil {
		limit = 10
	} else {
		limit = *params.Limit
	}
	identity, err := DecodeIdentity(string(params.XRhIdentity))
	if err != nil {
		return "", -1, -1, err
	}
	return identity.OrgID, offset, limit, nil
}

// TODO Document method
func (i domainInteractor) GetById(uuid string, params *public.ReadDomainParams) (string, string, error) {
	if uuid == "" {
		return "", "", fmt.Errorf("'in' cannot be an empty string")
	}

	identity, err := DecodeIdentity(string(params.XRhIdentity))
	if err != nil {
		return "", "", err
	}

	return identity.OrgID, uuid, nil
}
