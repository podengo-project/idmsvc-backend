package interactor

import (
	"fmt"

	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/openlyinc/pointy"
)

type domainInteractor struct{}

func NewTodoInteractor() interactor.DomainInteractor {
	return domainInteractor{}
}

// TODO Document method
func (i domainInteractor) Create(params *api_public.CreateDomainParams, in *api_public.CreateDomain, out *model.Domain) (error) {
	if params == nil {
		return fmt.Errorf("'params' cannot be nil")
	}
	if in == nil {
		return fmt.Errorf("'in' cannot be nil")
	}
	if out == nil {
		return fmt.Errorf("'out' cannot be nil")
	}
	// FIXME Add a middleware that decode the X-Rh-Identity and store
	//       the structure into the request context, so we can use directly
	//       into the specific service handler and pass the information to
	//       the interactors
	identity, err := DecodeIdentity(string(params.XRhIdentity))
	if err != nil {
		return err
	}
	out.OrgId = identity.OrgID
	out.DomainName = pointy.String(in.DomainName)
	switch in.DomainType {
	case api_public.CreateDomainDomainTypeIpa:
		out.DomainType = pointy.Uint(model.DomainTypeIpa)
	}
	out.AutoEnrollmentEnabled = pointy.Bool(in.AutoEnrollmentEnabled)
	return nil
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
func (i domainInteractor) Delete(id string, params *api_public.DeleteDomainParams, out *string) error {
	if params == nil {
		return fmt.Errorf("'params' cannot be nil")
	}
	if out == nil {
		return fmt.Errorf("'out' cannot be nil")
	}
	*out = id
	return nil
}

// TODO Document method
func (i domainInteractor) List(params *api_public.ListDomainsParams, offset *int, limit *int) error {
	if params == nil {
		return fmt.Errorf("'in' cannot be nil")
	}
	if offset == nil {
		return fmt.Errorf("'offset' cannot be nil")
	}
	if limit == nil {
		return fmt.Errorf("'limit' cannot be nil")
	}
	*offset = *params.Offset
	*limit = *params.Limit
	return nil
}

// TODO Document method
func (i domainInteractor) GetById(params *api_public.ReadDomainParams, in *string) (string, string, error) {
	if in == nil {
		return "", "", fmt.Errorf("'in' cannot be nil")
	}
	params.

	return nil
}
