package model

import "gorm.io/gorm"

// See: https://gorm.io/docs/models.html

type IpaServer struct {
	gorm.Model
	IpaID               uint
	FQDN                string
	RHSMId              string `gorm:"unique;column:rhsm_id"`
	CaServer            bool
	HCCEnrollmentServer bool
	PKInitServer        bool

	Ipa Ipa
}
