package model

import (
	"time"

	"gorm.io/gorm"
)

type IpaCert struct {
	gorm.Model
	IpaID          uint
	Issuer         string
	Nickname       string
	NotValidAfter  time.Time
	NotValidBefore time.Time
	Pem            string
	SerialNumber   string
	Subject        string

	Ipa Ipa
}
