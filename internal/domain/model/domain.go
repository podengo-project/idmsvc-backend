package model

import (
	"fmt"
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

	DomainTypeIpaString       = "rhel-idm"
	DomainTypeUndefinedString = ""
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
	Title                 *string
	Description           *string
	Type                  *uint
	AutoEnrollmentEnabled *bool
	IpaDomain             *Ipa `gorm:"foreignKey:id"`
}

func DomainTypeString(data uint) string {
	switch data {
	case DomainTypeIpa:
		return DomainTypeIpaString
	default:
		return DomainTypeUndefinedString
	}
}

func DomainTypeUint(data string) uint {
	switch data {
	case DomainTypeIpaString:
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

func (d *Domain) AfterCreate(tx *gorm.DB) (err error) {
	if d.Type == nil {
		return fmt.Errorf("'DomainType' cannot be nil")
	}
	switch *d.Type {
	case DomainTypeIpa:
		{
			if d.IpaDomain != nil {
				d.IpaDomain.ID = d.ID
			}
		}
	}
	return nil
}
