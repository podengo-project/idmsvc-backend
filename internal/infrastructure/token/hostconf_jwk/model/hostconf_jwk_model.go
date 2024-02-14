package model

import (
	"time"

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

// Before create hook
func (hc *HostconfJwk) BeforeCreate(tx *gorm.DB) (err error) {
	var currentTime = time.Now()
	hc.CreatedAt = currentTime
	hc.UpdatedAt = currentTime

	return nil
}
