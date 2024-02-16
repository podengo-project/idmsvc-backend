package secrets

import (
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/hkdf"
)

type PRK []byte

type HkdfInfo struct {
	Info   []byte
	Length int
}

const (
	Salt = "idmsvc-backend"
)

var (
	// MAC key for domain registration token
	DomainRegKeyInfo = HkdfInfo{[]byte("domain registration key"), 32}
	// hex string to identify AES encryption keys for encrypted private JWKs
	HostconfEncryptionIdInfo = HkdfInfo{[]byte("hostconf JWK encryption id"), 8}
	// AES-GCM encryption keys for private JWKs
	HostconfEncryptionKeyInfo = HkdfInfo{[]byte("hostconf JWK encryption key"), 16}
)

// Extract pseudo random key from a secret
func HkdfExtract(mainSecret []byte) PRK {
	var hash = sha256.New
	return PRK(hkdf.Extract(hash, mainSecret, []byte(Salt)))
}

// Expand pseudo random key into a secret
func HkdfExpand(prk PRK, hi HkdfInfo) (secret []byte, err error) {
	var hash = sha256.New
	reader := hkdf.Expand(hash, prk, hi.Info)
	secret = make([]byte, hi.Length)
	if _, err := io.ReadFull(reader, secret); err != nil {
		return nil, err
	}
	return secret, err
}
