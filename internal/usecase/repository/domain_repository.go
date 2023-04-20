package repository

import (
	"fmt"
	"time"

	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *domainRepository) createIpaDomain(db *gorm.DB, domainID uint, data *model.Ipa) (err error) {
	if data == nil {
		return fmt.Errorf("'data' of type '*model.Ipa' is nil")
	}
	data.Model.ID = domainID
	tokenExpiration := &time.Time{}
	*tokenExpiration = time.Now().Add(model.DefaultTokenExpiration()).UTC()
	data.TokenExpiration = tokenExpiration
	if err = db.Omit(clause.Associations).Create(data).Error; err != nil {
		return err
	}
	for idx := range data.CaCerts {
		data.CaCerts[idx].IpaID = data.ID
		if err = db.Create(&data.CaCerts[idx]).Error; err != nil {
			return err
		}
	}
	for idx := range data.Servers {
		data.Servers[idx].IpaID = data.ID
		if err = db.Create(&data.Servers[idx]).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *domainRepository) Create(db *gorm.DB, orgId string, data *model.Domain) (err error) {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if data == nil {
		return fmt.Errorf("data is nil")
	}
	if data.Type == nil {
		return fmt.Errorf("'Type' is nil")
	}
	data.OrgId = orgId
	err = db.Omit(clause.Associations).Create(data).Error
	if err != nil {
		return err
	}
	switch *data.Type {
	case model.DomainTypeIpa:
		err = r.createIpaDomain(db, data.ID, data.IpaDomain)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("'Type' is invalid")
	}
	return nil
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

// Update save the Domain record into the database. It only update
// data for the current organization.
func (r *domainRepository) Update(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error) {
	if db == nil {
		return model.Domain{}, fmt.Errorf("'db' cannot be nil")
	}
	if orgId == "" {
		return model.Domain{}, fmt.Errorf("'orgId' cannot be an empty string")
	}
	if data == nil {
		return model.Domain{}, fmt.Errorf("'data' is nil")
	}
	data.OrgId = orgId
	output = *data
	err = db.Model(&output).Updates(output).Error
	if err != nil {
		return model.Domain{}, err
	}
	return output, nil
}

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
	if output.Type != nil && *output.Type == model.DomainTypeIpa {
		output.IpaDomain = &model.Ipa{}
		if err = db.Preload("CaCerts").Preload("Servers").First(output.IpaDomain, "id = ?", output.ID).Count(&count).Error; err != nil {
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
		return fmt.Errorf("'db' is nil")
	}
	if orgId == "" {
		return fmt.Errorf("'orgId' cannot be an empty string")
	}
	if uuid == "" {
		return fmt.Errorf("'uuid' cannot be an empty string")
	}
	if err = db.First(&data, "org_id = ? AND domain_uuid = ?", orgId, uuid).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Register not found")
	}
	if err = db.Unscoped().Delete(&data, "org_id = ? AND domain_uuid = ?", orgId, uuid).Error; err != nil {
		return err
	}
	return nil
}

func (r *domainRepository) RhelIdmClearToken(db *gorm.DB, orgId string, uuid string) (err error) {
	if db == nil {
		return fmt.Errorf("'db' is nil")
	}
	if orgId == "" {
		return fmt.Errorf("'orgId' is empty")
	}
	if uuid == "" {
		return fmt.Errorf("'uuid' is empty")
	}
	var dataDomain model.Domain
	if dataDomain, err = r.FindById(db, orgId, uuid); err != nil {
		return err
	}
	if dataDomain.OrgId != orgId {
		return fmt.Errorf("'OrgId' mistmatch")
	}
	switch *dataDomain.Type {
	case model.DomainTypeIpa:
		if err = db.Table("ipas").
			Where("id = ?", dataDomain.IpaDomain.Model.ID).
			Update("token", nil).
			Error; err != nil {
			return err
		}
		if err = db.Table("ipas").
			Where("id = ?", dataDomain.IpaDomain.Model.ID).
			Update("token_expiration", nil).
			Error; err != nil {
			return err
		}

	default:
		return fmt.Errorf("'Type=%d' invalid", *dataDomain.Type)
	}

	return nil
}
