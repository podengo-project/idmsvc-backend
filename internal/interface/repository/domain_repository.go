package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
)

type DomainRegToken struct {
	DomainId     uuid.UUID
	DomainToken  string
	DomainType   public.DomainType
	ExpirationNS uint64
}

// DomainRepository interface
type DomainRepository interface {
	List(ctx context.Context, orgID string, offset, limit int) (output []model.Domain, count int64, err error)
	// PartialUpdate(ctx context.Context, orgId string, data *model.Domain) (output model.Domain, err error)
	// Update(ctx context.Context, orgId string, data *model.Domain) (output model.Domain, err error)
	FindByID(ctx context.Context, orgID string, UUID uuid.UUID) (output *model.Domain, err error)
	DeleteById(ctx context.Context, orgID string, UUID uuid.UUID) (err error)
	Register(ctx context.Context, orgID string, data *model.Domain) (err error)
	UpdateAgent(ctx context.Context, orgID string, data *model.Domain) (err error)
	UpdateUser(ctx context.Context, orgID string, data *model.Domain) (err error)
	CreateDomainToken(ctx context.Context, key []byte, validity time.Duration, orgID string, domainType public.DomainType) (token *DomainRegToken, err error)
}
