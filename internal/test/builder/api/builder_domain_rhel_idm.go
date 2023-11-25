package api

import (
	"strings"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
)

type RhelIdmDomain interface {
	Build() *public.DomainIpa
	WithRealmDomains(values []string) RhelIdmDomain
	AddRealmDomain(value string) RhelIdmDomain
	WithCaCerts(values []public.Certificate) RhelIdmDomain
	AddCaCert(value public.Certificate) RhelIdmDomain
	WithLocations(values []public.Location) RhelIdmDomain
	AddLocation(value public.Location) RhelIdmDomain
	WithServers(values []public.DomainIpaServer) RhelIdmDomain
	AddServer(value public.DomainIpaServer) RhelIdmDomain
	WithAutomounLocations(values *[]string) RhelIdmDomain
	AddAutomountLocation(value string) RhelIdmDomain
}

type rhelIdmDomain public.DomainIpa

func NewRhelIdmDomain(domainName string) RhelIdmDomain {
	data := &rhelIdmDomain{
		RealmName:          strings.ToUpper(domainName),
		RealmDomains:       []string{domainName},
		CaCerts:            []public.Certificate{},
		Locations:          []public.Location{},
		Servers:            []public.DomainIpaServer{},
		AutomountLocations: nil,
	}
	// TODO Implement to fill some random values
	data.CaCerts = append(data.CaCerts,
		NewCertificate(data.RealmName).
			Build(),
	)
	locationLabel := builder_helper.GenRandLocationLabel()
	data.Locations = append(data.Locations,
		NewLocation().
			WithName(locationLabel).
			WithDescription(locationLabel).
			Build(),
	)
	data.Servers = append(data.Servers,
		NewDomainIpaServer(builder_helper.GenRandFQDNWithDomain(domainName)).
			Build(),
	)
	// TODO Add some AutomountLocation
	return data
}

func (b *rhelIdmDomain) Build() *public.DomainIpa {
	return (*public.DomainIpa)(b)
}

func (b *rhelIdmDomain) WithRealmDomains(values []string) RhelIdmDomain {
	b.RealmDomains = values
	return b
}

func (b *rhelIdmDomain) AddRealmDomain(value string) RhelIdmDomain {
	b.RealmDomains = append(b.RealmDomains, value)
	return b
}

func (b *rhelIdmDomain) WithCaCerts(values []public.Certificate) RhelIdmDomain {
	b.CaCerts = values
	return b
}

func (b *rhelIdmDomain) AddCaCert(value public.Certificate) RhelIdmDomain {
	b.CaCerts = append(b.CaCerts, value)
	return b
}

func (b *rhelIdmDomain) WithLocations(values []public.Location) RhelIdmDomain {
	b.Locations = values
	return b
}

func (b *rhelIdmDomain) AddLocation(value public.Location) RhelIdmDomain {
	b.Locations = append(b.Locations, value)
	return b
}

func (b *rhelIdmDomain) WithServers(values []public.DomainIpaServer) RhelIdmDomain {
	b.Servers = values
	return b
}

func (b *rhelIdmDomain) AddServer(value public.DomainIpaServer) RhelIdmDomain {
	b.Servers = append(b.Servers, value)
	return b
}

func (b *rhelIdmDomain) WithAutomounLocations(values *[]string) RhelIdmDomain {
	b.AutomountLocations = values
	return b
}

func (b *rhelIdmDomain) AddAutomountLocation(value string) RhelIdmDomain {
	if b.AutomountLocations == nil {
		b.AutomountLocations = &[]string{}
	}
	*b.AutomountLocations = append(*b.AutomountLocations, value)
	return b
}
