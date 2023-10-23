package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/stretchr/testify/assert"
)

func TestJWKEncryption(t *testing.T) {
	raw, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	jwkey, err := jwk.FromRaw(raw)
	assert.NoError(t, err)

	aeskey1 := make([]byte, 16)
	_, err = rand.Reader.Read(aeskey1)
	assert.NoError(t, err)

	encrypted, err := EncryptJWK(aeskey1, jwkey)
	assert.NoError(t, err)

	decrypted, err := DecryptJWK(aeskey1, encrypted)
	assert.NoError(t, err)
	assert.Equal(t, decrypted, jwkey)

	// same key and value, different nonce
	encrypted2, err := EncryptJWK(aeskey1, jwkey)
	assert.NoError(t, err)
	assert.NotEqual(t, encrypted, encrypted2)

	decrypted, err = DecryptJWK(aeskey1, encrypted2)
	assert.NoError(t, err)
	assert.Equal(t, decrypted, jwkey)

	// Invalid key
	aeskey2 := make([]byte, 16)
	decrypted, err = DecryptJWK(aeskey2, encrypted)
	assert.Nil(t, decrypted)
	assert.EqualError(t, err, "cipher: message authentication failed")
}
