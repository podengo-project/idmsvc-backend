package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHkdfExpand(t *testing.T) {
	mainSecret := []byte("mainSecret")
	prk := HkdfExtract(mainSecret)
	assert.Equal(t, prk, PRK{0xfc, 0x66, 0x8, 0xeb, 0x2f, 0x83, 0x77, 0x93, 0x2, 0x1a, 0x57, 0xf2, 0x1d, 0x39, 0x8a, 0x5, 0xa, 0x48, 0x9a, 0x63, 0x8e, 0xf9, 0x57, 0xfb, 0xac, 0x7a, 0x21, 0x63, 0x50, 0x3f, 0xac, 0x69})
}

// from cryptography.hazmat.primitives import hashes
// from cryptography.hazmat.primitives.kdf.hkdf import HKDF
// HKDF(hashes.SHA256(), 8, b"idmsvc-backend", b"test").derive(b"mainSecret")

func TestHkdfExtract(t *testing.T) {
	mainSecret := []byte("mainSecret")
	prk := HkdfExtract(mainSecret)
	info := HkdfInfo{[]byte("test"), 8}
	secret, err := HkdfExpand(prk, info)
	assert.NoError(t, err)
	assert.Equal(t, secret, []byte{0xe7, 0x84, 0x18, 0xae, 0xc6, 0x4d, 0xe5, 0x42})
}
