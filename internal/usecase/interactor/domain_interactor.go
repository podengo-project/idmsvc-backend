package interactor

import (
	"fmt"
	"time"

	"github.com/hmsidm/internal/api/public"
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/openlyinc/pointy"
)

type domainInteractor struct{}

// NewDomainInteractor Create an interactor for the /domain endpoint handler
// Return an initialized instance of interactor.DomainInteractor
func NewDomainInteractor() interactor.DomainInteractor {
	return domainInteractor{}
}

// helperDomainTypeToUint transform public.CreateDomainDomainType to an uint const
// Return the uint representation or model.DomainTypeUndefined if it does not match
// the current types.
func helperDomainTypeToUint(domainType public.CreateDomainDomainType) uint {
	switch domainType {
	case public.CreateDomainDomainTypeIpa:
		return model.DomainTypeIpa
	default:
		return model.DomainTypeUndefined
	}
}

func (i domainInteractor) FillCert(to *model.IpaCert, from *api_public.CreateDomainIpaCert) error {
	if to == nil {
		return fmt.Errorf("'to' cannot be nil")
	}
	if from == nil {
		return fmt.Errorf("'from' cannot be nil")
	}
	if from.Nickname == nil {
		to.Nickname = ""
	} else {
		to.Nickname = *from.Nickname
	}
	if from.Issuer == nil {
		to.Issuer = ""
	} else {
		to.Issuer = *from.Issuer
	}
	if from.Subject == nil {
		to.Subject = ""
	} else {
		to.Subject = *from.Subject
	}
	if from.NotValidAfter == nil {
		to.NotValidAfter = time.Time{}
	} else {
		to.NotValidAfter = *from.NotValidAfter
	}
	if from.NotValidBefore == nil {
		to.NotValidBefore = time.Time{}
	} else {
		to.NotValidBefore = *from.NotValidBefore
	}
	if from.Pem == nil {
		to.Pem = ""
	} else {
		to.Pem = *from.Pem
	}
	if from.SerialNumber == nil {
		to.SerialNumber = ""
	} else {
		to.SerialNumber = *from.SerialNumber
	}
	return nil
}

func (i domainInteractor) FillServer(to *model.IpaServer, from *api_public.CreateDomainIpaServer) error {
	if to == nil {
		return fmt.Errorf("'to' cannot be nil")
	}
	if from == nil {
		return fmt.Errorf("'from' cannot be nil")
	}
	to.FQDN = from.Fqdn
	to.CaServer = from.CaServer
	to.HCCEnrollmentServer = from.HccEnrollmentServer
	to.HCCUpdateServer = from.HccUpdateServer
	to.PKInitServer = from.PkinitServer
	to.RHSMId = from.RhsmId
	return nil
}

// Create translate api request to modle.Domain internal representation.
// params is the x-rh-identity as base64 and the x-rh-insights-request-id value for this request.
// body is the CreateDomain schema received for the POST request.
// Return a reference to a new model.Domain instance and nil error for
// a success transformation, else nil and an error instance.
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
	domain.DomainDescription = pointy.String(body.DomainDescription)
	domain.DomainType = pointy.Uint(helperDomainTypeToUint(body.DomainType))

	domain.IpaDomain = &model.Ipa{}
	if body.Ipa.RealmName != "" {
		domain.IpaDomain.RealmName = pointy.String(body.Ipa.RealmName)
	}
	if body.Ipa.Servers != nil {
		domain.IpaDomain.Servers = make([]model.IpaServer, len(*body.Ipa.Servers))
		for idx, server := range *body.Ipa.Servers {
			i.FillServer(&domain.IpaDomain.Servers[idx], &server)
		}
	} else {
		domain.IpaDomain.Servers = []model.IpaServer{}
	}
	if body.Ipa.CaCerts != nil {
		domain.IpaDomain.CaCerts = make([]model.IpaCert, len(body.Ipa.CaCerts))
		for idx, cert := range body.Ipa.CaCerts {
			i.FillCert(&domain.IpaDomain.CaCerts[idx], &cert)
		}
	} else {
		domain.IpaDomain.CaCerts = []model.IpaCert{}
	}
	if body.Ipa.RealmDomains == nil {
		domain.IpaDomain.RealmDomains = []string{}
	} else {
		domain.IpaDomain.RealmDomains = body.Ipa.RealmDomains
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
