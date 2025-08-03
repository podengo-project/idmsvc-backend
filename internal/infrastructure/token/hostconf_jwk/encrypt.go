package hostconf_jwk

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"

	"github.com/lestrrat-go/jwx/v3/jwk"
)

// Encrypt a JWK Key with encryption key
// Uses AES-GCM with a random nonce. The nonce is pre-pended to ciphertext.
func EncryptJWK(aeskey []byte, jwkey jwk.Key) (enc []byte, err error) {
	// create AES-GCM AEAD and nonce
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// serialize JWK to JSON bytes
	jwkbuf, err := json.Marshal(jwkey)
	if err != nil {
		return nil, err
	}

	// encrypt and seal, nonce is pre-pended
	return aesgcm.Seal(nonce, nonce, jwkbuf, nil), nil
}

// Decrypt an AES-GCM encrypted JWK
func DecryptJWK(aeskey []byte, encrypted []byte) (jwkey jwk.Key, err error) {
	// create AES-GCM AEAD and get nonce
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(encrypted) < (aesgcm.NonceSize() + aesgcm.Overhead()) {
		return nil, errors.New("encrypted blob is too short")
	}
	nonce := encrypted[:aesgcm.NonceSize()]
	ciphertext := encrypted[aesgcm.NonceSize():]

	// unseal and decrypt
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return jwk.ParseKey(plaintext)
}
