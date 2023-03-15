package impl

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hmsidm/internal/domain/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (a *application) findIpaById(tx *gorm.DB, orgId string, uuid string) (data *model.Domain, err error) {
	data = &model.Domain{}
	*data, err = a.domain.repository.FindById(tx, orgId, uuid)
	if err != nil {
		return nil, err
	}
	if *data.DomainType != model.DomainTypeIpa {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Wrong domain type")
	}

	if data.IpaDomain == nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "No IPA data found for the domain")
	}
	return data, nil
}

func (a *application) checkToken(token string, ipa *model.Ipa) error {
	if token == "" {
		return fmt.Errorf("'token' cannot be empty")
	}
	if ipa == nil {
		return fmt.Errorf("'ipa' cannot be nil")
	}
	if ipa.Token == nil || *ipa.Token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "OTP token is required")
	}
	if ipa.TokenExpiration == nil || (*ipa.TokenExpiration == time.Time{}) {
		return echo.NewHTTPError(http.StatusBadRequest, "OTP token expiration not found")
	}
	if *ipa.Token != token {
		return echo.NewHTTPError(http.StatusBadRequest, "OTP token does not match")
	}
	if ipa.TokenExpiration.After(time.Now()) {
		return echo.NewHTTPError(http.StatusBadRequest, "OTP token expired")
	}
	return nil
}

func (a *application) checkSubscriptionIdWithInventory(xrhiden string, subscription_manager_id string) (string, error) {
	hostList, err := a.inventory.ListHost(xrhiden)
	if err != nil {
		return "", err
	}
	for _, host := range hostList {
		if host.SubscriptionManagerId == subscription_manager_id {
			return host.FQDN, nil
		}
	}
	return "", fmt.Errorf("'subscription_manager_id' not found")
}

func (a *application) checkEnrollmentIpaServer(fqdn string, ipa *model.Ipa) error {
	for _, item := range ipa.Servers {
		if item.FQDN == fqdn && item.HCCEnrollmentServer {
			return nil
		}
	}
	return fmt.Errorf("IPA server not found or not enabled to register")
}
