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

func (i hostInteractor) HostConf(xrhid *identity.XRHID, fqdn string, params *api_public.HostConfParams, body *api_public.HostConf) (*interactor.HostConfOptions, error) {
	return nil, fmt.Errorf("TODO: not implemented")
}
