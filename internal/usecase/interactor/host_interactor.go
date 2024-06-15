package interactor

import (
	"fmt"

	"github.com/google/uuid"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

type hostInteractor struct{}

func NewHostInteractor() interactor.HostInteractor {
	return hostInteractor{}
}

func (i hostInteractor) HostConf(
	xrhid *identity.XRHID,
	inventoryId api_public.HostId,
	fqdn api_public.Fqdn,
	params *api_public.HostConfParams,
	body *api_public.HostConf,
) (*interactor.HostConfOptions, error) {
	if xrhid == nil {
		return nil, internal_errors.NilArgError("xrhid")
	}
	if xrhid.Identity.Type != "System" {
		return nil, fmt.Errorf("invalid 'xrhid' type '%s'", xrhid.Identity.Type)
	}
	if fqdn == "" {
		return nil, fmt.Errorf("'fqdn' is empty")
	}
	if params == nil {
		return nil, internal_errors.NilArgError("params")
	}
	if body == nil {
		return nil, internal_errors.NilArgError("body")
	}
	cn, err := uuid.Parse(xrhid.Identity.System.CommonName)
	if err != nil {
		return nil, err
	}

	options := &interactor.HostConfOptions{
		OrgId:       xrhid.Identity.OrgID,
		CommonName:  cn,
		InventoryId: inventoryId,
		Fqdn:        fqdn,
		DomainId:    body.DomainId,
		DomainName:  body.DomainName,
		DomainType:  body.DomainType,
	}
	return options, nil
}
