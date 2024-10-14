package api

import (
	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
)

type UpdateDomainAgent interface {
	Build() *public.UpdateDomainAgentRequest
	WithDomainName(value string) UpdateDomainAgent
	WithDomainType(value public.DomainType) UpdateDomainAgent
	WithDomainRhelIdm(value public.DomainIpa) UpdateDomainAgent
	WithSubscriptionManagerID(value string) UpdateDomainAgent
	WithHCCUpdate(value bool) UpdateDomainAgent
}

type updateDomainAgentRequest public.UpdateDomainAgentRequest

func NewUpdateDomainAgent(domainName string) UpdateDomainAgent {
	return &updateDomainAgentRequest{
		DomainName: domainName,
		DomainType: public.RhelIdm,
		RhelIdm:    *NewRhelIdmDomain(domainName).Build(),
	}
}

func (b *updateDomainAgentRequest) Build() *public.UpdateDomainAgentJSONRequestBody {
	return (*public.UpdateDomainAgentRequest)(b)
}

func (b *updateDomainAgentRequest) WithDomainName(value string) UpdateDomainAgent {
	b.DomainName = value
	return b
}

func (b *updateDomainAgentRequest) WithDomainType(value public.DomainType) UpdateDomainAgent {
	b.DomainType = value
	return b
}

func (b *updateDomainAgentRequest) WithDomainRhelIdm(value public.DomainIpa) UpdateDomainAgent {
	b.RhelIdm = value
	return b
}

func (b *updateDomainAgentRequest) WithSubscriptionManagerID(value string) UpdateDomainAgent {
	if len(b.RhelIdm.Servers) > 0 {
		*b.RhelIdm.Servers[0].SubscriptionManagerId = uuid.MustParse(value)
	}
	return b
}

func (b *updateDomainAgentRequest) WithHCCUpdate(value bool) UpdateDomainAgent {
	if len(b.RhelIdm.Servers) > 0 {
		b.RhelIdm.Servers[0].HccUpdateServer = value
	}
	return b
}
