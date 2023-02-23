package model

import "gorm.io/gorm"

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

type Domain struct {
	gorm.Model
	DomainUuid            *string `gorm:"unique"`
	DomainName            *string
	DomainType            *string
	AutoEnrollmentEnabled bool
	Title                 *string
	Description           *string
}
