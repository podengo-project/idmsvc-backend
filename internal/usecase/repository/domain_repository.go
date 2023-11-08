package repository

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token"
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
		return nil, 0, internal_errors.NilArgError("db")
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

// Register a new domain
func (r *domainRepository) Register(
	db *gorm.DB,
	orgID string,
	data *model.Domain,
) (err error) {
	if err = r.checkCommonAndData(db, orgID, data); err != nil {
		return err
	}

	if err = db.Omit(clause.Associations).
		Create(data).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			msg := fmt.Sprintf("domain id '%s' is already registered.", data.DomainUuid.String())
			return echo.NewHTTPError(http.StatusConflict, msg)
		} else {
			return err
		}
	}

	// Specific
	switch *data.Type {
	case model.DomainTypeIpa:
		if data.IpaDomain == nil {
			return internal_errors.NilArgError("IpaDomain")
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

// UpdateAgent save the Domain record into the database.  The org_id
// and domain_uuid from the data object are used to select the target
// record.
//
// It is assumed that the record is complete and changes have already been
// "merged" and are ready to save to the database.
func (r *domainRepository) UpdateAgent(
	db *gorm.DB,
	orgID string,
	oldData *model.Domain,
	newData *model.Domain,
) (err error) {
	if err = r.checkCommonAndDataUpdateAgent(db, orgID, oldData, newData); err != nil {
		return err
	}

	var (
		remCerts     []uint
		remServers   []uint
		remLocations []uint
	)

	oldData.OrgId = orgID
	newData.OrgId = orgID
	if remCerts, remServers, remLocations, err = r.ipaFillUpdateAgent(oldData, newData); err != nil {
		return err
	}

	if err = db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Where("org_id = ? AND id = ?", orgID, newData.ID).
		Save(newData).
		Error; err != nil {
		return err
	}

	if len(remCerts) > 0 {
		if err := db.
			Unscoped().
			Model(&model.IpaCert{}).
			Delete(&model.IpaCert{}, "ipa_certs.id in ?", remCerts).
			Error; err != nil {
			return err
		}
	}

	if len(remServers) > 0 {
		if err := db.
			Unscoped().
			Model(&model.IpaServer{}).
			Delete(&model.IpaServer{}, "ipa_servers.id in ?", remServers).
			Error; err != nil {
			return err
		}
	}

	if len(remLocations) > 0 {
		if err := db.
			Unscoped().
			Model(&model.IpaLocation{}).
			Delete(&model.IpaLocation{}, "ipa_locations.id in ?", remLocations).
			Error; err != nil {
			return err
		}
	}

	return nil
}

// prepareUpdateUser fill the hashmap with the not nil
// fields that we need to update into the database. It
// only consider the fields that a user can update.
// data is the partial fields to be updated.
// Return a hashmap with the values of the fields.
func (r *domainRepository) prepareUpdateUser(data *model.Domain) map[string]interface{} {
	res := make(map[string]interface{}, 3)
	if data.Title != nil {
		res["title"] = data.Title
	}
	if data.Description != nil {
		res["description"] = data.Description
	}
	if data.AutoEnrollmentEnabled != nil {
		res["auto_enrollment_enabled"] = data.AutoEnrollmentEnabled
	}
	return res
}

// UpdateUser save the Domain record, but only the provided
// information for the user update.
// data for the current organization.
func (r *domainRepository) UpdateUser(
	db *gorm.DB,
	orgID string,
	data *model.Domain,
) (err error) {
	var currentDomain *model.Domain
	if err = r.checkCommonAndData(db, orgID, data); err != nil {
		return err
	}

	// Check the entity exists
	if currentDomain, err = r.FindByID(
		db,
		orgID,
		data.DomainUuid,
	); err != nil {
		return err
	}

	fields := r.prepareUpdateUser(data)

	if err = db.Omit(clause.Associations).
		Model(data).
		Where("org_id = ? AND domain_uuid = ?", orgID, currentDomain.DomainUuid).
		Updates(fields).
		Error; err != nil {
		return err
	}

	return nil
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
	UUID uuid.UUID,
) (output *model.Domain, err error) {
	// See: https://gorm.io/docs/query.html
	if err = r.checkCommonAndUUID(db, orgID, UUID); err != nil {
		return nil, err
	}
	output = &model.Domain{}
	if err = db.Model(&model.Domain{}).
		First(output, "org_id = ? AND domain_uuid = ?", orgID, UUID).
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
	UUID uuid.UUID,
) (err error) {
	var (
		data  model.Domain
		count int64
	)
	if err = r.checkCommonAndUUID(db, orgID, UUID); err != nil {
		return err
	}
	if err = db.First(&data, "org_id = ? AND domain_uuid = ?", orgID, UUID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Register not found")
	}
	if err = db.Unscoped().Delete(&data, "org_id = ? AND domain_uuid = ?", orgID, UUID).Error; err != nil {
		return err
	}
	return nil
}

func (r *domainRepository) CreateDomainToken(
	key []byte,
	validity time.Duration,
	orgID string,
	domainType public.DomainType,
) (drt *repository.DomainRegToken, err error) {
	tok, expireNS, err := token.NewDomainRegistrationToken(key, string(domainType), orgID, validity)
	if err != nil {
		return nil, err
	}
	domainId := token.TokenDomainId(tok)
	drt = &repository.DomainRegToken{
		DomainId:     domainId,
		DomainToken:  string(tok),
		DomainType:   domainType,
		ExpirationNS: expireNS,
	}
	return drt, nil
}

// ------- PRIVATE METHODS --------

func (r *domainRepository) checkCommon(
	db *gorm.DB,
	orgID string,
) error {
	if db == nil {
		return internal_errors.NilArgError("db")
	}
	if orgID == "" {
		return fmt.Errorf("'orgID' is empty")
	}
	return nil
}

func (r *domainRepository) checkCommonAndUUID(
	db *gorm.DB,
	orgID string,
	UUID uuid.UUID,
) error {
	if err := r.checkCommon(db, orgID); err != nil {
		return err
	}
	if UUID == uuid.Nil {
		return fmt.Errorf("'uuid' is invalid")
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
		return internal_errors.NilArgError("data")
	}
	return nil
}

func (r *domainRepository) checkCommonAndDataUpdateAgent(
	db *gorm.DB,
	orgID string,
	oldData *model.Domain,
	newData *model.Domain,
) error {
	if err := r.checkCommon(db, orgID); err != nil {
		return err
	}
	if oldData == nil {
		return internal_errors.NilArgError("oldData")
	}
	if newData == nil {
		return internal_errors.NilArgError("newData")
	}
	if orgID != oldData.OrgId || orgID != newData.OrgId {
		return internal_errors.NewHTTPErrorF(http.StatusInternalServerError, "orgID mismatch")
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
		return internal_errors.NilArgError("Type")
	}
	return nil
}

func (r *domainRepository) createIpaDomain(
	db *gorm.DB,
	domainID uint,
	data *model.Ipa,
) (err error) {
	if data == nil {
		return internal_errors.NilArgError("data' of type '*model.Ipa")
	}
	data.Model.ID = domainID
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
		data.Locations[idx].IpaID = data.ID
		if err = db.Create(&data.Locations[idx]).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *domainRepository) ipaFillUpdateAgentRemainders(oldData *model.Domain, newData *model.Domain) (remCerts []uint, remServers []uint, remLocations []uint) {
	remCerts = nil
	remServers = nil
	remLocations = nil

	// Check cert remainders
	OldLenCerts := len(oldData.IpaDomain.CaCerts)
	NewLenCerts := len(newData.IpaDomain.CaCerts)
	if OldLenCerts > NewLenCerts {
		remCerts = make([]uint, OldLenCerts-NewLenCerts)
		for i := NewLenCerts; i < OldLenCerts; i++ {
			remCerts[i-NewLenCerts] = oldData.IpaDomain.CaCerts[i].Model.ID
		}
	}

	// Check servers remainders
	OldLenServers := len(oldData.IpaDomain.Servers)
	NewLenServers := len(newData.IpaDomain.Servers)
	if OldLenServers > NewLenServers {
		remServers = make([]uint, OldLenServers-NewLenServers)
		for i := NewLenServers; i < OldLenServers; i++ {
			remServers[i-NewLenServers] = oldData.IpaDomain.Servers[i].Model.ID
		}
	}

	// Check location remainders
	OldLenLocations := len(oldData.IpaDomain.Locations)
	NewLenLocations := len(newData.IpaDomain.Locations)
	if OldLenLocations > NewLenLocations {
		remLocations = make([]uint, OldLenLocations-NewLenLocations)
		for i := NewLenLocations; i < OldLenLocations; i++ {
			remLocations[i-NewLenLocations] = oldData.IpaDomain.Locations[i].Model.ID
		}
	}
	return
}

func (r *domainRepository) ipaFillUpdateAgent(oldData *model.Domain, newData *model.Domain) (remCerts []uint, remServers []uint, remLocations []uint, err error) {
	remCerts, remServers, remLocations = r.ipaFillUpdateAgentRemainders(oldData, newData)

	newData.Model = oldData.Model
	newData.IpaDomain.Model = oldData.IpaDomain.Model

	newData.DomainName = oldData.DomainName
	newData.Title = oldData.Title
	newData.Description = oldData.Description
	newData.AutoEnrollmentEnabled = oldData.AutoEnrollmentEnabled

	// Certs
	for i := range newData.IpaDomain.CaCerts {
		newData.IpaDomain.CaCerts[i].Model = oldData.IpaDomain.CaCerts[i].Model
	}

	// Servers
	for i := range newData.IpaDomain.Servers {
		newData.IpaDomain.Servers[i].Model = oldData.IpaDomain.Servers[i].Model
	}

	// Locations
	for i := range newData.IpaDomain.Locations {
		newData.IpaDomain.Locations[i].Model = oldData.IpaDomain.Locations[i].Model
	}

	return remCerts, remServers, remLocations, nil
}
