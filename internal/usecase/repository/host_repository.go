package repository

import (
	"fmt"

	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/hmsidm/internal/interface/repository"
	"gorm.io/gorm"
)

type hostRepository struct{}

func NewHostRepository() repository.HostRepository {
	return &hostRepository{}
}

func (r *hostRepository) MatchDomain(db *gorm.DB, options *interactor.HostConfOptions) (output *model.Domain, err error) {
	if db == nil {
		return nil, fmt.Errorf("'db' is nil")
	}

	return nil, fmt.Errorf("TODO: not implemented")
}
