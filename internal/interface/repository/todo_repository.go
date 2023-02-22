package repository

import (
	"github.com/hmsidm/internal/domain/model"
	"gorm.io/gorm"
)

type TodoRepository interface {
	FindAll(db *gorm.DB, offset int64, count int32) (output []model.Todo, err error)
	Create(db *gorm.DB, data *model.Todo) (err error)
	PartialUpdate(db *gorm.DB, data *model.Todo) (output model.Todo, err error)
	Update(db *gorm.DB, data *model.Todo) (output model.Todo, err error)
	FindById(db *gorm.DB, id uint) (output model.Todo, err error)
	DeleteById(db *gorm.DB, id uint) (err error)
}
