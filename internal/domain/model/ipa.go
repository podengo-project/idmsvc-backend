package model

import "gorm.io/gorm"

// See: https://gorm.io/docs/models.html

type Ipa struct {
	gorm.Model
	DomainID  uint
	RealmName *string
	CaList    *string
	// TODO Do we want to create a ipa_server_list table
	//      related with this Ipa entry?
	// NOTE Thinking about this as a comma separated list
	//      of servers
	ServerList *string
	Domain     Domain `gorm:"foreignKey:id"`
}
