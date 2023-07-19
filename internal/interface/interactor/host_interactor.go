package interactor

import (
	"github.com/google/uuid"
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type HostConfOptions struct {
	OrgId      string
	CommonName string
	Fqdn       string
	DomainId   *uuid.UUID
	DomainName *string
	DomainType *api_public.DomainType
}

type HostInteractor interface {
	HostConf(xrhid *identity.XRHID, fqdn string, params *api_public.HostConfParams, body *api_public.HostConf) (*HostConfOptions, error)
}
