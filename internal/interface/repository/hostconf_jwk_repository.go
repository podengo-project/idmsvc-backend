package repository

import (
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
	"gorm.io/gorm"
)

type HostconfJwkRepository interface {
	InsertJWK(db *gorm.DB, hcjwk *model.HostconfJwk) (err error)
	RevokeJWK(db *gorm.DB, kid string) (hcjwk *model.HostconfJwk, err error)
	ListJWKs(db *gorm.DB) (hcjwks []model.HostconfJwk, err error)
	PurgeExpiredJWKs(db *gorm.DB) (hcjwks []model.HostconfJwk, err error)
	GetPublicKeyArray(db *gorm.DB) (pubkeys []string, revokedKids []string, err error)
	// TODO: refactor code to use jwk.Set instead of []jwk.Key
	GetPrivateSigningKeys(db *gorm.DB) (privkeys []jwk.Key, err error)
}
