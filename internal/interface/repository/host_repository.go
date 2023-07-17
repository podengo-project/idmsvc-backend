package repository

import (
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"gorm.io/gorm"
)

// HostRepository interface
type HostRepository interface {
	MatchDomain(db *gorm.DB, options *interactor.HostConfOptions) (output *model.Domain, err error)
}
