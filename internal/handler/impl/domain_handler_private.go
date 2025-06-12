package impl

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

func (a *application) findIpaById(tx *gorm.DB, orgId string, UUID uuid.UUID) (data *model.Domain, err error) {
	ctx := app_context.CtxWithDB(context.Background(), tx)
	if data, err = a.domain.repository.FindByID(ctx, orgId, UUID); err != nil {
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

func subscriptionManagerIDIncluded(
	ctx context.Context,
	subscriptionManagerID string,
	servers []model.IpaServer,
) (bool, error) {
	logger := app_context.LogFromCtx(ctx)
	if subscriptionManagerID == "" {
		logger.Error("'subscriptionManagerID' is an empty string")
		return false, fmt.Errorf("'subscriptionManagerID' is empty")
	}
	if servers == nil {
		logger.Error("'servers' is nil")
		return false, internal_errors.NilArgError("servers")
	}
	for i := range servers {
		rhsmid := "nil"
		if servers[i].RHSMId != nil {
			rhsmid = *servers[i].RHSMId
		}

		logger.Debug("Checking server",
			slog.Bool("HCCUpdateServer", servers[i].HCCUpdateServer),
			slog.String("RHSMId", rhsmid),
			slog.String("subscriptionManagerID", subscriptionManagerID),
		)
		if servers[i].HCCUpdateServer &&
			servers[i].RHSMId != nil &&
			*servers[i].RHSMId == subscriptionManagerID {
			logger.Debug("server found in the list of enabled servers",
				slog.String("subscriptionManagerID", subscriptionManagerID),
				slog.Bool("HCCUpdateServer", true),
			)
			return true, nil
		}
	}
	logger.Debug("server not found in the list of enabled servers",
		slog.String("subscriptionManagerID", subscriptionManagerID),
	)
	return false, nil
}

// ensureSubscriptionManagerIDAuthorizedToUpdate checks if a server with
// the subscription manager ID is authorized to update the domain.
// Returns a Forbidden error if it is not.
func ensureSubscriptionManagerIDAuthorizedToUpdate(
	ctx context.Context,
	subscriptionManagerID string,
	servers []model.IpaServer,
) error {
	included, err := subscriptionManagerIDIncluded(ctx, subscriptionManagerID, servers)
	if err != nil {
		return err
	}
	if included {
		return nil
	}
	return internal_errors.NewHTTPErrorF(
		http.StatusForbidden,
		"update server is not authorized to update the domain",
	)
}

// ensureUpdateServerEnabledForUpdates checks if the update server with the subscription
// manager ID is included in the list of servers and is enabled for updates.
// Returns a BadRequest error if it is not.
func ensureUpdateServerEnabledForUpdates(
	ctx context.Context,
	subscriptionManagerID string,
	servers []model.IpaServer,
) error {
	included, err := subscriptionManagerIDIncluded(ctx, subscriptionManagerID, servers)
	if err != nil {
		return err
	}
	if !included {
		return internal_errors.NewHTTPErrorF(
			http.StatusBadRequest,
			"update server's 'Subscription Manager ID' not found in the authorized list of rhel-idm servers",
		)
	}
	return nil
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
