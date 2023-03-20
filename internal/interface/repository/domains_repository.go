package repository

import (
	"github.com/hmsidm/internal/domain/model"
	"gorm.io/gorm"
)

type DomainRepository interface {
	FindAll(db *gorm.DB, orgId string, offset int64, count int32) (output []model.Domain, err error)
	Create(db *gorm.DB, orgId string, data *model.Domain) (err error)
	// PartialUpdate(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error)
	// Update(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error)
	FindById(db *gorm.DB, orgId string, uuid string) (output model.Domain, err error)
	DeleteById(db *gorm.DB, orgId string, uuid string) (err error)
	Update(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error)
}
