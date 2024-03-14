package repository

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"gorm.io/gorm"
)

var notImplementedError = fmt.Errorf("TODO: not implemented")

type hostconfJwkRepository struct {
	config *config.Config
}

func NewHostconfJwkRepository(cfg *config.Config) repository.HostconfJwkRepository {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	r := &hostconfJwkRepository{config: cfg}
	return r
}

// CreateJWK generates and inserts a new JWK into the database. The private
// key is encrypted with the current app secret.
func (r *hostconfJwkRepository) InsertJWK(db *gorm.DB, hcjwk *model.HostconfJwk) (err error) {
	if db == nil {
		return internal_errors.NilArgError("db")
	}

	if err = db.Create(&hcjwk).Error; err != nil {
		return err
	}
	return nil
}

// RevokeJWK revokes a JWK with key identifier `kid`
func (r *hostconfJwkRepository) RevokeJWK(db *gorm.DB, kid string) (hcjwk *model.HostconfJwk, err error) {
	if db == nil {
		return nil, internal_errors.NilArgError("db")
	}
	// find JKW by unique kid
	if err = db.
		Where("key_id = ?", kid).
		First(&hcjwk).
		Error; err != nil {
		return nil, err
	}

	hcjwk.Revoke()
	if err = db.Save(hcjwk).Error; err != nil {
		return nil, err
	}

	return hcjwk, nil
}

// ListJWKs all JWKs, including expired and revoked JWKs
func (r *hostconfJwkRepository) ListJWKs(db *gorm.DB) (hcjwks []model.HostconfJwk, err error) {
	if db == nil {
		return nil, internal_errors.NilArgError("db")
	}
	// list JWKs, order by id to have determinstic sorting
	if err = db.
		Order("id").
		Find(&hcjwks).
		Error; err != nil {
		return nil, err
	}
	return hcjwks, nil
}

// PurgeExpiredJWKs find and removes all JWKs that are expired
func (r *hostconfJwkRepository) PurgeExpiredJWKs(db *gorm.DB) (hcjwks []model.HostconfJwk, err error) {
	if db == nil {
		return nil, internal_errors.NilArgError("db")
	}

	// list JWKs, order by id to have determinstic sorting
	now := time.Now()
	if err = db.
		Order("id").
		Where("expires_at <= ?", now). // use SQL NOW()?
		Find(&hcjwks).
		Error; err != nil {
		return nil, err
	}
	if err = db.
		Unscoped(). // do not use GORM's soft delete for purging
		Delete(&hcjwks).
		Error; err != nil {
		return nil, err
	}
	return hcjwks, nil
}

// GetPublicKeyArray returns an array of string with all valid, non-expired
// public JWKs as serialized JSON. Expired or invalid keys are ignored
func (r *hostconfJwkRepository) GetPublicKeyArray(db *gorm.DB) (pubkeys []string, revokedKids []string, err error) {
	if db == nil {
		return nil, nil, internal_errors.NilArgError("db")
	}
	var hcjwks []model.HostconfJwk

	now := time.Now()
	if err = db.
		Where("expires_at > ?", now). // use SQL NOW()?
		Order("id").
		Find(&hcjwks).Error; err != nil {
		return nil, nil, err
	}

	for _, hcjwk := range hcjwks {
		state, _ := hcjwk.GetPublicKeyState()
		switch state {
		case hostconf_jwk.ValidKey:
			pubkeys = append(pubkeys, hcjwk.PublicJwk)
		case hostconf_jwk.RevokedKey:
			revokedKids = append(revokedKids, hcjwk.KeyId)
		}
	}
	return pubkeys, revokedKids, nil
}

// GetPrivateSigningKeys returns a array of jwk.Keys with all valid, non-expired
// private JWKs for signing that can be decrypted with the current main app
// secret. Expired, invalid keys, and keys encrypted for a different main app
// secret are ignored.
func (r *hostconfJwkRepository) GetPrivateSigningKeys(db *gorm.DB) (privkeys []jwk.Key, err error) {
	if db == nil {
		return nil, internal_errors.NilArgError("db")
	}
	var hcjwks []model.HostconfJwk
	now := time.Now()
	if err = db.
		Where("encrypted_jwk is not NULL").
		Where("encryption_id = ?", r.config.Secrets.HostconfEncryptionId).
		Where("expires_at > ?", now). // use SQL NOW()?
		Order("id").
		Find(&hcjwks).Error; err != nil {
		return nil, err
	}
	for _, hcjwk := range hcjwks {
		privkey, _, err := hcjwk.GetPrivateJWK(r.config.Secrets)
		if err == nil {
			privkeys = append(privkeys, privkey)
		}
	}
	return privkeys, nil
}
