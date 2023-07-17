package interactor

import (
	"github.com/google/uuid"
	"github.com/hmsidm/internal/api/header"
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type HostConfOptions struct {
	OrgId         string
	CommonName    string
	Fqdn          string
	ClientVersion *header.XRHIDMVersion
	DomainId      *uuid.UUID
	DomainName    *string
	DomainType    *string
}

type HostInteractor interface {
	HostConf(xrhid *identity.XRHID, fqdn string, params *api_public.HostConfParams, body *api_public.HostConf) (*HostConfOptions, error)
}
