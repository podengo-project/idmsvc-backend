package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type AppSecrets struct {
	DomainRegKey []byte
}

const (
	MainSecretMinLength = 16
)

// Parse main secret and get sub secrets
func NewAppSecrets(mainSecret string) (sec *AppSecrets, err error) {
	var secret []byte
	if mainSecret == "random" {
		secret = make([]byte, MainSecretMinLength)
		if _, err = rand.Read(secret); err != nil {
			return nil, err
		}
	} else {
		if secret, err = base64.RawURLEncoding.DecodeString(mainSecret); err != nil {
			return nil, fmt.Errorf("Failed to main secret: %v", err)
		}
		if len(secret) < MainSecretMinLength {
			return nil, fmt.Errorf("Main secret is too short, expected at least %d bytes.", MainSecretMinLength)
		}
	}

	// extract PRK from main secret
	prk := HkdfExtract(secret)

	sec = &AppSecrets{}
	sec.DomainRegKey, err = HkdfExpand(prk, DomainRegKeyInfo)
	if err != nil {
		return nil, err
	}
	return sec, nil
}
