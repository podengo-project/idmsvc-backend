package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"gorm.io/gorm"
)

type DomainRegToken struct {
	DomainId     uuid.UUID
	DomainToken  string
	DomainType   public.DomainType
	ExpirationNS uint64
}

// DomainRepository interface
type DomainRepository interface {
	List(db *gorm.DB, orgID string, offset int, limit int) (output []model.Domain, count int64, err error)
	// PartialUpdate(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error)
	// Update(db *gorm.DB, orgId string, data *model.Domain) (output model.Domain, err error)
	FindByID(db *gorm.DB, orgID string, UUID uuid.UUID) (output *model.Domain, err error)
	DeleteById(db *gorm.DB, orgID string, UUID uuid.UUID) (err error)
	Update(db *gorm.DB, orgID string, data *model.Domain) (err error)
	CreateDomainToken(key []byte, validity time.Duration, orgID string, domainType public.DomainType) (token *DomainRegToken, err error)
}
