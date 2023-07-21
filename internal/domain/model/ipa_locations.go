package model

import "gorm.io/gorm"

// IpaLocation represent the possible locations for
// a rhel-idm instance.
type IpaLocation struct {
	gorm.Model
	IpaID       uint
	Name        string
	Description *string
	Ipa         Ipa `gorm:"foreignKey:ID;references:IpaID"`
}
