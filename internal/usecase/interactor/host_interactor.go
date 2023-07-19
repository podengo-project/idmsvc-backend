package interactor

import (
	"fmt"

	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type hostInteractor struct{}

func NewHostInteractor() interactor.HostInteractor {
	return hostInteractor{}
}

func (i hostInteractor) getXRHIDParams(xrhid *identity.XRHID) (orgID string, cn string, err error) {
	if xrhid == nil {
		return "", "", fmt.Errorf("'xrhid' is nil")
	}
	if xrhid.Identity.Type != "System" {
		return "", "", fmt.Errorf("invalid 'xrhid' type '%s'", xrhid.Identity.Type)
	}
	return xrhid.Identity.OrgID, xrhid.Identity.System.CommonName, nil
}

func (i hostInteractor) HostConf(xrhid *identity.XRHID, fqdn string, params *api_public.HostConfParams, body *api_public.HostConf) (*interactor.HostConfOptions, error) {
	orgID, cn, err := i.getXRHIDParams(xrhid)
	if err != nil {
		return nil, err
	}
	if fqdn == "" {
		return nil, fmt.Errorf("'fqdn' is empty")
	}
	if params == nil {
		return nil, fmt.Errorf("'params' is nil")
	}
	if body == nil {
		return nil, fmt.Errorf("'body' is nil")
	}

	options := &interactor.HostConfOptions{
		OrgId:      orgID,
		CommonName: cn,
		Fqdn:       fqdn,
		DomainId:   body.DomainId,
		DomainName: body.DomainName,
		DomainType: body.DomainType,
	}
	return options, nil
}
