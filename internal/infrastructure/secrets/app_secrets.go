package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type AppSecrets struct {
	DomainRegKey          []byte
	HostconfEncryptionId  string
	HostConfEncryptionKey []byte
}

const (
	MainSecretMinLength = 16
)

// Generate random value for main secret, used in tests
func GenerateRandomMainSecret() string {
	random := make([]byte, MainSecretMinLength)
	if _, err := rand.Read(random); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(random)
}

// Parse main secret and get sub secrets
func NewAppSecrets(mainSecret string) (sec *AppSecrets, err error) {
	var secret []byte

	if mainSecret == "random" {
		return nil, fmt.Errorf("random main secret is no longer supported")
	}
	if secret, err = base64.RawURLEncoding.DecodeString(mainSecret); err != nil {
		return nil, fmt.Errorf("Failed to decode main secret: %v", err)
	}
	if len(secret) < MainSecretMinLength {
		return nil, fmt.Errorf("Main secret is too short, expected at least %d bytes.", MainSecretMinLength)
	}

	// extract PRK from main secret
	prk := HkdfExtract(secret)

	sec = &AppSecrets{}
	sec.DomainRegKey, err = HkdfExpand(prk, DomainRegKeyInfo)
	if err != nil {
		return nil, err
	}
	sec.HostConfEncryptionKey, err = HkdfExpand(prk, HostconfEncryptionKeyInfo)
	if err != nil {
		return nil, err
	}
	encid, err := HkdfExpand(prk, HostconfEncryptionIdInfo)
	if err != nil {
		return nil, err
	}
	sec.HostconfEncryptionId = hex.EncodeToString(encid)

	return sec, nil
}
