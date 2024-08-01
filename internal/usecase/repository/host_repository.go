package repository

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_token"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
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
// ctx is the current request context with db and slog instances.
// options provide filtering information to select the domain.
// Return the matched domain and nil on success, else nil and the error
// instance with additional information.
func (r *hostRepository) MatchDomain(ctx context.Context, options *interactor.HostConfOptions) (output *model.Domain, err error) {
	db := app_context.DBFromCtx(ctx)
	log := app_context.LogFromCtx(ctx)
	if db == nil {
		err = internal_errors.NilArgError("db")
		log.Error(err.Error())
		return nil, err
	}
	if options == nil {
		err = internal_errors.NilArgError("options")
		log.Error(err.Error())
		return nil, err
	}

	logCriteria := make([]string, 0, 3)
	var domains []model.Domain
	tx := db.Model(&model.Domain{}).
		Joins("left join ipas on domains.id = ipas.id").
		Where("domains.org_id = ?", options.OrgId)
	if options.DomainId != nil {
		domainUUID := options.DomainId.String()
		tx = tx.Where("domains.domain_uuid = ?", domainUUID)
		logCriteria = append(logCriteria, fmt.Sprintf("domains.domain_uuid = %s", domainUUID))
	}
	if options.DomainName != nil {
		tx = tx.Where("domains.domain_name = ?", *options.DomainName)
		logCriteria = append(logCriteria, fmt.Sprintf("domains.domain_name = %s", *options.DomainName))
	}
	if options.DomainType != nil {
		tx = tx.Where("domains.type = ?", model.DomainTypeUint(string(*options.DomainType)))
		logCriteria = append(logCriteria, fmt.Sprintf("domains.type = %s", string(*options.DomainType)))
	}
	if err = tx.Find(&domains).Error; err != nil {
		// empty result set is not an error here, but is handled below
		log.Error(fmt.Sprintf("finding a 'rhel-idm' domain which match the criteria: %s", strings.Join(logCriteria, ", ")))
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
		err = internal_errors.NewHTTPErrorF(
			http.StatusNotFound,
			"no matching domains",
		)
		log.Error("no matching domains")
		return nil, err
	} else if len(matchedDomains) > 1 {
		err = internal_errors.NewHTTPErrorF(
			http.StatusConflict,
			"matched %d domains, only one expected",
			len(matchedDomains),
		)
		log.Error("more than one domain found but only one expected")
		return nil, err
	}

	// verify and fill domain object
	output = &domains[0]
	if err = output.FillAndPreload(db); err != nil {
		log.Error(fmt.Sprintf("preloading domain data for output.domain_id = %s", output.DomainUuid.String()))
		return nil, err
	}
	return output, nil
}

// SignHostConfToken
// ctx is the current request context with db and slog instances.
// privs
// options
// Return the matched domain and nil on success, else nil and the error
// instance with additional information.
func (r *hostRepository) SignHostConfToken(
	ctx context.Context,
	privs []jwk.Key,
	options *interactor.HostConfOptions,
	domain *model.Domain,
) (hctoken public.HostToken, err error) {
	log := app_context.LogFromCtx(ctx)
	if options == nil {
		err = internal_errors.NilArgError("options")
		log.Error(err.Error())
		return "", err
	}
	if domain == nil {
		err = internal_errors.NilArgError("domain")
		log.Error(err.Error())
		return "", err
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
		log.Error("error building hostconf token")
		return "", err
	}
	b, err := hostconf_token.SignToken(tok, privs)
	if err != nil {
		log.Error("error signing hostconf token")
		return "", err
	}
	return public.HostToken(b), nil
}
