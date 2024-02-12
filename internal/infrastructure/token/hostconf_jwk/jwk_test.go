package hostconf_jwk

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateJWK(t *testing.T) {
	expiration := time.Now().Add(time.Hour)
	key, err := GeneratePrivateJWK(expiration)
	assert.NoError(t, err)

	raw, ok := key.(jwk.ECDSAPrivateKey)
	assert.True(t, ok)
	assert.Equal(t, raw.Crv(), jwa.P256)

	assert.Equal(t, key.KeyType(), jwa.EC)
	assert.Equal(t, len(key.KeyID()), 8)

	// exp is an int64 almost equal to current Unix time
	expif, ok := key.Get("exp")
	assert.True(t, ok)
	exp, ok := expif.(int64)
	assert.True(t, ok)
	assert.Equal(t, exp, expiration.Unix())

	state, err := checkJWK(key)
	assert.Equal(t, state, ValidKey)
	assert.NoError(t, err)
}

func TestGetPublicK(t *testing.T) {
	expiration := time.Now().Add(time.Hour)
	key, err := GeneratePrivateJWK(expiration)
	assert.NoError(t, err)

	pub, err := GetPublicJWK(key)
	_, ok := pub.(jwk.ECDSAPublicKey)
	assert.True(t, ok)
	_, ok = pub.(jwk.ECDSAPrivateKey)
	assert.False(t, ok)

	pub2, err := GetPublicJWK(pub)
	assert.NoError(t, err)
	assert.Equal(t, pub, pub2)
}

func TestParseJWK(t *testing.T) {
	expiration := time.Now().Add(time.Hour)
	key, err := GeneratePrivateJWK(expiration)
	assert.NoError(t, err)

	s, err := json.Marshal(key)
	assert.NoError(t, err)
	parsed, state, err := ParseJWK(s)
	assert.NoError(t, err)
	assert.Equal(t, state, ValidKey)
	assert.Equal(t, parsed, key)

	pub, err := GetPublicJWK(key)
	assert.NoError(t, err)

	s, err = json.Marshal(pub)
	assert.NoError(t, err)
	parsed, state, err = ParseJWK(s)
	assert.Equal(t, state, ValidKey)
	assert.Equal(t, parsed, pub)

	// TODO add tests for invalid and expired keys
}
