package interactor

import (
	"fmt"
	"time"

	"github.com/hmsidm/internal/api/header"
	"github.com/hmsidm/internal/api/public"
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/redhatinsights/platform-go-middlewares/identity"
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
	xrhid, err := header.DecodeXRHID(string(params.XRhIdentity))
	if err != nil {
		return "", nil, err
	}
	domain.OrgId = xrhid.Identity.OrgID
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
	return xrhid.Identity.OrgID, domain, nil
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
	xrhid, err := header.DecodeXRHID(string(params.XRhIdentity))
	if err != nil {
		return "", "", err
	}
	return xrhid.Identity.OrgID, uuid, nil
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
	xrhid, err := header.DecodeXRHID(string(params.XRhIdentity))
	if err != nil {
		return "", -1, -1, err
	}
	return xrhid.Identity.OrgID, offset, limit, nil
}

// TODO Document method
func (i domainInteractor) GetById(uuid string, params *public.ReadDomainParams) (string, string, error) {
	if uuid == "" {
		return "", "", fmt.Errorf("'in' cannot be an empty string")
	}

	xrhid, err := header.DecodeXRHID(string(params.XRhIdentity))
	if err != nil {
		return "", "", err
	}

	return xrhid.Identity.OrgID, uuid, nil
}

// RegisterIpa translates the API input format into the business
// data models for the PUT /domains/{uuid}/ipa/register endpoint.
// params contains the header parameters.
// body contains the input payload.
// Return the orgId and the business model for Ipa information,
// when success translation, else it returns empty string for orgId,
// nil for the Ipa data, and an error filled.
func (i domainInteractor) RegisterIpa(xrhid *identity.XRHID, params *api_public.RegisterIpaDomainParams, body *public.RegisterIpaDomainJSONRequestBody) (string, *model.Ipa, error) {
	if xrhid == nil {
		return "", nil, fmt.Errorf("'xrhid' cannot be nil")
	}
	if params == nil {
		return "", nil, fmt.Errorf("'params' cannot be nil")
	}
	if body == nil {
		return "", nil, fmt.Errorf("'body' cannot be nil")
	}
	orgId := xrhid.Identity.Internal.OrgID

	// Read the body payload
	domainIpa := &model.Ipa{}
	// TODO Update api schema for this endpoint to don't include RealmName

	// Translate realm domains
	i.registerIpaRealmDomains(body, domainIpa)

	// Certificate list
	i.registerIpaCaCerts(body, domainIpa)

	// Server list
	i.registerIpaServers(body, domainIpa)

	return orgId, domainIpa, nil
}

func (i domainInteractor) registerIpaRealmDomains(body *public.RegisterIpaDomainJSONRequestBody, domainIpa *model.Ipa) {
	if body.RealmDomains == nil {
		domainIpa.RealmDomains = pq.StringArray{}
		return
	}
	domainIpa.RealmDomains = make(pq.StringArray, 0)
	domainIpa.RealmDomains = append(
		domainIpa.RealmDomains,
		body.RealmDomains...,
	)
}

func (i domainInteractor) registerIpaCaCerts(body *public.RegisterIpaDomainJSONRequestBody, domainIpa *model.Ipa) {
	if body.CaCerts == nil {
		domainIpa.CaCerts = []model.IpaCert{}
		return
	}
	domainIpa.CaCerts = make([]model.IpaCert, len(body.CaCerts))
	for idx := range body.CaCerts {
		i.registerIpaCaCertOne(&domainIpa.CaCerts[idx], &body.CaCerts[idx])
	}
}

func (i domainInteractor) registerIpaCaCertOne(caCert *model.IpaCert, cert *api_public.CreateDomainIpaCert) {
	if cert.Nickname != nil {
		caCert.Nickname = *cert.Nickname
	}
	if cert.Issuer != nil {
		caCert.Issuer = *cert.Issuer
	}
	if cert.Subject != nil {
		caCert.Subject = *cert.Subject
	}
	if cert.SerialNumber != nil {
		caCert.SerialNumber = *cert.SerialNumber
	}
	if cert.NotValidBefore != nil {
		caCert.NotValidBefore = *cert.NotValidBefore
	}
	if cert.NotValidAfter != nil {
		caCert.NotValidAfter = *cert.NotValidAfter
	}
	if cert.Pem != nil {
		caCert.Pem = *cert.Pem
	}
}

func (i domainInteractor) registerIpaServers(body *public.RegisterIpaDomainJSONRequestBody, domainIpa *model.Ipa) {
	if body.Servers == nil {
		domainIpa.Servers = []model.IpaServer{}
		return
	}
	// FIXME Set body.Servers as required into the openapi specification
	domainIpa.Servers = make([]model.IpaServer, len(*body.Servers))
	for idx, server := range *body.Servers {
		domainIpa.Servers[idx].FQDN = server.Fqdn
		domainIpa.Servers[idx].RHSMId = server.RhsmId
		domainIpa.Servers[idx].PKInitServer = server.PkinitServer
		domainIpa.Servers[idx].CaServer = server.CaServer
		domainIpa.Servers[idx].HCCEnrollmentServer = server.HccEnrollmentServer
		domainIpa.Servers[idx].HCCUpdateServer = server.HccUpdateServer
	}
}
