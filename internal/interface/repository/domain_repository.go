package repository

import (
	"github.com/hmsidm/internal/domain/model"
	"gorm.io/gorm"
)

// DomainRepository interface
type DomainRepository interface {
	List(db *gorm.DB, orgID string, offset int, limit int) (output []model.Domain, count int64, err error)
	Create(db *gorm.DB, orgID string, data *model.Domain) (err error)
	// PartialUpdate(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error)
	// Update(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error)
	FindById(db *gorm.DB, orgID string, uuid string) (output *model.Domain, err error)
	DeleteById(db *gorm.DB, orgID string, uuid string) (err error)
	Update(db *gorm.DB, orgID string, data *model.Domain) (err error)
	RhelIdmClearToken(db *gorm.DB, orgID string, uuid string) (err error)
}
