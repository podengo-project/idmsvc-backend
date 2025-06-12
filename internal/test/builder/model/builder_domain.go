package model

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

const (
	// MinOrgID the minimal org id for random values
	minOrgID = 1
	// MaxOrgID the maximum org id for random values
	maxOrgID = 999999
)

// Domain is a builder to fill model.Domain structures
// with random data to make life easier during the tests.
type Domain interface {
	Build() *model.Domain
	WithModel(value gorm.Model) Domain
	WithOrgID(value string) Domain
	WithDomainName(value string) Domain
	WithDomainUUID(value uuid.UUID) Domain
	WithTitle(value *string) Domain
	WithDescription(value *string) Domain
	WithAutoEnrollmentEnabled(value *bool) Domain
	WithIpaDomain(value *model.Ipa) Domain
}

// domain is the specific builder implementation
type domain model.Domain

// NewDomain create a new builder for mode.Domain data.
// Return the Domain builder interface.
func NewDomain(gormModel gorm.Model) Domain {
	var autoEnrollmentEnabled *bool
	var title *string
	var description *string
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	switch builder_helper.GenRandNum(1, 3) {
	case 1:
		autoEnrollmentEnabled = pointy.Bool(false)
	case 2:
		autoEnrollmentEnabled = pointy.Bool(true)
	case 3:
		autoEnrollmentEnabled = nil
	default:
		panic("something went wrong for random AutoEnrollmentEnabled")
	}

	switch builder_helper.GenRandNum(1, 2) {
	case 1:
		title = pointy.String(builder_helper.GenRandString(letters, int(builder_helper.GenRandNum(5, 20))))
	case 2:
		title = nil
	default:
		panic("something went wrong for random title")
	}

	switch builder_helper.GenRandNum(1, 2) {
	case 1:
		description = pointy.String(builder_helper.GenRandString(letters, int(builder_helper.GenRandNum(10, 300))))
	case 2:
		description = nil
	default:
		panic("something went wrong for random title")
	}

	return &domain{
		Model:                 gormModel,
		OrgId:                 strconv.Itoa(int(builder_helper.GenRandNum(minOrgID, maxOrgID))),
		DomainUuid:            uuid.New(),
		DomainName:            pointy.String(builder_helper.GenRandDomainName(2)),
		AutoEnrollmentEnabled: autoEnrollmentEnabled,
		Title:                 title,
		Description:           description,
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain:             NewIpaDomain().WithModel(gormModel).Build(),
	}
}

// Build generate the mode.Domain instance based into
// the current inputs and random data to fill the gaps.
// Return the generated model.Domain instance.
func (b *domain) Build() *model.Domain {
	return (*model.Domain)(b)
}

func (b *domain) WithModel(value gorm.Model) Domain {
	b.Model = value
	return b
}

func (b *domain) WithOrgID(value string) Domain {
	b.OrgId = value
	return b
}

func (b *domain) WithAutoEnrollmentEnabled(value *bool) Domain {
	b.AutoEnrollmentEnabled = value
	return b
}

func (b *domain) WithDomainUUID(value uuid.UUID) Domain {
	b.DomainUuid = value
	return b
}

func (b *domain) WithTitle(value *string) Domain {
	b.Title = value
	return b
}

func (b *domain) WithDescription(value *string) Domain {
	b.Description = value
	return b
}

func (b *domain) WithIpaDomain(value *model.Ipa) Domain {
	b.Type = pointy.Uint(model.DomainTypeIpa)
	b.IpaDomain = value
	return b
}

func (b *domain) WithDomainName(value string) Domain {
	b.DomainName = pointy.String(value)
	return b
}
