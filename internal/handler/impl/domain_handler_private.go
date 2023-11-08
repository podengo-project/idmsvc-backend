package impl

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"gorm.io/gorm"
)

func (a *application) findIpaById(tx *gorm.DB, orgId string, UUID uuid.UUID) (data *model.Domain, err error) {
	if data, err = a.domain.repository.FindByID(tx, orgId, UUID); err != nil {
		return nil, err
	}
	if *data.Type != model.DomainTypeIpa {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Wrong domain type")
	}
	if data.IpaDomain == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "No IPA data found for the domain")
	}
	return data, nil
}

func ensureSubscriptionManagerIDAuthorizedToUpdate(
	subscriptionManagerID string,
	servers []model.IpaServer,
) error {
	if subscriptionManagerID == "" {
		return fmt.Errorf("'subscriptionManagerID' is empty")
	}
	if servers == nil {
		return internal_errors.NilArgError("servers")
	}
	for i := range servers {
		if servers[i].HCCUpdateServer &&
			servers[i].RHSMId != nil &&
			*servers[i].RHSMId == subscriptionManagerID {
			return nil
		}
	}
	return internal_errors.NewHTTPErrorF(
		http.StatusForbidden,
		"'subscriptionManagerID' not found into the authorized list of rhel-idm servers",
	)
}

// fillDomain is a helper function to copy Ipa domain
// data between structures, to be used at register IPA domain endpoint.
// target is the destination Ipa structure, it cannot be nil.
// source is the source Ipa structure, it cannot be nil.
// Return nil if it is copied succesfully, else an error.
func (a *application) fillDomain(
	target *model.Domain,
	source *model.Domain,
) error {
	if source == nil {
		return fmt.Errorf("'target' cannot be nil")
	}
	if target == nil {
		return fmt.Errorf("'source' cannot be nil")
	}
	if source.Type == nil {
		return internal_errors.NilArgError("Type")
	}
	if source.AutoEnrollmentEnabled != nil {
		target.AutoEnrollmentEnabled = pointy.Bool(*source.AutoEnrollmentEnabled)
	}
	if source.DomainName != nil {
		target.DomainName = pointy.String(*source.DomainName)
	}
	if source.Title != nil {
		target.Title = pointy.String(*source.Title)
	}
	if source.Description != nil {
		target.Description = pointy.String(*source.Description)
	}
	target.OrgId = source.OrgId
	if source.Type != nil {
		target.Type = pointy.Uint(*source.Type)
	}

	switch *target.Type {
	case model.DomainTypeIpa:
		if source.IpaDomain == nil {
			return internal_errors.NilArgError("source.IpaDomain")
		}
		target.IpaDomain = &model.Ipa{}
		return a.fillDomainIpa(target.IpaDomain, source.IpaDomain)
	default:
		return fmt.Errorf("'model.DomainTypeIpa' ")
	}
}

// fillDomainUser is a helper function to copy domain
// data between structures when patching the domain information.
// target is the destination Ipa structure, it cannot be nil.
// source is the source Ipa structure, it cannot be nil.
// Return nil if it is copied succesfully, else an error.
func (a *application) fillDomainUser(
	target *model.Domain,
	source *model.Domain,
) error {
	if source == nil {
		return fmt.Errorf("'target' cannot be nil")
	}
	if target == nil {
		return fmt.Errorf("'source' cannot be nil")
	}
	if source.AutoEnrollmentEnabled != nil {
		target.AutoEnrollmentEnabled = pointy.Bool(*source.AutoEnrollmentEnabled)
	}
	if source.Title != nil {
		target.Title = pointy.String(*source.Title)
	}
	if source.Description != nil {
		target.Description = pointy.String(*source.Description)
	}
	return nil
}

func (a *application) fillDomainIpa(target *model.Ipa, source *model.Ipa) error {
	target.Model = source.Model
	if source.RealmName != nil {
		target.RealmName = pointy.String(*source.RealmName)
	}
	target.CaCerts = make([]model.IpaCert, len(source.CaCerts))
	for i := range source.CaCerts {
		target.CaCerts[i] = source.CaCerts[i]
		target.CaCerts[i].IpaID = target.ID
	}
	target.Servers = make([]model.IpaServer, len(source.Servers))
	for i := range source.Servers {
		target.Servers[i] = source.Servers[i]
		target.Servers[i].IpaID = target.ID
	}
	target.Locations = make([]model.IpaLocation, len(source.Locations))
	for i := range source.Locations {
		target.Locations[i] = source.Locations[i]
		target.Locations[i].IpaID = target.ID
	}
	target.RealmDomains = source.RealmDomains
	return nil
}
