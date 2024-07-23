package repository

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
)

// HostRepository interface
type HostRepository interface {
	MatchDomain(ctx context.Context, options *interactor.HostConfOptions) (output *model.Domain, err error)
	// TODO: hack, actual implementation will take gorm.DB argument
	SignHostConfToken(ctx context.Context, privs []jwk.Key, options *interactor.HostConfOptions, domain *model.Domain) (hctoken public.HostToken, err error)
}
