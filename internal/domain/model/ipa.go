package model

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// See: https://gorm.io/docs/models.html
// See: https://gorm.io/docs/conventions.html

// Ipa represent the specific rhel-idm domain information
type Ipa struct {
	gorm.Model
	CaCerts      []IpaCert
	Servers      []IpaServer
	Locations    []IpaLocation
	RealmName    *string
	RealmDomains pq.StringArray `gorm:"type:text[]"`

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
