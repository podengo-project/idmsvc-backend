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

func (r *domainRepository) FindAll(db *gorm.DB, orgId string, offset int64, count int32) (output []model.Domain, err error) {
	if db == nil {
		return []model.Domain{}, fmt.Errorf("db is nil")
	}
	if offset < 0 {
		return []model.Domain{}, fmt.Errorf("offset lower than 0")
	}
	if count < 0 {
		return []model.Domain{}, fmt.Errorf("limit lower than 0")
	}
	err = db.Offset(int(offset)).Limit(int(count)).Find(&output).Error
	if err != nil {
		return []model.Domain{}, err
	}
	return
}

func (r *domainRepository) Create(db *gorm.DB, orgId string, data *model.Domain) (err error) {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if data == nil {
		return fmt.Errorf("data is nil")
	}
	if data.IpaDomain == nil {
		return fmt.Errorf("data.IpaDomain is nil")
	}
	data.OrgId = orgId
	err = db.Create(data).Error
	if err != nil {
		return err
	}
	data.IpaDomain.DomainID = data.ID
	if err = db.Create(data.IpaDomain).Error; err != nil {
		return err
	}
	for _, cert := range data.IpaDomain.CaCerts {
		cert.IpaID = data.IpaDomain.ID
		if err = db.Create(&cert).Error; err != nil {
			return err
		}
	}
	for _, server := range data.IpaDomain.Servers {
		server.IpaID = data.IpaDomain.ID
		if err = db.Create(&server).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *domainRepository) getColumnsToUpdate(data *model.Domain) []string {
	cols := []string{"id"}
	if data.AutoEnrollmentEnabled != nil {
		cols = append(cols, "auto_enrollment_enabled")
	}
	if data.DomainName != nil {
		cols = append(cols, "domain_name")
	}
	if data.DomainType != nil {
		cols = append(cols, "domain_type")
	}
	return cols
}

// func (r *domainRepository) PartialUpdate(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error) {
// 	if db == nil {
// 		return model.Domain{}, fmt.Errorf("db is nil")
// 	}
// 	if data == nil {
// 		return model.Domain{}, fmt.Errorf("data is nil")
// 	}
// 	data.OrgId = orgId
// 	cols := r.getColumnsToUpdate(data)
// 	if err = db.Model(data).Select(cols).Updates(*data).Error; err != nil {
// 		return model.Domain{}, err
// 	}
// 	return *data, nil
// }

// func (r *domainRepository) Update(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error) {
// 	if db == nil {
// 		return model.Domain{}, fmt.Errorf("db is nil")
// 	}
// 	if data == nil {
// 		return model.Domain{}, fmt.Errorf("data is nil")
// 	}
// 	data.OrgId = orgId
// 	cols := []string{"id"}
// 	if err = db.Model(data).Select(cols).Updates(*data).Error; err != nil {
// 		return model.Domain{}, err
// 	}
// 	return *data, nil
// }

// See: https://gorm.io/docs/query.html
// TODO Document the method
func (r *domainRepository) FindById(db *gorm.DB, orgId string, uuid string) (output model.Domain, err error) {
	var count int64
	if db == nil {
		return model.Domain{}, fmt.Errorf("db is nil")
	}
	if err = db.First(&output, "org_id = ? AND domain_uuid = ?", orgId, uuid).Count(&count).Error; err != nil {
		return model.Domain{}, err
	}
	if count == 0 {
		return model.Domain{}, fmt.Errorf("Not found")
	}
	if output.DomainType != nil && *output.DomainType == model.DomainTypeIpa {
		if err = db.First(&output.IpaDomain, "domain_id = ?", output.ID).Count(&count).Error; err != nil {
			return model.Domain{}, err
		}
		if count == 0 {
			return model.Domain{}, fmt.Errorf("Not found")
		}
	}
	return output, nil
}

// See: https://gorm.io/docs/delete.html
func (r *domainRepository) DeleteById(db *gorm.DB, orgId string, uuid string) (err error) {
	var (
		data  model.Domain
		count int64
	)
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if orgId == "" {
		return fmt.Errorf("orgId cannot be an empty string")
	}
	if uuid == "" {
		return fmt.Errorf("uuid cannot be an empty string")
	}
	if err = db.First(&data, "org_id = ? AND domain_uuid = ?", orgId, uuid).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Register not found")
	}
	if err = db.Unscoped().Delete(&data).Error; err != nil {
		return err
	}
	return nil
}
