package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// See: https://gorm.io/docs/models.html

/*
{
  "domain_uuid": "1aa15eae-a88b-11ed-a2cb-482ae3863d30",
  "domain_name": "mydomain.example",
  "domain_type": "ipa",
  "auto_enrollment_enabled": true,
  "ipa": {
    "realm_name": "IPA.EXAMPLE",
    "ca_list": "base64",
    "server_list": [
      "server1.mydomain.example",
      "Server2.mydomain.example"
    ],
    "client_options": {}
  }
}
*/

const (
	DomainTypeIpa uint = iota + 1
	DomainTypeAzure
	DomainTypeActiveDirector
)

// NOTE https://samu.space/uuids-with-postgres-and-gorm/
//      thanks @anschnei
// NOTE hmscontent can be an example of this; they redefine
//      the base model of the gorm models to use uuid as
//      the primary key

type Domain struct {
	gorm.Model
	OrgId                 string
	DomainUuid            *uuid.UUID `gorm:"unique"`
	DomainName            *string
	DomainType            uint
	AutoEnrollmentEnabled bool
	Title                 *string
	Description           *string
}
