package repository

import (
	"github.com/hmsidm/internal/domain/model"
	"gorm.io/gorm"
)

type DomainRepository interface {
	FindAll(db *gorm.DB, offset int64, count int32) (output []model.Domain, err error)
	Create(db *gorm.DB, data *model.Domain) (err error)
	// PartialUpdate(db *gorm.DB, data *model.Domain) (output model.Domain, err error)
	// Update(db *gorm.DB, data *model.Domain) (output model.Domain, err error)
	FindById(db *gorm.DB, id uint) (output model.Domain, err error)
	DeleteById(db *gorm.DB, id uint) (err error)
}
