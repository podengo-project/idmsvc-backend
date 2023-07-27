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

func (i hostInteractor) HostConf(xrhid *identity.XRHID, inventoryId api_public.HostId, fqdn string, params *api_public.HostConfParams, body *api_public.HostConf) (*interactor.HostConfOptions, error) {
	if xrhid == nil {
		return nil, fmt.Errorf("'xrhid' is nil")
	}
	if xrhid.Identity.Type != "System" {
		return nil, fmt.Errorf("invalid 'xrhid' type '%s'", xrhid.Identity.Type)
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
		OrgId:       xrhid.Identity.OrgID,
		CommonName:  xrhid.Identity.System.CommonName,
		InventoryId: inventoryId,
		Fqdn:        fqdn,
		DomainId:    body.DomainId,
		DomainName:  body.DomainName,
		DomainType:  body.DomainType,
	}
	return options, nil
}
