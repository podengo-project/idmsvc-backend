package model

import (
	"fmt"

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
	if i.Locations != nil {
		for idx := range i.Locations {
			i.Locations[idx].IpaID = i.ID
		}
	}
	return nil
}
