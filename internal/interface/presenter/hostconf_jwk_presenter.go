package presenter

import "github.com/podengo-project/idmsvc-backend/internal/api/public"

type HostconfJwkPresenter interface {
	PublicSigningKeys(keys []string) (*public.SigningKeysResponse, error)
}
