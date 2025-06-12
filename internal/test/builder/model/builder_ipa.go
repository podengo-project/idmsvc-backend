package model

import (
	"strings"

	"github.com/lib/pq"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

type IpaDomain interface {
	Build() *model.Ipa
	WithModel(value gorm.Model) IpaDomain
	WithCaCerts(value []model.IpaCert) IpaDomain
	WithServers(value []model.IpaServer) IpaDomain
	WithLocations(value []model.IpaLocation) IpaDomain
	WithRealmName(value *string) IpaDomain
	WithRealmDomains(value pq.StringArray) IpaDomain
}

type ipaDomain struct {
	IpaDomain *model.Ipa
}

func NewIpaDomain() IpaDomain {
	var (
		certs        []model.IpaCert
		locations    []model.IpaLocation
		servers      []model.IpaServer
		realmName    *string
		realmDomains pq.StringArray = pq.StringArray{}
	)
	realmName = pointy.String(strings.ToUpper(builder_helper.GenRandDomainName(2)))
	if realmName != nil {
		realmDomains = pq.StringArray{*realmName}
	}
	return &ipaDomain{
		IpaDomain: &model.Ipa{
			Model:        NewModel().Build(),
			CaCerts:      certs,
			Servers:      servers,
			Locations:    locations,
			RealmName:    realmName,
			RealmDomains: realmDomains,
		},
	}
}

func (b *ipaDomain) Build() *model.Ipa {
	return b.IpaDomain
}

func (b *ipaDomain) WithModel(value gorm.Model) IpaDomain {
	b.IpaDomain.Model = value
	return b
}
func (b *ipaDomain) WithCaCerts(value []model.IpaCert) IpaDomain {
	b.IpaDomain.CaCerts = value
	return b
}

func (b *ipaDomain) WithServers(value []model.IpaServer) IpaDomain {
	b.IpaDomain.Servers = value
	return b
}

func (b *ipaDomain) WithLocations(value []model.IpaLocation) IpaDomain {
	b.IpaDomain.Locations = value
	return b
}

func (b *ipaDomain) WithRealmName(value *string) IpaDomain {
	b.IpaDomain.RealmName = value
	return b
}

func (b *ipaDomain) WithRealmDomains(value pq.StringArray) IpaDomain {
	b.IpaDomain.RealmDomains = value
	return b
}
