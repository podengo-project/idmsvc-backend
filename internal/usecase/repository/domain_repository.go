package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/repository"
	"github.com/openlyinc/pointy"
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
	if err = r.checkCommonAndDataAndType(db, orgId, data); err != nil {
		return err
	}
	data.OrgId = orgId
	if err = db.Omit(clause.Associations).
		Create(data).Error; err != nil {
		return err
	}
	switch *data.Type {
	case model.DomainTypeIpa:
		if data.IpaDomain == nil {
			return nil
		}
		if err = r.createIpaDomain(db, data.ID, data.IpaDomain); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("'Type' is invalid")
	}
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
func (r *domainRepository) Update(
	db *gorm.DB,
	orgId string,
	data *model.Domain,
) (err error) {
	var currentDomain *model.Domain
	if err = r.checkCommonAndData(db, orgId, data); err != nil {
		return err
	}

	uuid := data.DomainUuid.String()
	// Check the entity exists
	if currentDomain, err = r.FindById(
		db,
		orgId,
		uuid,
	); err != nil {
		return err
	}

	// gorm.Model
	// OrgId                 string
	// DomainUuid            uuid.UUID `gorm:"unique"`
	// DomainName            *string
	// Title                 *string
	// Description           *string
	// Type                  *uint
	// AutoEnrollmentEnabled *bool

	currentDomain.OrgId = orgId
	currentDomain.DomainUuid = data.DomainUuid

	if data.DomainName != nil {
		currentDomain.DomainName = data.DomainName
	} else {
		currentDomain.DomainName = pointy.String("")
	}

	if data.Title != nil {
		currentDomain.Title = data.Title
	} else {
		currentDomain.Title = pointy.String("")
	}

	if data.Description != nil {
		currentDomain.Description = data.Description
	} else {
		currentDomain.Description = pointy.String("")
	}

	if data.AutoEnrollmentEnabled != nil {
		currentDomain.AutoEnrollmentEnabled = data.AutoEnrollmentEnabled
	} else {
		currentDomain.AutoEnrollmentEnabled = pointy.Bool(false)
	}

	if err = db.Omit(clause.Associations).
		Where("org_id = ?", orgId).
		Updates(data).
		Error; err != nil {
		return err
	}

	// Specific
	switch *data.Type {
	case model.DomainTypeIpa:
		if data.IpaDomain == nil {
			return fmt.Errorf("'IpaDomain' is nil")
		}
		data.IpaDomain.ID = data.ID
		return r.updateIpaDomain(db, data.IpaDomain)
	default:
		return fmt.Errorf("'Type' is invalid")
	}
}

func (r *domainRepository) updateIpaDomain(db *gorm.DB, data *model.Ipa) (err error) {
	if data == nil {
		return fmt.Errorf("'data' of type '*model.Ipa' is nil")
	}
	if err = db.Unscoped().
		Delete(
			data,
			"id = ?",
			data.ID,
		).Error; err != nil {
		return err
	}

	// Clean token and create new entry
	data.Model.CreatedAt = time.Time{}
	data.Model.UpdatedAt = time.Time{}
	data.Model.DeletedAt = gorm.DeletedAt{}
	data.Token = nil
	data.TokenExpiration = nil
	if err = db.Omit(clause.Associations).
		Create(data).
		Error; err != nil {
		return err
	}

	// CaCerts
	for i := range data.CaCerts {
		data.CaCerts[i].IpaID = data.ID
		if err = db.Create(&data.CaCerts[i]).Error; err != nil {
			return err
		}
	}

	// Servers
	for i := range data.CaCerts {
		data.Servers[i].IpaID = data.ID
		if err = db.Create(&data.Servers[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

// See: https://gorm.io/docs/query.html
// TODO Document the method
func (r *domainRepository) FindById(db *gorm.DB, orgId string, uuid string) (output *model.Domain, err error) {
	if err = r.checkCommonAndUUID(db, orgId, uuid); err != nil {
		return nil, err
	}
	if err = db.First(&output, "org_id = ? AND domain_uuid = ?", orgId, uuid).Error; err != nil {
		return nil, err
	}
	if output.Type == nil {
		return output, nil
	}
	switch *output.Type {
	case model.DomainTypeIpa:
		output.IpaDomain = &model.Ipa{}
		output.IpaDomain.ID = output.ID
		if err = db.Preload("CaCerts").
			Preload("Servers").
			// First(output.IpaDomain, "id = ?", output.ID).
			First(output.IpaDomain).
			Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}
		return output, nil
	default:
		return nil, fmt.Errorf("'Type' is invalid")
	}
}

// See: https://gorm.io/docs/delete.html
func (r *domainRepository) DeleteById(db *gorm.DB, orgId string, uuid string) (err error) {
	var (
		data  model.Domain
		count int64
	)
	if err = r.checkCommonAndUUID(db, orgId, uuid); err != nil {
		return err
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
	if err = r.checkCommonAndUUID(db, orgId, uuid); err != nil {
		return err
	}
	var dataDomain *model.Domain
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

// ------- PRIVATE METHODS --------
func (r *domainRepository) checkCommon(db *gorm.DB, orgId string) error {
	if db == nil {
		return fmt.Errorf("'db' is nil")
	}
	if orgId == "" {
		return fmt.Errorf("'orgId' is empty")
	}
	return nil
}

func (r *domainRepository) checkCommonAndUUID(db *gorm.DB, orgId string, uuid string) error {
	if err := r.checkCommon(db, orgId); err != nil {
		return err
	}
	if uuid == "" {
		return fmt.Errorf("'uuid' is empty")
	}
	return nil
}

func (r *domainRepository) checkCommonAndData(db *gorm.DB, orgId string, data *model.Domain) error {
	if err := r.checkCommon(db, orgId); err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("'data' is nil")
	}
	return nil
}

func (r *domainRepository) checkCommonAndDataAndType(db *gorm.DB, orgId string, data *model.Domain) error {
	if err := r.checkCommonAndData(db, orgId, data); err != nil {
		return err
	}
	if data.Type == nil {
		return fmt.Errorf("'Type' is nil")
	}
	return nil
}
