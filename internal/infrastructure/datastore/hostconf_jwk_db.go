package datastore

import (
	"log/slog"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
	interface_repository "github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"github.com/podengo-project/idmsvc-backend/internal/usecase/repository"
	"gorm.io/gorm"
)

type HostconfJwkDb struct {
	cfg        *config.Config
	repository interface_repository.HostconfJwkRepository
}

// Create new HostconfJwkDb
func NewHostconfJwkDb(cfg *config.Config) *HostconfJwkDb {
	return &HostconfJwkDb{
		cfg:        cfg,
		repository: repository.NewHostconfJwkRepository(cfg),
	}
}

// Calculate and log renew after and expires at times
func (r *HostconfJwkDb) timestamps() (renewAfter, expiresAfter time.Time) {
	utcnow := time.Now().UTC().Truncate(time.Second)
	renewAfter = utcnow.Add(r.cfg.Application.HostconfJwkRenewalThreshold)
	expiresAfter = utcnow.Add(r.cfg.Application.HostconfJwkValidity)
	slog.Info(
		"Hostconf JWK configuration",
		slog.Time("now", utcnow),
		slog.Time("renewAfter", renewAfter),
		slog.Duration("renewThreshold", r.cfg.Application.HostconfJwkRenewalThreshold),
		slog.Time("expiresAfter", expiresAfter),
		slog.Duration("validity", r.cfg.Application.HostconfJwkValidity),
	)
	return renewAfter, expiresAfter

}

// Refresh and create JWKs in database
// A new JWK is created whenever the database has no valid JWK or all JWKs
// expire within the renewal threshold period.
func (r *HostconfJwkDb) Refresh() (err error) {
	var (
		db     *gorm.DB
		tx     *gorm.DB
		hcjwks []model.HostconfJwk
	)
	db = NewDB(r.cfg)
	defer Close(db)

	if tx = db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if hcjwks, err = r.repository.ListJWKs(db); err != nil {
		return err
	}

	// Create new JWK unless there is one or more JWK that is not expired,
	// not revoked, encrypted with current app secret, and which does expires
	// after renewal threshold.
	create := true
	valid := 0
	revoked := 0
	expired := 0
	renewAfter, expiresAfter := r.timestamps()

	for _, hcjwk := range hcjwks {
		privstate, _ := hcjwk.GetPrivateKeyState(r.cfg.Secrets)
		log := slog.With(
			slog.String("kid", hcjwk.KeyId),
			slog.String("privatekey", hostconf_jwk.KeyStateString(privstate)),
			slog.Time("expires", hcjwk.ExpiresAt),
			slog.Time("expiresAfter", expiresAfter),
			slog.Time("renewalAfter", renewAfter),
		)
		switch privstate {
		case hostconf_jwk.ValidKey:
			valid += 1
			if hcjwk.ExpiresAt.Unix() >= renewAfter.Unix() {
				log.Info("Valid Hostconf JWK")
				create = false
			} else {
				log.Warn("Valid Hostconf JWK is after renewal threshold")
			}
		case hostconf_jwk.RevokedKey:
			log.Info("Hostconf JWK is revoked")
			revoked += 1
		case hostconf_jwk.ExpiredKey:
			log.Info("Hostconf JWK is expired")
			expired += 1
		default:
			log.Info("Hostconf JWK is invalid")
		}
	}

	slog.Info(
		"Current JWKs in database",
		slog.Int("total", len(hcjwks)),
		slog.Int("valid", valid),
		slog.Int("expired", expired),
		slog.Int("revoked", revoked),
	)

	if create {
		var newjwk *model.HostconfJwk

		if valid == 0 {
			slog.Warn("No valid JWK found in database")
		} else {
			slog.Warn("All valid JWKs expire in the renewal threshold period")
		}

		if newjwk, err = model.NewHostconfJwk(r.cfg.Secrets, expiresAfter); err != nil {
			return err
		}
		r.repository.InsertJWK(db, newjwk)
		slog.Info(
			"Created new hostconf JWK",
			slog.String("kid", newjwk.KeyId),
			slog.Time("expires", newjwk.ExpiresAt),
		)
	}

	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Mark a JWK as revoked
func (r *HostconfJwkDb) Revoke(kid string) (err error) {
	var (
		db    *gorm.DB
		tx    *gorm.DB
		hcjwk *model.HostconfJwk
	)
	db = NewDB(r.cfg)
	defer Close(db)

	if tx = db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if hcjwk, err = r.repository.RevokeJWK(tx, kid); err != nil {
		return err
	}
	slog.Info("Revoked JWK", slog.String("kid", hcjwk.KeyId))

	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Purge and remove expired JWKs from database
func (r *HostconfJwkDb) Purge() (err error) {
	var (
		db     *gorm.DB
		tx     *gorm.DB
		hcjwks []model.HostconfJwk
	)
	db = NewDB(r.cfg)
	defer Close(db)

	if tx = db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if hcjwks, err = r.repository.PurgeExpiredJWKs(tx); err != nil {
		return err
	}
	if len(hcjwks) > 0 {
		slog.Info("Purged keys from DB", slog.Int("purged", len(hcjwks)))
		for _, hcjwk := range hcjwks {
			slog.Info(
				"Purged key",
				slog.String("kid", hcjwk.KeyId),
				slog.Time("expires", hcjwk.ExpiresAt),
			)
		}
	} else {
		slog.Info("Nothing to purge")
	}

	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}
	return nil
}

// List all JWKs in database
func (r *HostconfJwkDb) ListKeys() (err error) {
	var (
		db     *gorm.DB
		tx     *gorm.DB
		hcjwks []model.HostconfJwk
		hcjwk  model.HostconfJwk
	)
	db = NewDB(r.cfg)
	defer Close(db)

	if tx = db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	renewAfter, expiresAfter := r.timestamps()

	if hcjwks, err = r.repository.ListJWKs(db); err != nil {
		return err
	}
	slog.Info("JWKs in database", slog.Int("length", len(hcjwks)))
	for _, hcjwk = range hcjwks {
		pubstate, _ := hcjwk.GetPublicKeyState()
		privstate, _ := hcjwk.GetPrivateKeyState(r.cfg.Secrets)
		slog.Info(
			"Hostconf JWK",
			slog.String("kid", hcjwk.KeyId),
			slog.String("publickey", hostconf_jwk.KeyStateString(pubstate)),
			slog.String("privatekey", hostconf_jwk.KeyStateString(privstate)),
			slog.Time("expires", hcjwk.ExpiresAt),
			slog.Time("expiresAfter", expiresAfter),
			slog.Time("renewAfter", renewAfter),
		)
	}
	return nil
}
