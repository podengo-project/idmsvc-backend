package repository

import (
	"fmt"

	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/repository"
	"gorm.io/gorm"
)

type domainRepository struct{}

func NewDomainRepository() repository.DomainRepository {
	return &domainRepository{}
}

func (r *domainRepository) FindAll(db *gorm.DB, offset int64, limit int32) (data []model.Domain, err error) {
	if db == nil {
		return []model.Domain{}, fmt.Errorf("db is nil")
	}
	if offset < 0 {
		return []model.Domain{}, fmt.Errorf("offset lower than 0")
	}
	if limit < 0 {
		return []model.Domain{}, fmt.Errorf("limit lower than 0")
	}
	if err = db.Offset(int(offset)).Limit(int(limit)).Find(&data).Error; err != nil {
		return []model.Domain{}, err
	}
	return
}

func (r *domainRepository) Create(db *gorm.DB, data *model.Domain) (err error) {
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

func (r *domainRepository) getColumnsToUpdate(data *model.Domain) []string {
	cols := []string{"id"}
	if data.Title != nil {
		cols = append(cols, "title")
	}
	if data.Description != nil {
		cols = append(cols, "description")
	}
	return cols
}

func (r *domainRepository) PartialUpdate(db *gorm.DB, data *model.Domain) (output model.Domain, err error) {
	if db == nil {
		return model.Domain{}, fmt.Errorf("db is nil")
	}
	if data == nil {
		return model.Domain{}, fmt.Errorf("data is nil")
	}
	cols := r.getColumnsToUpdate(data)
	if err = db.Model(data).Select(cols).Updates(*data).Error; err != nil {
		return model.Domain{}, err
	}
	return *data, nil
}

func (r *domainRepository) Update(db *gorm.DB, data *model.Domain) (output model.Domain, err error) {
	if db == nil {
		return model.Domain{}, fmt.Errorf("db is nil")
	}
	if data == nil {
		return model.Domain{}, fmt.Errorf("data is nil")
	}
	cols := []string{"id"}
	if err = db.Model(data).Select(cols).Updates(*data).Error; err != nil {
		return model.Domain{}, err
	}
	return *data, nil
}

func (r *domainRepository) FindById(db *gorm.DB, id uint) (output model.Domain, err error) {
	if db == nil {
		return model.Domain{}, fmt.Errorf("db is nil")
	}
	if err = db.First(&output, int(id)).Error; err != nil {
		return model.Domain{}, err
	}
	return output, nil
}

func (r *domainRepository) DeleteById(db *gorm.DB, id uint) (err error) {
	var data model.Domain
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
