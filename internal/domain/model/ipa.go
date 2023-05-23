package model

import (
	"crypto/rand"
	"fmt"
	"time"

	b64 "encoding/base64"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"gorm.io/gorm"
)

// See: https://gorm.io/docs/models.html

type Ipa struct {
	gorm.Model
	CaCerts         []IpaCert
	Servers         []IpaServer
	RealmName       *string
	RealmDomains    pq.StringArray `gorm:"type:text[]"`
	Token           *string
	TokenExpiration *time.Time

	Domain Domain `gorm:"foreignKey:ID;references:ID"`
}

// Set by the default tokenExpiration to 24hours once it is created
var tokenExpirationDuration time.Duration = time.Hour * 24

func SetDefaultTokenExpiration(d time.Duration) {
	tokenExpirationDuration = d
}

func DefaultTokenExpiration() time.Duration {
	return tokenExpirationDuration
}

func GenerateToken(length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	sEnc := b64.StdEncoding.EncodeToString([]byte(b))
	return sEnc
}

func (i *Ipa) BeforeCreate(tx *gorm.DB) (err error) {
	if i == nil {
		return fmt.Errorf("'BeforeCreate' cannot be invoked on nil")
	}
	tokenExpiration := &time.Time{}
	*tokenExpiration = time.Now().Add(tokenExpirationDuration)
	i.Token = pointy.String(uuid.NewString())
	i.TokenExpiration = tokenExpiration

	return nil
}

func (i *Ipa) AfterCreate(tx *gorm.DB) (err error) {
	if i == nil {
		return fmt.Errorf("'AfterCreate' cannot be invoked on nil")
	}
	if i.CaCerts != nil {
		for idx := range i.CaCerts {
			i.CaCerts[idx].IpaID = i.ID
		}
	}
	if i.Servers != nil {
		for idx := range i.Servers {
			i.Servers[idx].IpaID = i.ID
		}
	}
	return nil
}
