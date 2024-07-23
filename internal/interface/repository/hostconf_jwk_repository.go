package repository

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
)

type HostconfJwkRepository interface {
	InsertJWK(ctx context.Context, hcjwk *model.HostconfJwk) (err error)
	RevokeJWK(ctx context.Context, kid string) (hcjwk *model.HostconfJwk, err error)
	ListJWKs(ctx context.Context) (hcjwks []model.HostconfJwk, err error)
	PurgeExpiredJWKs(ctx context.Context) (hcjwks []model.HostconfJwk, err error)
	GetPublicKeyArray(ctx context.Context) (pubkeys, revokedKids []string, err error)
	// TODO: refactor code to use jwk.Set instead of []jwk.Key
	GetPrivateSigningKeys(ctx context.Context) (privkeys []jwk.Key, err error)
}
