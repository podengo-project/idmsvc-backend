package repository

import (
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"gorm.io/gorm"
)

// HostRepository interface
type HostRepository interface {
	MatchDomain(db *gorm.DB, options *interactor.HostConfOptions) (output *model.Domain, err error)
}
