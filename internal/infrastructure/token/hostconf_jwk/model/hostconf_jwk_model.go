package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/secrets"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk"
	"gorm.io/gorm"
)

// HostconfJwks hold public and private JWKs
type HostconfJwk struct {
	gorm.Model
	KeyId        string
	ExpiresAt    time.Time
	PublicJwk    string
	EncryptionId string
	EncryptedJwk []byte
}

var (
	ErrExpiredKey          = errors.New("expired key")
	ErrInvalidKey          = errors.New("invalid key")
	ErrRevokedKey          = errors.New("revoked key")
	ErrKeyDecryptionFailed = errors.New("decryption failed")
)

// Create a new Hostconf JWK entry with public and encrypted private JWK
func NewHostconfJwk(secrets secrets.AppSecrets, expiresAt time.Time) (hc *HostconfJwk, err error) {
	var (
		encryptedJwk []byte
		pubkey       jwk.Key
		pubkeybytes  []byte
		privkey      jwk.Key
	)

	// create an encrypt private key
	if privkey, err = hostconf_jwk.GeneratePrivateJWK(expiresAt); err != nil {
		return nil, err
	}
	if encryptedJwk, err = hostconf_jwk.EncryptJWK(secrets.HostConfEncryptionKey, privkey); err != nil {
		return nil, err
	}

	// serialize public key
	if pubkey, err = hostconf_jwk.GetPublicJWK(privkey); err != nil {
		return nil, err
	}
	if pubkeybytes, err = json.Marshal(pubkey); err != nil {
		return nil, err
	}

	hc = &HostconfJwk{
		KeyId:        pubkey.KeyID(),
		ExpiresAt:    expiresAt,
		PublicJwk:    string(pubkeybytes),
		EncryptedJwk: encryptedJwk,
		EncryptionId: secrets.HostconfEncryptionId,
	}

	return hc, nil
}

// Get public key state (invalid, expired, revoked, valid)
// A public key can be valid although its private key cannot be decrypted
// by current secret.
func (hc *HostconfJwk) GetPublicKeyState() (hostconf_jwk.KeyState, error) {
	if hc.PublicJwk == "" || hc.KeyId == "" {
		return hostconf_jwk.InvalidKey, ErrInvalidKey
	}
	if hc.ExpiresAt.Unix() <= time.Now().Unix() {
		return hostconf_jwk.ExpiredKey, ErrExpiredKey
	}
	if hc.EncryptedJwk == nil {
		return hostconf_jwk.RevokedKey, ErrRevokedKey
	}
	return hostconf_jwk.ValidKey, nil
}

// Get public jwk.Key from entry
// Fails if key is invalid, expired, or revoked.
func (hc *HostconfJwk) GetPublicJWK() (pubkey jwk.Key, state hostconf_jwk.KeyState, err error) {
	if state, err = hc.GetPublicKeyState(); err != nil {
		return nil, state, err
	}
	return hostconf_jwk.ParseJWK([]byte(hc.PublicJwk))
}

// Get private key state (invalid, expired, revoked, mismatch, valid)
func (hc *HostconfJwk) GetPrivateKeyState(secrets secrets.AppSecrets) (state hostconf_jwk.KeyState, err error) {
	state, err = hc.GetPublicKeyState()
	if state == hostconf_jwk.ValidKey {
		if hc.EncryptionId != secrets.HostconfEncryptionId {
			return hostconf_jwk.EncryptionIdMismatch, ErrKeyDecryptionFailed
		}
	}
	return state, err
}

// Decrypt and return private jwk.Key from entry
// Fails if key is invalid, expired, revoked, or not encrypted with secret.
func (hc *HostconfJwk) GetPrivateJWK(secrets secrets.AppSecrets) (privkey jwk.Key, state hostconf_jwk.KeyState, err error) {
	if state, err = hc.GetPrivateKeyState(secrets); err != nil {
		return nil, state, err
	}
	if privkey, err = hostconf_jwk.DecryptJWK(secrets.HostConfEncryptionKey, hc.EncryptedJwk); err != nil {
		return nil, hostconf_jwk.KeyDecryptionFailed, err
	}
	return privkey, state, err
}
