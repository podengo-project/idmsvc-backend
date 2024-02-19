package presenter

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/interface/presenter"
)

type hostconfJwkPresenter struct {
	cfg *config.Config
}

func NewHostconfJwkPresenter(cfg *config.Config) presenter.HostconfJwkPresenter {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	return &hostconfJwkPresenter{cfg}
}

func (p *hostconfJwkPresenter) PublicSigningKeys(keys []string, revokedKids []string) (*public.SigningKeysResponse, error) {
	if keys == nil {
		return nil, internal_errors.NilArgError("keys")
	}

	response := &public.SigningKeysResponse{
		Keys:        keys,
		RevokedKids: nil,
	}
	if len(revokedKids) > 0 {
		response.RevokedKids = &revokedKids
	}

	return response, nil
}
