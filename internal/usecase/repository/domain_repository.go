package repository

import (
	"fmt"
	"time"

	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type domainRepository struct{}

// NewDomainRepository create a new repository component for the
// idm domainds backend service.
// Return a repository.DomainRepository interface.
func NewDomainRepository() repository.DomainRepository {
	return &domainRepository{}
}

// List retrieve the list of domains for the given orgID and
// the pagination info.
// db is the gorm database connector.
// orgID is the organization id that we belongs.
// offset is the starting record for the given ordered result.
// limit is the number of items for the current requested page.
// Return the list of items for the current page, the total number
// of items for the given organization and nil error if the call
// is successful, else it return nil slice, 0 for count and a
// filled error interface with the details.
func (r *domainRepository) List(
	db *gorm.DB,
	orgID string,
	offset int,
	limit int,
) (output []model.Domain, count int64, err error) {
	if db == nil {
		return nil, 0, fmt.Errorf("'db' is nil")
	}
	if orgID == "" {
		return nil, 0, fmt.Errorf("'orgID' is empty")
	}
	if offset < 0 {
		return nil, 0, fmt.Errorf("'offset' is lower than 0")
	}
	if limit < 0 {
		return nil, 0, fmt.Errorf("'limit' is lower than 0")
	}
	if err = db.
		Table("domains").
		Limit(limit).
		Where("org_id = ?", orgID).
		Offset(int(offset)).
		Count(&count).
		Find(&output).
		Error; err != nil {
		return nil, 0, err
	}

	return output, count, nil
}

// Create add a new record to the domains entity.
// db is the gorm database connector.
// orgID is the organization to segment the data.
// data is the business model to be created into the database.
// Return nil for a successful operation, or an error interface
// with details about the situation.
func (r *domainRepository) Create(
	db *gorm.DB,
	orgID string,
	data *model.Domain,
) (err error) {
	if err = r.checkCommonAndDataAndType(
		db,
		orgID,
		data,
	); err != nil {
		return err
	}
	data.OrgId = orgID
	if err = db.Omit(clause.Associations).
		Create(data).Error; err != nil {
		return err
	}
	switch *data.Type {
	case model.DomainTypeIpa:
		if data.IpaDomain == nil {
			return nil
		}
		if err = r.createIpaDomain(
			db,
			data.ID,
			data.IpaDomain,
		); err != nil {
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
	orgID string,
	data *model.Domain,
) (err error) {
	var currentDomain *model.Domain
	if err = r.checkCommonAndData(db, orgID, data); err != nil {
		return err
	}

	uuid := data.DomainUuid.String()
	// Check the entity exists
	if currentDomain, err = r.FindByID(
		db,
		orgID,
		uuid,
	); err != nil {
		return err
	}

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
		Where("org_id = ? AND domain_uuid = ?", orgID, currentDomain.DomainUuid).
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

// FindByID retrieve the model.Domain specified by its uuid that
// belongs to the specified organization.
// db is the gorm database connector.
// orgID is the organization id which the statement is executed for.
// uuid is the uuid that identify the domain record as provided for
// the API.
// Return a filled model.Domain structure and nil error for success
// scenario, else nil reference and an error with details about the
// situation.
func (r *domainRepository) FindByID(
	db *gorm.DB,
	orgID string,
	uuid string,
) (output *model.Domain, err error) {
	// See: https://gorm.io/docs/query.html
	if err = r.checkCommonAndUUID(db, orgID, uuid); err != nil {
		return nil, err
	}
	if err = db.Model(&model.Domain{}).
		First(&output, "org_id = ? AND domain_uuid = ?", orgID, uuid).
		Error; err != nil {
		return nil, err
	}
	if err = output.FillAndPreload(db); err != nil {
		return nil, err
	}
	return output, nil
}

// See: https://gorm.io/docs/delete.html
func (r *domainRepository) DeleteById(
	db *gorm.DB,
	orgID string,
	uuid string,
) (err error) {
	var (
		data  model.Domain
		count int64
	)
	if err = r.checkCommonAndUUID(db, orgID, uuid); err != nil {
		return err
	}
	if err = db.First(&data, "org_id = ? AND domain_uuid = ?", orgID, uuid).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Register not found")
	}
	if err = db.Unscoped().Delete(&data, "org_id = ? AND domain_uuid = ?", orgID, uuid).Error; err != nil {
		return err
	}
	return nil
}

func (r *domainRepository) RhelIdmClearToken(
	db *gorm.DB,
	orgID string,
	uuid string,
) (err error) {
	if err = r.checkCommonAndUUID(db, orgID, uuid); err != nil {
		return err
	}
	var dataDomain *model.Domain
	if dataDomain, err = r.FindByID(db, orgID, uuid); err != nil {
		return err
	}
	if dataDomain.OrgId != orgID {
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
			Update("token_expiration_ts", nil).
			Error; err != nil {
			return err
		}
	}

	return nil
}

// ------- PRIVATE METHODS --------

func (r *domainRepository) checkCommon(
	db *gorm.DB,
	orgID string,
) error {
	if db == nil {
		return fmt.Errorf("'db' is nil")
	}
	if orgID == "" {
		return fmt.Errorf("'orgID' is empty")
	}
	return nil
}

func (r *domainRepository) checkCommonAndUUID(
	db *gorm.DB,
	orgID string,
	uuid string,
) error {
	if err := r.checkCommon(db, orgID); err != nil {
		return err
	}
	if uuid == "" {
		return fmt.Errorf("'uuid' is empty")
	}
	return nil
}

func (r *domainRepository) checkCommonAndData(
	db *gorm.DB,
	orgID string,
	data *model.Domain,
) error {
	if err := r.checkCommon(db, orgID); err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("'data' is nil")
	}
	return nil
}

func (r *domainRepository) checkCommonAndDataAndType(
	db *gorm.DB,
	orgID string,
	data *model.Domain,
) error {
	if err := r.checkCommonAndData(db, orgID, data); err != nil {
		return err
	}
	if data.Type == nil {
		return fmt.Errorf("'Type' is nil")
	}
	return nil
}

func (r *domainRepository) createIpaDomain(
	db *gorm.DB,
	domainID uint,
	data *model.Ipa,
) (err error) {
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
	for idx := range data.Locations {
		data.Servers[idx].IpaID = data.ID
		if err = db.Create(&data.Locations[idx]).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *domainRepository) updateIpaDomain(
	db *gorm.DB,
	data *model.Ipa,
) (err error) {
	if db == nil {
		return fmt.Errorf("'db' is nil")
	}
	if data == nil {
		return fmt.Errorf("'data' of type '*model.Ipa' is nil")
	}
	if err = db.Unscoped().
		Delete(data).Error; err != nil {
		return err
	}

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
	for i := range data.Servers {
		data.Servers[i].IpaID = data.ID
		if err = db.Create(&data.Servers[i]).Error; err != nil {
			return err
		}
	}

	// Locations
	for i := range data.Locations {
		data.Locations[i].IpaID = data.ID
		if err = db.Create(&data.Locations[i]).Error; err != nil {
			return err
		}
	}

	return nil
}
