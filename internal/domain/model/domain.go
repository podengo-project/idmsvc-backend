package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// See: https://gorm.io/docs/models.html

const (
	DomainTypeUndefined uint = iota
	DomainTypeIpa
	// DomainTypeAzure
	// DomainTypeActiveDirector
)

// NOTE https://samu.space/uuids-with-postgres-and-gorm/
//      thanks @anschnei
// NOTE hmscontent can be an example of this; they redefine
//      the base model of the gorm models to use uuid as
//      the primary key

type Domain struct {
	gorm.Model
	OrgId                 string
	DomainUuid            uuid.UUID `gorm:"unique"`
	DomainName            *string
	RealmName             *string
	DomainType            *uint
	AutoEnrollmentEnabled *bool
	IpaDomain             *Ipa
}

func DomainTypeString(data uint) string {
	switch data {
	case DomainTypeIpa:
		return "ipa"
	default:
		return ""
	}
}

func DomainTypeUint(data string) uint {
	switch data {
	case "ipa":
		return DomainTypeIpa
	default:
		return DomainTypeUndefined
	}
}

// See: https://gorm.io/docs/hooks.html

func (d *Domain) BeforeCreate(tx *gorm.DB) (err error) {
	d.DomainUuid = uuid.New()
	var currentTime = time.Now()
	d.CreatedAt = currentTime
	d.UpdatedAt = currentTime

	return nil
}
