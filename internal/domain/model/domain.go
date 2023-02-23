package model

import (
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
	DomainType            *uint
	AutoEnrollmentEnabled *bool

	IpaDomain *Ipa `gorm:"foreignKey:DomainID"`
}
