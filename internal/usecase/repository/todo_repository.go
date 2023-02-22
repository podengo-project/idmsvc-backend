package repository

import (
	"fmt"

	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/repository"
	"gorm.io/gorm"
)

type todoRepository struct{}

func NewTodoRepository() repository.TodoRepository {
	return &todoRepository{}
}

func (r *todoRepository) FindAll(db *gorm.DB, offset int64, limit int32) (data []model.Todo, err error) {
	if db == nil {
		return []model.Todo{}, fmt.Errorf("db is nil")
	}
	if offset < 0 {
		return []model.Todo{}, fmt.Errorf("offset lower than 0")
	}
	if limit < 0 {
		return []model.Todo{}, fmt.Errorf("limit lower than 0")
	}
	if err = db.Offset(int(offset)).Limit(int(limit)).Find(&data).Error; err != nil {
		return []model.Todo{}, err
	}
	return
}

func (r *todoRepository) Create(db *gorm.DB, data *model.Todo) (err error) {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if data == nil {
		return fmt.Errorf("data is nil")
	}
	if err = db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *todoRepository) getColumnsToUpdate(data *model.Todo) []string {
	cols := []string{"id"}
	if data.Title != nil {
		cols = append(cols, "title")
	}
	if data.Description != nil {
		cols = append(cols, "description")
	}
	return cols
}

func (r *todoRepository) PartialUpdate(db *gorm.DB, data *model.Todo) (output model.Todo, err error) {
	if db == nil {
		return model.Todo{}, fmt.Errorf("db is nil")
	}
	if data == nil {
		return model.Todo{}, fmt.Errorf("data is nil")
	}
	cols := r.getColumnsToUpdate(data)
	if err = db.Model(data).Select(cols).Updates(*data).Error; err != nil {
		return model.Todo{}, err
	}
	return *data, nil
}

func (r *todoRepository) Update(db *gorm.DB, data *model.Todo) (output model.Todo, err error) {
	if db == nil {
		return model.Todo{}, fmt.Errorf("db is nil")
	}
	if data == nil {
		return model.Todo{}, fmt.Errorf("data is nil")
	}
	cols := []string{"id"}
	if err = db.Model(data).Select(cols).Updates(*data).Error; err != nil {
		return model.Todo{}, err
	}
	return *data, nil
}

func (r *todoRepository) FindById(db *gorm.DB, id uint) (output model.Todo, err error) {
	if db == nil {
		return model.Todo{}, fmt.Errorf("db is nil")
	}
	if err = db.First(&output, int(id)).Error; err != nil {
		return model.Todo{}, err
	}
	return output, nil
}

func (r *todoRepository) DeleteById(db *gorm.DB, id uint) (err error) {
	var data model.Todo
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if err = db.First(&data, id).Error; err != nil {
		return err
	}
	if err = db.Delete(&data, id).Error; err != nil {
		return err
	}
	return nil
}
