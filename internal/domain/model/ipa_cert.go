package model

import (
	"time"

	"gorm.io/gorm"
)

type IpaCert struct {
	gorm.Model
	IpaID        uint
	Issuer       string
	Nickname     string
	NotAfter     time.Time
	NotBefore    time.Time
	Pem          string
	SerialNumber string
	Subject      string

	Ipa Ipa `gorm:"foreignKey:ID;references:IpaID"`
}
