package model

import (
	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

type IpaServer interface {
	Build() model.IpaServer
	WithIpaID(value uint) IpaServer
	WithFQDN(value string) IpaServer
	WithRHSMID(value *string) IpaServer
	WithLocation(value *string) IpaServer
	WithCaServer(value bool) IpaServer
	WithHCCEnrollmentServer(value bool) IpaServer
	WithHCCUpdateServer(value bool) IpaServer
	WithPKInitServer(value bool) IpaServer
	// IpaID               uint
	// FQDN                string
	// RHSMId              *string `gorm:"unique;column:rhsm_id"`
	// Location            *string
	// CaServer            bool
	// HCCEnrollmentServer bool
	// HCCUpdateServer     bool
	// PKInitServer        bool
}

type ipaServer struct {
	IpaServer model.IpaServer
}

func NewIpaServer(gormModel gorm.Model) IpaServer {
	var (
		rhsmID   *string
		location *string
	)
	if builder_helper.GenRandBool() {
		rhsmID = pointy.String(uuid.NewString())
	}
	if builder_helper.GenRandBool() {
		location = pointy.String(builder_helper.GenRandDomainLabel())
	}
	return &ipaServer{
		IpaServer: model.IpaServer{
			Model:               gormModel,
			IpaID:               0,
			FQDN:                builder_helper.GenRandFQDN(),
			RHSMId:              rhsmID,
			Location:            location,
			CaServer:            builder_helper.GenRandBool(),
			HCCEnrollmentServer: builder_helper.GenRandBool(),
			HCCUpdateServer:     builder_helper.GenRandBool(),
			PKInitServer:        builder_helper.GenRandBool(),
		},
	}
}

func (b *ipaServer) Build() model.IpaServer {
	return b.IpaServer
}

func (b *ipaServer) WithIpaID(value uint) IpaServer {
	b.IpaServer.IpaID = value
	return b
}

func (b *ipaServer) WithFQDN(value string) IpaServer {
	b.IpaServer.FQDN = value
	return b
}

func (b *ipaServer) WithRHSMID(value *string) IpaServer {
	b.IpaServer.RHSMId = value
	return b
}

func (b *ipaServer) WithLocation(value *string) IpaServer {
	b.IpaServer.Location = value
	return b
}

func (b *ipaServer) WithCaServer(value bool) IpaServer {
	b.IpaServer.CaServer = value
	return b
}

func (b *ipaServer) WithHCCEnrollmentServer(value bool) IpaServer {
	b.IpaServer.HCCEnrollmentServer = value
	return b
}

func (b *ipaServer) WithHCCUpdateServer(value bool) IpaServer {
	b.IpaServer.HCCUpdateServer = value
	return b
}

func (b *ipaServer) WithPKInitServer(value bool) IpaServer {
	b.IpaServer.PKInitServer = value
	return b
}
