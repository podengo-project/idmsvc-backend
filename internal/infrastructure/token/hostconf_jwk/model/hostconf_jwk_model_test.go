package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/secrets"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHostconfJwk(t *testing.T) {
	config := test.GetTestConfig()
	expiresAt := time.
		Now().
		Add(config.Application.HostconfJwkValidity).
		Truncate(time.Second)
	hc, err := NewHostconfJwk(config.Secrets, expiresAt)
	assert.Nil(t, err)

	assert.NotNil(t, hc.CreatedAt)
	assert.Equal(t, hc.CreatedAt, hc.UpdatedAt)
	assert.NotNil(t, hc.KeyId)
	assert.Equal(t, hc.ExpiresAt, expiresAt)
	assert.NotNil(t, hc.PublicJwk)
	assert.NotNil(t, hc.EncryptedJwk)
	assert.Equal(t, hc.EncryptionId, config.Secrets.HostconfEncryptionId)

	k, err := jwk.ParseKey([]byte(hc.PublicJwk))
	assert.Nil(t, err)
	assert.Equal(t, k.KeyID(), hc.KeyId)

	state, err := hc.GetPublicKeyState()
	assert.Nil(t, err)
	assert.Equal(t, state, hostconf_jwk.ValidKey)

	pubkey, state, err := hc.GetPublicJWK()
	assert.Nil(t, err)
	assert.Equal(t, state, hostconf_jwk.ValidKey)
	assert.Equal(t, k, pubkey)

	state, err = hc.GetPrivateKeyState(config.Secrets)
	assert.Nil(t, err)
	assert.Equal(t, state, hostconf_jwk.ValidKey)

	privkey, state, err := hc.GetPrivateJWK(config.Secrets)
	assert.Nil(t, err)
	assert.Equal(t, state, hostconf_jwk.ValidKey)

	privpubkey, err := privkey.PublicKey()
	assert.Nil(t, err)
	assert.Equal(t, pubkey, privpubkey)
}

func TestHostconfJwkMethods(t *testing.T) {
	var (
		err   error
		state hostconf_jwk.KeyState
	)
	type TestCaseGiven struct {
		HC     *HostconfJwk
		Secret *secrets.AppSecrets
	}
	type TestCaseExpected struct {
		PubState  hostconf_jwk.KeyState
		PubError  error
		PrivState hostconf_jwk.KeyState
		PrivError error
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}

	expiresFuture := time.Now().Add(time.Hour)
	expiresPast := time.Now().Add(-time.Hour)

	sec, err := secrets.NewAppSecrets("3cBBUQSnlKHQO7-5hyxJRQ")
	assert.Nil(t, err)

	sec2, err := secrets.NewAppSecrets("MLaVBnV5kadqAasiUmtEwg")
	assert.Nil(t, err)

	privkey, err := hostconf_jwk.GeneratePrivateJWK(expiresFuture)
	assert.Nil(t, err)
	encryptedJwk, err := hostconf_jwk.EncryptJWK(sec.HostConfEncryptionKey, privkey)
	assert.Nil(t, err)

	pubkey, err := hostconf_jwk.GetPublicJWK(privkey)
	assert.Nil(t, err)
	pubkeybytes, err := json.Marshal(pubkey)
	assert.Nil(t, err)
	kid := privkey.KeyID()

	testCases := []TestCase{
		{
			Name: "All valid",
			Given: TestCaseGiven{
				HC: &HostconfJwk{
					KeyId:        kid,
					ExpiresAt:    expiresFuture,
					PublicJwk:    string(pubkeybytes),
					EncryptionId: sec.HostconfEncryptionId,
					EncryptedJwk: encryptedJwk,
				},
				Secret: sec,
			},
			Expected: TestCaseExpected{
				PubState:  hostconf_jwk.ValidKey,
				PubError:  nil,
				PrivState: hostconf_jwk.ValidKey,
				PrivError: nil,
			},
		},
		{
			Name: "Invalid (missing pubkey)",
			Given: TestCaseGiven{
				HC: &HostconfJwk{
					KeyId:        kid,
					ExpiresAt:    expiresFuture,
					PublicJwk:    "",
					EncryptionId: sec.HostconfEncryptionId,
					EncryptedJwk: encryptedJwk,
				},
				Secret: sec,
			},
			Expected: TestCaseExpected{
				PubState:  hostconf_jwk.InvalidKey,
				PubError:  ErrInvalidKey,
				PrivState: hostconf_jwk.InvalidKey,
				PrivError: ErrInvalidKey,
			},
		},
		{
			Name: "Invalid (missing kid)",
			Given: TestCaseGiven{
				HC: &HostconfJwk{
					KeyId:        "",
					ExpiresAt:    expiresFuture,
					PublicJwk:    string(pubkeybytes),
					EncryptionId: sec.HostconfEncryptionId,
					EncryptedJwk: encryptedJwk,
				},
				Secret: sec,
			},
			Expected: TestCaseExpected{
				PubState:  hostconf_jwk.InvalidKey,
				PubError:  ErrInvalidKey,
				PrivState: hostconf_jwk.InvalidKey,
				PrivError: ErrInvalidKey,
			},
		},
		{
			Name: "Expired",
			Given: TestCaseGiven{
				HC: &HostconfJwk{
					KeyId:        kid,
					ExpiresAt:    expiresPast,
					PublicJwk:    string(pubkeybytes),
					EncryptionId: sec.HostconfEncryptionId,
					EncryptedJwk: encryptedJwk,
				},
				Secret: sec,
			},
			Expected: TestCaseExpected{
				PubState:  hostconf_jwk.ExpiredKey,
				PubError:  ErrExpiredKey,
				PrivState: hostconf_jwk.ExpiredKey,
				PrivError: ErrExpiredKey,
			},
		},
		{
			Name: "Revoked private key",
			Given: TestCaseGiven{
				HC: &HostconfJwk{
					KeyId:        kid,
					ExpiresAt:    expiresFuture,
					PublicJwk:    string(pubkeybytes),
					EncryptionId: sec.HostconfEncryptionId,
					EncryptedJwk: nil,
				},
				Secret: sec,
			},
			Expected: TestCaseExpected{
				PubState:  hostconf_jwk.RevokedKey,
				PubError:  ErrRevokedKey,
				PrivState: hostconf_jwk.RevokedKey,
				PrivError: ErrRevokedKey,
			},
		},
		{
			Name: "Revoked and expired key returns ExpiredKey",
			Given: TestCaseGiven{
				HC: &HostconfJwk{
					KeyId:        kid,
					ExpiresAt:    expiresPast,
					PublicJwk:    string(pubkeybytes),
					EncryptionId: sec.HostconfEncryptionId,
					EncryptedJwk: nil,
				},
				Secret: sec,
			},
			Expected: TestCaseExpected{
				PubState:  hostconf_jwk.ExpiredKey,
				PubError:  ErrExpiredKey,
				PrivState: hostconf_jwk.ExpiredKey,
				PrivError: ErrExpiredKey,
			},
		},
		{
			Name: "Private key with different encryption id",
			Given: TestCaseGiven{
				HC: &HostconfJwk{
					KeyId:        kid,
					ExpiresAt:    expiresFuture,
					PublicJwk:    string(pubkeybytes),
					EncryptionId: sec.HostconfEncryptionId,
					EncryptedJwk: encryptedJwk,
				},
				// different secret
				Secret: sec2,
			},
			Expected: TestCaseExpected{
				PubState:  hostconf_jwk.ValidKey,
				PubError:  nil,
				PrivState: hostconf_jwk.EncryptionIdMismatch,
				PrivError: ErrKeyDecryptionFailed,
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		hc := testCase.Given.HC

		state, err = hc.GetPublicKeyState()
		require.Equal(t, testCase.Expected.PubState, state)
		require.Equal(t, testCase.Expected.PubError, err)

		pubkeyout, state, err := hc.GetPublicJWK()
		require.Equal(t, testCase.Expected.PubState, state)
		require.Equal(t, testCase.Expected.PubError, err)

		if testCase.Expected.PubState == hostconf_jwk.ValidKey {
			assert.Equal(t, pubkeyout, pubkey)
		} else {
			assert.Nil(t, pubkeyout)
		}

		state, err = hc.GetPrivateKeyState(*testCase.Given.Secret)
		require.Equal(t, testCase.Expected.PrivState, state)
		require.Equal(t, testCase.Expected.PrivError, err)

		privkeyout, state, err := hc.GetPrivateJWK(*testCase.Given.Secret)
		require.Equal(t, testCase.Expected.PrivState, state)
		require.Equal(t, testCase.Expected.PrivError, err)

		if testCase.Expected.PrivState == hostconf_jwk.ValidKey {
			assert.Equal(t, privkeyout, privkey)
		} else {
			assert.Nil(t, privkeyout)
		}
	}
}
