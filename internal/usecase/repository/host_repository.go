package repository

import (
	"fmt"

	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
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
		return nil, fmt.Errorf("'db' is nil")
	}
	if options == nil {
		return nil, fmt.Errorf("'options' is nil")
	}

	// look through domains and find domains with non-NULL token
	// TODO: match FQDN with domain realms
	var domains []model.Domain
	tx := db.Model(&model.Domain{}).
		Joins("left join ipas on domains.id = ipas.id").
		Where("ipas.token is NULL").
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
