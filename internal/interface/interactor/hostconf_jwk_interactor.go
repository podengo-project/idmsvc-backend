package interactor

import (
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

type HostconfJwkInteractor interface {
	GetSigningKeys(rhid *identity.XRHID, params *api_public.GetSigningKeysParams) (orgID string, err error)
}
