package api

import (
	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
)

type DomainIpaServer interface {
	Build() public.DomainIpaServer
	// TODO Add methods and implement them as they are needed
	// With...() DomainIpaServer
	WithHccUpdateServer(value bool) DomainIpaServer
	WithSubscriptionManagerId(value string) DomainIpaServer
}

type domainIpaServer public.DomainIpaServer

func NewDomainIpaServer(fqdn string) DomainIpaServer {
	subscriptionManagerID := &uuid.UUID{}
	*subscriptionManagerID = uuid.New()
	return &domainIpaServer{
		Fqdn:                  fqdn,
		CaServer:              builder_helper.GenRandBool(),
		HccEnrollmentServer:   builder_helper.GenRandBool(),
		HccUpdateServer:       builder_helper.GenRandBool(),
		PkinitServer:          builder_helper.GenRandBool(),
		Location:              nil,
		SubscriptionManagerId: subscriptionManagerID,
	}
}

func (b *domainIpaServer) WithHccUpdateServer(value bool) DomainIpaServer {
	b.HccUpdateServer = value
	return b
}

func (b *domainIpaServer) WithSubscriptionManagerId(value string) DomainIpaServer {
	subscriptionManagerID := &uuid.UUID{}
	*subscriptionManagerID = uuid.MustParse(value)
	b.SubscriptionManagerId = subscriptionManagerID
	return b
}

func (b *domainIpaServer) Build() public.DomainIpaServer {
	return public.DomainIpaServer(*b)
}
