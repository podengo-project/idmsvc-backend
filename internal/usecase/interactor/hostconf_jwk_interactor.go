package interactor

import (
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

type hostconfJwkInteractor struct{}

func NewHostconfJwkInteractor() interactor.HostconfJwkInteractor {
	return hostconfJwkInteractor{}
}

func (i hostconfJwkInteractor) GetSigningKeys(xrhid *identity.XRHID, params *api_public.GetSigningKeysParams) (orgID string, err error) {
	if xrhid == nil {
		return "", internal_errors.NilArgError("xrhid")
	}
	if params == nil {
		return "", internal_errors.NilArgError("params")
	}
	return xrhid.Identity.OrgID, nil
}
