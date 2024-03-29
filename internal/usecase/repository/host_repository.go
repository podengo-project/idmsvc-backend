package repository

import (
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_token"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"gorm.io/gorm"
)

type hostRepository struct{}

func NewHostRepository() repository.HostRepository {
	return &hostRepository{}
}

// MatchDomain uses information from `options` to find a matching domain.
// Return an error when either no matching domain is found or multiple
// domains are matching.
//
// Exclude domains with auto_enrollment_enabled = FALSE.
func (r *hostRepository) MatchDomain(db *gorm.DB, options *interactor.HostConfOptions) (output *model.Domain, err error) {
	if db == nil {
		return nil, internal_errors.NilArgError("db")
	}
	if options == nil {
		return nil, internal_errors.NilArgError("options")
	}

	var domains []model.Domain
	tx := db.Model(&model.Domain{}).
		Joins("left join ipas on domains.id = ipas.id").
		Where("domains.org_id = ?", options.OrgId)
	if options.DomainId != nil {
		tx = tx.Where("domains.domain_uuid = ?", options.DomainId.String())
	}
	if options.DomainName != nil {
		tx = tx.Where("domains.domain_name = ?", *options.DomainName)
	}
	if options.DomainType != nil {
		tx = tx.Where("domains.type = ?", model.DomainTypeUint(string(*options.DomainType)))
	}
	if err = tx.Find(&domains).Error; err != nil {
		// empty result set is not an error here, but is handled below
		return nil, err
	}

	matchedDomains := make([]model.Domain, 0, len(domains))
	for _, domain := range domains {
		if domain.AutoEnrollmentEnabled == nil || !(*domain.AutoEnrollmentEnabled) {
			continue
		}
		// TODO: match FQDN with domain realms
		matchedDomains = append(matchedDomains, domain)
	}

	// only one domain is currently supported. Fail if query found multiple doamins.
	if len(matchedDomains) < 1 {
		return nil, internal_errors.NewHTTPErrorF(
			http.StatusNotFound,
			"no matching domains",
		)
	} else if len(matchedDomains) > 1 {
		return nil, internal_errors.NewHTTPErrorF(
			http.StatusConflict,
			"matched %d domains, only one expected",
			len(matchedDomains),
		)
	}

	// verify and fill domain object
	output = &domains[0]
	if err = output.FillAndPreload(db); err != nil {
		return nil, err
	}
	return output, nil
}

func (r *hostRepository) SignHostConfToken(
	privs []jwk.Key, options *interactor.HostConfOptions, domain *model.Domain,
) (hctoken public.HostToken, err error) {
	if options == nil {
		return "", internal_errors.NilArgError("options")
	}
	if domain == nil {
		return "", internal_errors.NilArgError("domain")
	}

	validity := time.Hour
	tok, err := hostconf_token.BuildHostconfToken(
		options.CommonName,
		options.OrgId,
		options.InventoryId,
		options.Fqdn,
		domain.DomainUuid,
		validity,
	)
	if err != nil {
		return "", err
	}
	b, err := hostconf_token.SignToken(tok, privs)
	if err != nil {
		return "", err
	}
	return public.HostToken(b), nil
}
