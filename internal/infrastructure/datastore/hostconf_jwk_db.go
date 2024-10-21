package datastore

import (
	"context"
	"log/slog"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
	interface_repository "github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"github.com/podengo-project/idmsvc-backend/internal/usecase/repository"
	"gorm.io/gorm"
)

type HostconfJwkDb struct {
	cfg        *config.Config
	repository interface_repository.HostconfJwkRepository
	log        *slog.Logger
}

// NewHostconfJwkDb Create new HostconfJwkDb
func NewHostconfJwkDb(cfg *config.Config, log *slog.Logger) *HostconfJwkDb {
	return &HostconfJwkDb{
		cfg:        cfg,
		repository: repository.NewHostconfJwkRepository(cfg),
		log:        log,
	}
}

// Calculate and log renew after and expires at times
func (r *HostconfJwkDb) timestamps() (renewAfter, expiresAfter time.Time) {
	utcnow := time.Now().UTC().Truncate(time.Second)
	renewAfter = utcnow.Add(r.cfg.Application.HostconfJwkRenewalThreshold)
	expiresAfter = utcnow.Add(r.cfg.Application.HostconfJwkValidity)
	r.log.Info(
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

	ctx := app_context.CtxWithDB(app_context.CtxWithLog(context.Background(), r.log), tx)
	if hcjwks, err = r.repository.ListJWKs(ctx); err != nil {
		r.log.Error(err.Error())
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
		logHCJW := r.log.With(
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
				logHCJW.Info("Valid Hostconf JWK")
				create = false
			} else {
				logHCJW.Warn("Valid Hostconf JWK is after renewal threshold")
			}
		case hostconf_jwk.RevokedKey:
			logHCJW.Info("Hostconf JWK is revoked")
			revoked += 1
		case hostconf_jwk.ExpiredKey:
			logHCJW.Info("Hostconf JWK is expired")
			expired += 1
		default:
			logHCJW.Info("Hostconf JWK is invalid")
		}
	}

	r.log.Info(
		"Current JWKs in database",
		slog.Int("total", len(hcjwks)),
		slog.Int("valid", valid),
		slog.Int("expired", expired),
		slog.Int("revoked", revoked),
	)

	if create {
		var newjwk *model.HostconfJwk

		if valid == 0 {
			r.log.Warn("No valid JWK found in database")
		} else {
			r.log.Warn("All valid JWKs expire in the renewal threshold period")
		}

		if newjwk, err = model.NewHostconfJwk(r.cfg.Secrets, expiresAfter); err != nil {
			r.log.Error(err.Error())
			return err
		}
		if err = r.repository.InsertJWK(ctx, newjwk); err != nil {
			r.log.Error(err.Error())
			return err
		}
		r.log.Info(
			"Created new hostconf JWK",
			slog.String("kid", newjwk.KeyId),
			slog.Time("expires", newjwk.ExpiresAt),
		)
	}

	if tx.Commit(); tx.Error != nil {
		r.log.Error(tx.Error.Error())
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
		r.log.Error(tx.Error.Error())
		return tx.Error
	}
	defer tx.Rollback()

	ctx := app_context.CtxWithDB(app_context.CtxWithLog(context.Background(), r.log), tx)
	if hcjwk, err = r.repository.RevokeJWK(ctx, kid); err != nil {
		r.log.Error(err.Error())
		return err
	}
	r.log.Info("Revoked JWK", slog.String("kid", hcjwk.KeyId))

	if tx.Commit(); tx.Error != nil {
		r.log.Error(tx.Error.Error())
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
		r.log.Error(tx.Error.Error())
		return tx.Error
	}
	defer tx.Rollback()

	ctx := app_context.CtxWithDB(app_context.CtxWithLog(context.Background(), r.log), tx)
	if hcjwks, err = r.repository.PurgeExpiredJWKs(ctx); err != nil {
		r.log.Error(err.Error())
		return err
	}
	if len(hcjwks) > 0 {
		r.log.Info("Purged keys from DB", slog.Int("purged", len(hcjwks)))
		for _, hcjwk := range hcjwks {
			r.log.Info(
				"Purged key",
				slog.String("kid", hcjwk.KeyId),
				slog.Time("expires", hcjwk.ExpiresAt),
			)
		}
	} else {
		r.log.Info("Nothing to purge")
	}

	if tx.Commit(); tx.Error != nil {
		r.log.Error(tx.Error.Error())
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
		r.log.Error(tx.Error.Error())
		return tx.Error
	}
	defer tx.Rollback()

	renewAfter, expiresAfter := r.timestamps()

	ctx := app_context.CtxWithDB(app_context.CtxWithLog(context.Background(), r.log), tx)
	if hcjwks, err = r.repository.ListJWKs(ctx); err != nil {
		r.log.Error(err.Error())
		return err
	}
	r.log.Info("JWKs in database", slog.Int("length", len(hcjwks)))
	for _, hcjwk = range hcjwks {
		pubstate, _ := hcjwk.GetPublicKeyState()
		privstate, _ := hcjwk.GetPrivateKeyState(r.cfg.Secrets)
		r.log.Info(
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
