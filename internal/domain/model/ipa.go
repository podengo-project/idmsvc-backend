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
// See: https://gorm.io/docs/conventions.html

// Ipa represent the specific rhel-idm domain information
type Ipa struct {
	gorm.Model
	CaCerts         []IpaCert
	Servers         []IpaServer
	RealmName       *string
	RealmDomains    pq.StringArray `gorm:"type:text[]"`
	Token           *string
	TokenExpiration *time.Time `gorm:"column:token_expiration_ts"`

	Domain Domain `gorm:"foreignKey:ID;references:ID"`
}

// tokenExpirationDuration is the duration to be used
// when a domain is created into the database.
var tokenExpirationDuration time.Duration = time.Hour * 2

// SetDefaultTokenExpiration update the default expiration
// period. This value is global value and is intendeed to be
// set on initialization time.
// d is the period of time during the token is valid.
func SetDefaultTokenExpiration(d time.Duration) {
	tokenExpirationDuration = d
}

// DefaultTokenExpiration get the value used to set
// the expiration period for a new domain created.
// Return the duration to be used for new creted domains
// into the database.
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
	*tokenExpiration = time.Now().
		Add(DefaultTokenExpiration())
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
