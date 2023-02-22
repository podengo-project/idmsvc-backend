package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title       *string
	Description *string
}
