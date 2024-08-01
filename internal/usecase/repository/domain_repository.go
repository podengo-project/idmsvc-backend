package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/domain_token"
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
// ctx is the current request context with db and slog instances.
// orgID is the organization id that we belongs.
// offset is the starting record for the given ordered result.
// limit is the number of items for the current requested page.
// Return the list of items for the current page, the total number
// of items for the given organization and nil error if the call
// is successful, else it return nil slice, 0 for count and a
// filled error interface with the details.
func (r *domainRepository) List(
	ctx context.Context,
	orgID string,
	offset int,
	limit int,
) (output []model.Domain, count int64, err error) {
	log, db, err := r.checkList(ctx, orgID, offset, limit)
	if err != nil {
		log.ErrorContext(ctx, err.Error())
		return nil, 0, err
	}
	if err = db.
		Table("domains").
		Where("org_id = ?", orgID).
		Count(&count).
		Offset(int(offset)).
		Limit(limit).
		Find(&output).
		Error; err != nil {
		log.ErrorContext(ctx, err.Error())
		return nil, 0, err
	}

	return output, count, nil
}

// Register a new domain
func (r *domainRepository) Register(
	ctx context.Context,
	orgID string,
	data *model.Domain,
) (err error) {
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if err = r.checkCommonAndData(db, orgID, data); err != nil {
		log.Error(err.Error())
		return err
	}

	if err = db.Omit(clause.Associations).
		Create(data).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			msg := fmt.Sprintf("domain id '%s' is already registered.", data.DomainUuid.String())
			err = echo.NewHTTPError(http.StatusConflict, msg)
			log.Error(err.Error())
			return err
		} else {
			log.Error(err.Error())
			return err
		}
	}

	// Specific
	switch *data.Type {
	case model.DomainTypeIpa:
		if data.IpaDomain == nil {
			err = internal_errors.NilArgError("IpaDomain")
			log.Error(err.Error())
			return err
		}
		if err = r.createIpaDomain(
			log,
			db,
			data.ID,
			data.IpaDomain,
		); err != nil {
			log.Error(err.Error())
			return err
		}
		return nil
	default:
		err = fmt.Errorf("'Type' is invalid")
		log.Error(err.Error())
		return err
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

// UpdateAgent save the Domain record into the database. It only update
// data for the current organization.
func (r *domainRepository) UpdateAgent(
	ctx context.Context,
	orgID string,
	data *model.Domain,
) (err error) {
	var currentDomain *model.Domain
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if err = r.checkCommonAndData(db, orgID, data); err != nil {
		log.Error(err.Error())
		return err
	}

	// Check the entity exists
	if currentDomain, err = r.FindByID(
		ctx,
		orgID,
		data.DomainUuid,
	); err != nil {
		log.Error(err.Error())
		return err
	}

	if err = db.Omit(clause.Associations).
		Where("org_id = ? AND domain_uuid = ?", orgID, currentDomain.DomainUuid).
		Updates(data).
		Error; err != nil {
		log.Error(err.Error())
		return err
	}

	// Specific
	switch *data.Type {
	case model.DomainTypeIpa:
		if data.IpaDomain == nil {
			err = internal_errors.NilArgError("IpaDomain")
			log.Error(err.Error())
			return err
		}
		data.IpaDomain.ID = data.ID
		return r.updateIpaDomain(log, db, data.IpaDomain)
	default:
		err = fmt.Errorf("'Type' is invalid")
		log.Error(err.Error())
		return err
	}
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
	ctx context.Context,
	orgID string,
	data *model.Domain,
) (err error) {
	var currentDomain *model.Domain
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if err = r.checkCommonAndData(db, orgID, data); err != nil {
		log.Error(err.Error())
		return err
	}

	// Check the entity exists
	if currentDomain, err = r.FindByID(
		ctx,
		orgID,
		data.DomainUuid,
	); err != nil {
		log.Error(err.Error())
		return err
	}

	fields := r.prepareUpdateUser(data)

	if err = db.Omit(clause.Associations).
		Model(data).
		Where("org_id = ? AND domain_uuid = ?", orgID, currentDomain.DomainUuid).
		Updates(fields).
		Error; err != nil {
		log.Error(err.Error())
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
	ctx context.Context,
	orgID string,
	UUID uuid.UUID,
) (output *model.Domain, err error) {
	// See: https://gorm.io/docs/query.html
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if err = r.checkCommonAndUUID(db, orgID, UUID); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	output = &model.Domain{}
	if err = db.Model(&model.Domain{}).
		First(output, "org_id = ? AND domain_uuid = ?", orgID, UUID).
		Error; err != nil {
		err = r.wrapErrNotFound(err, UUID)
		log.Error(err.Error())
		return nil, err
	}
	if err = output.FillAndPreload(db); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return output, nil
}

// See: https://gorm.io/docs/delete.html
func (r *domainRepository) DeleteById(
	ctx context.Context,
	orgID string,
	UUID uuid.UUID,
) (err error) {
	var (
		data  model.Domain
		count int64
	)
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if err = r.checkCommonAndUUID(db, orgID, UUID); err != nil {
		log.Error(err.Error())
		return err
	}
	if err = db.First(&data, "org_id = ? AND domain_uuid = ?", orgID, UUID).Count(&count).Error; err != nil {
		err = r.wrapErrNotFound(err, UUID)
		log.Error("deleting domain when checking that the record exist")
		return err
	}
	if count == 0 {
		err = r.wrapErrNotFound(gorm.ErrRecordNotFound, UUID)
		log.Error("deleting domain because no record found to delete")
		return err
	}
	if err = db.Unscoped().Delete(&data, "org_id = ? AND domain_uuid = ?", orgID, UUID).Error; err != nil {
		err = r.wrapErrNotFound(err, UUID)
		log.Error("deleting domain when removing record")
		return err
	}
	return nil
}

func (r *domainRepository) CreateDomainToken(
	ctx context.Context,
	key []byte,
	validity time.Duration,
	orgID string,
	domainType public.DomainType,
) (drt *repository.DomainRegToken, err error) {
	log := app_context.LogFromCtx(ctx)
	tok, expireNS, err := domain_token.NewDomainRegistrationToken(key, string(domainType), orgID, validity)
	if err != nil {
		log.Error("creating ipa domain token")
		return nil, err
	}
	domainId := domain_token.TokenDomainId(tok)
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
	log *slog.Logger,
	db *gorm.DB,
	domainID uint,
	data *model.Ipa,
) (err error) {
	if data == nil {
		err = internal_errors.NilArgError("data")
		log.Error(err.Error())
		return err
	}
	data.Model.ID = domainID
	if err = db.Omit(clause.Associations).Create(data).Error; err != nil {
		log.Error("creating ipa record")
		return err
	}
	for idx := range data.CaCerts {
		data.CaCerts[idx].IpaID = data.ID
		if err = db.Create(&data.CaCerts[idx]).Error; err != nil {
			log.Error("creating ipa cert record")
			return err
		}
	}
	for idx := range data.Servers {
		data.Servers[idx].IpaID = data.ID
		if err = db.Create(&data.Servers[idx]).Error; err != nil {
			log.Error("creating ipa server record")
			return err
		}
	}
	for idx := range data.Locations {
		data.Locations[idx].IpaID = data.ID
		if err = db.Create(&data.Locations[idx]).Error; err != nil {
			log.Error("creating ipa location record")
			return err
		}
	}
	return nil
}

func (r *domainRepository) updateIpaDomain(
	log *slog.Logger,
	db *gorm.DB,
	data *model.Ipa,
) (err error) {
	if db == nil {
		err = internal_errors.NilArgError("db")
		log.Error(err.Error())
		return err
	}
	if data == nil {
		err = internal_errors.NilArgError("data")
		log.Error(err.Error())
		return err
	}
	if err = db.Unscoped().
		Delete(data).Error; err != nil {
		log.Error("updating ipa domain when deleting old record")
		return err
	}

	if err = db.Omit(clause.Associations).
		Create(data).
		Error; err != nil {
		log.Error("updating ipa domain when creating new ipa record")
		return err
	}

	// CaCerts
	for i := range data.CaCerts {
		data.CaCerts[i].IpaID = data.ID
		if err = db.Create(&data.CaCerts[i]).Error; err != nil {
			log.Error("updating ipa domain when creating new ipa certificate record")
			return err
		}
	}

	// Servers
	for i := range data.Servers {
		data.Servers[i].IpaID = data.ID
		if err = db.Create(&data.Servers[i]).Error; err != nil {
			log.Error("updating ipa domain when creating new ipa server record")
			return err
		}
	}

	// Locations
	for i := range data.Locations {
		data.Locations[i].IpaID = data.ID
		if err = db.Create(&data.Locations[i]).Error; err != nil {
			log.Error("updating ipa domain when creating new ipa location record")
			return err
		}
	}

	return nil
}

func (r *domainRepository) wrapErrNotFound(err error, UUID uuid.UUID) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return internal_errors.NewHTTPErrorF(
			http.StatusNotFound,
			"unknown domain '%s'",
			UUID.String(),
		)
	} else {
		return err
	}
}

func (r *domainRepository) checkList(
	ctx context.Context,
	orgID string,
	offset int,
	limit int,
) (*slog.Logger, *gorm.DB, error) {
	log := app_context.LogFromCtx(ctx)
	db := app_context.DBFromCtx(ctx)
	if orgID == "" {
		return log, db, fmt.Errorf("'orgID' is empty")
	}
	if offset < 0 {
		return log, db, fmt.Errorf("'offset' is lower than 0")
	}
	if limit < 0 {
		return log, db, fmt.Errorf("'limit' is lower than 0")
	}
	return log, db, nil
}
