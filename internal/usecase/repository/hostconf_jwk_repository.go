package repository

import (
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"gorm.io/gorm"
)

var notImplementedError = fmt.Errorf("TODO: not implemented")

type hostconfJwkRepository struct {
	config *config.Config
	// TODO: store JWKs in database
	signingKeys []jwk.Key
	publicKeys  []string
}

func NewHostconfJwkRepository(cfg *config.Config) repository.HostconfJwkRepository {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	return &hostconfJwkRepository{config: cfg}
}

// CreateJWK generates and inserts a new JWK into the database. The private
// key is encrypted with the current app secret.
func (r *hostconfJwkRepository) CreateJWK(db *gorm.DB) (model *model.HostconfJwk, err error) {
	return nil, notImplementedError
}

// RevokeJWK revokes a JWK with key identifier `kid`
func (r *hostconfJwkRepository) RevokeJWK(db *gorm.DB, kid string) (model *model.HostconfJwk, err error) {
	return nil, notImplementedError
}

// ListJWKs all JWKs, including expired and revoked JWKs
func (r *hostconfJwkRepository) ListJWKs(db *gorm.DB) (models []model.HostconfJwk, err error) {
	return nil, notImplementedError
}

// GetPublicKeyArray returns an array of string with all valid, non-expired
// public JWKs as serialized JSON. Expired or invalid keys are ignored
func (r *hostconfJwkRepository) GetPublicKeyArray(db *gorm.DB) (pubkeys []string, err error) {
	return nil, notImplementedError
}

// GetPrivateSigningKeys returns a array of jwk.Keys with all valid, non-expired
// private JWKs for signing that can be decrypted with the current main app
// secret. Expired, invalid keys, and keys encrypted for a different main app
// secret are ignored.
func (r *hostconfJwkRepository) GetPrivateSigningKeys(db *gorm.DB) (privkeys []jwk.Key, err error) {
	return nil, notImplementedError
}
