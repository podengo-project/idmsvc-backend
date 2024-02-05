package model

import (
	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"gorm.io/gorm"
)

// HostConfOptions is a builder to fill model.Domain structures
// with random data to make life easier during the tests.
type HostConfOptions interface {
	Build() *interactor.HostConfOptions
	WithModel(value gorm.Model) HostConfOptions
	WithOrgID(value string) HostConfOptions
	WithDomainUUID(value uuid.UUID) HostConfOptions
	WithTitle(value *string) HostConfOptions
	WithDescription(value *string) HostConfOptions
	WithAutoEnrollmentEnabled(value *bool) HostConfOptions
	WithIpaDomain(value *model.Ipa) HostConfOptions
}
