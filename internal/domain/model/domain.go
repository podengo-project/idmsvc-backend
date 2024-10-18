package model

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
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

var (
	NilUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
)

// NOTE https://samu.space/uuids-with-postgres-and-gorm/
//      thanks @anschnei
// NOTE hmscontent can be an example of this; they redefine
//      the base model of the gorm models to use uuid as
//      the primary key

type Domain struct {
	gorm.Model
	OrgId                 string    `gorm:"index:idx_domains_org_id"`
	DomainUuid            uuid.UUID `gorm:"unique"`
	DomainName            *string
	Title                 *string
	Description           *string
	Type                  *uint
	AutoEnrollmentEnabled *bool
	IpaDomain             *Ipa `gorm:"foreignKey:ID"`
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
func (d *Domain) AfterCreate(tx *gorm.DB) (err error) {
	if d.Type == nil {
		return internal_errors.NilArgError("Type")
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

// FillAndPreload is a helper method to fill and preload a domain object based
// TODO: use "AfterFind" hook instead?
func (d *Domain) FillAndPreload(db *gorm.DB) (err error) {
	if d.Type == nil {
		return internal_errors.NilArgError("Type")
	}
	switch *d.Type {
	case DomainTypeIpa:
		d.IpaDomain = &Ipa{}
		if err := db.
			Model(&Ipa{}).
			Preload("CaCerts").
			Preload("Servers").
			Preload("Locations").
			First(d.IpaDomain, "id = ?", d.ID).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// FIXME Something different to do here?
				return err
			} else {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("'Type' is invalid")
	}
}
