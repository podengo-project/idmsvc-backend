package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
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

// InsertJWK generates and inserts a new JWK into the database. The private
// key is encrypted with the current app secret.
func (r *hostconfJwkRepository) InsertJWK(ctx context.Context, hcjwk *model.HostconfJwk) (err error) {
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if db == nil {
		err = internal_errors.NilArgError("db")
		log.Error(err.Error())
		return err
	}

	if err = db.Create(&hcjwk).Error; err != nil {
		log.Error("creating JWK")
		return err
	}
	return nil
}

// RevokeJWK revokes a JWK with key identifier `kid`
// ctx is the current request context with db and slog instances.
// kid is the key id to revoke.
// Return
func (r *hostconfJwkRepository) RevokeJWK(ctx context.Context, kid string) (hcjwk *model.HostconfJwk, err error) {
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if db == nil {
		err = internal_errors.NilArgError("db")
		log.Error(err.Error())
		return nil, err
	}
	// find JKW by unique kid
	if err = db.
		Where("key_id = ?", kid).
		First(&hcjwk).
		Error; err != nil {
		log.Error("revoking JWK when finding the JWK to revoke")
		return nil, err
	}

	if err = hcjwk.Revoke(); err != nil {
		log.Error("revoking JWK when clean-up the current in memory data")
		return nil, err
	}

	if err = db.Save(hcjwk).Error; err != nil {
		log.Error("revoking JWK when saving the data")
		return nil, err
	}

	return hcjwk, nil
}

// ListJWKs all JWKs, including expired and revoked JWKs
// ctx is the current request context with db and slog instances.
func (r *hostconfJwkRepository) ListJWKs(ctx context.Context) (hcjwks []model.HostconfJwk, err error) {
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if db == nil {
		err = internal_errors.NilArgError("db")
		log.Error(err.Error())
		return nil, err
	}
	// list JWKs, order by id to have determinstic sorting
	if err = db.
		Order("id").
		Find(&hcjwks).
		Error; err != nil {
		log.Error("listing JWK when finding records")
		return nil, err
	}
	return hcjwks, nil
}

// PurgeExpiredJWKs find and removes all JWKs that are expired
// ctx is the current request context with db and slog instances.
func (r *hostconfJwkRepository) PurgeExpiredJWKs(ctx context.Context) (hcjwks []model.HostconfJwk, err error) {
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if db == nil {
		err = internal_errors.NilArgError("db")
		log.Error(err.Error())
		return nil, err
	}

	// list JWKs, order by id to have determinstic sorting
	now := time.Now()
	if err = db.
		Order("id").
		Where("expires_at <= ?", now). // use SQL NOW()?
		Find(&hcjwks).
		Error; err != nil {
		log.Error("purging expired JWK when finding set of records")
		return nil, err
	}
	if len(hcjwks) > 0 {
		if err = db.
			Unscoped(). // do not use GORM's soft delete for purging
			Delete(&hcjwks).
			Error; err != nil {
			log.Error("purging expired JWK when deleting records")
			return nil, err
		}
	}
	return hcjwks, nil
}

// GetPublicKeyArray returns an array of string with all valid, non-expired
// public JWKs as serialized JSON. Expired or invalid keys are ignored
// ctx is the current request context with db and slog instances.
func (r *hostconfJwkRepository) GetPublicKeyArray(ctx context.Context) (pubkeys, revokedKids []string, err error) {
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if db == nil {
		err := internal_errors.NilArgError("db")
		log.Error(err.Error())
		return nil, nil, err
	}
	var hcjwks []model.HostconfJwk

	now := time.Now()
	if err = db.
		Where("expires_at > ?", now). // use SQL NOW()?
		Order("id").
		Find(&hcjwks).Error; err != nil {
		log.Error("reading keys when finding not expired keys")
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
// ctx is the current request context with db and slog instances.
func (r *hostconfJwkRepository) GetPrivateSigningKeys(ctx context.Context) (privkeys []jwk.Key, err error) {
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if db == nil {
		err := internal_errors.NilArgError("db")
		log.Error(err.Error())
		return nil, err
	}
	var hcjwks []model.HostconfJwk
	now := time.Now()
	if err = db.
		Where("encrypted_jwk is not NULL").
		Where("encryption_id = ?", r.config.Secrets.HostconfEncryptionId).
		Where("expires_at > ?", now). // use SQL NOW()?
		Order("id").
		Find(&hcjwks).Error; err != nil {
		log.Error("reading private signing key when finding records")
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
