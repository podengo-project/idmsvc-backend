package repository

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/interface/repository"
	"gorm.io/gorm"
)

type hostRepository struct{}

func NewHostRepository() repository.HostRepository {
	return &hostRepository{}
}

// MatchDomain() uses information from `options` to find a matching domain.
// It returns an error when either no matching domain is found or multiple
// domains are matching.
func (r *hostRepository) MatchDomain(db *gorm.DB, options *interactor.HostConfOptions) (output *model.Domain, err error) {
	if db == nil {
		return nil, internal_errors.NilArgError("db")
	}
	if options == nil {
		return nil, internal_errors.NilArgError("options")
	}

	// look through domains and find domains with non-NULL token
	// TODO: match FQDN with domain realms
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
	if err = tx.Limit(2).Find(&domains).Error; err != nil {
		// returns gorm.ErrRecordNotFound when no domain is configured
		return nil, err
	}

	// only one domain is currently supported. Fail if query found multiple doamins.
	if len(domains) != 1 {
		return nil, fmt.Errorf("matched %d domains found, only one expected", len(domains))
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
	tok, err := token.BuildHostconfToken(
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
	b, err := token.SignToken(tok, privs)
	if err != nil {
		return "", err
	}
	return public.HostToken(b), nil
}
