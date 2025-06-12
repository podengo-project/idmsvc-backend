package api

import (
	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"go.openly.dev/pointy"
)

type Domain interface {
	Build() *public.Domain
	WithDomainID(value *uuid.UUID) Domain
	WithAutoEnrollmentEnabled(value *bool) Domain
	WithDescription(value *string) Domain
	WithTitle(value *string) Domain
	WithRhelIdm(value *public.DomainIpa) Domain
}

type domain public.Domain

func NewDomain(domainName string) Domain {
	domainID := uuid.New()
	return &domain{
		DomainId:              &domainID,
		AutoEnrollmentEnabled: builder_helper.GenRandPointyBool(),
		Description:           pointy.String(builder_helper.GenRandParagraph(0)),
		DomainName:            domainName,
		DomainType:            public.RhelIdm,
		Title:                 pointy.String(domainName),
		RhelIdm:               NewRhelIdmDomain(domainName).Build(),
	}
}

func (b *domain) Build() *public.Domain {
	return (*public.Domain)(b)
}

func (b *domain) WithDomainID(value *uuid.UUID) Domain {
	b.DomainId = value
	return b
}

func (b *domain) WithAutoEnrollmentEnabled(value *bool) Domain {
	b.AutoEnrollmentEnabled = value
	return b
}

func (b *domain) WithDescription(value *string) Domain {
	b.Description = value
	return b
}

func (b *domain) WithTitle(value *string) Domain {
	b.Title = value
	return b
}

func (b *domain) WithRhelIdm(value *public.DomainIpa) Domain {
	b.RhelIdm = value
	return b
}
