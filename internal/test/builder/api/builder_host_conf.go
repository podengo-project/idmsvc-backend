package api

import (
	"github.com/google/uuid"
	"github.com/openlyinc/pointy"
	"github.com/pioz/faker"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
)

// HostConf builder interface
type HostConf interface {
	Build() *public.HostConf
	WithDomainId(value *public.DomainId) HostConf
	WithDomainName(value *public.DomainName) HostConf
	WithDomainType(value *public.DomainType) HostConf
}

type hostConf public.HostConf

func NewHostConf() HostConf {
	id := &uuid.UUID{}
	*id = uuid.New()
	return &hostConf{
		DomainId:   id,
		DomainName: pointy.String(helper.GenRandDomainName(2)),
		DomainType: (*public.DomainType)(pointy.String(faker.Pick(string(public.RhelIdm)))),
	}
}

func (b *hostConf) Build() *public.HostConf {
	return (*public.HostConf)(b)
}

func (b *hostConf) WithDomainId(value *public.DomainId) HostConf {
	b.DomainId = value
	return b
}

func (b *hostConf) WithDomainName(value *public.DomainName) HostConf {
	b.DomainName = value
	return b
}

func (b *hostConf) WithDomainType(value *public.DomainType) HostConf {
	b.DomainType = value
	return b
}
