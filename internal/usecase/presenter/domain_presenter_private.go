package presenter

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"go.openly.dev/pointy"
)

func (p *domainPresenter) fillRhelIdmLocations(
	target *public.Domain,
	source *model.Domain,
) {
	if target == nil || target.RhelIdm == nil {
		panic("'target' or 'target.RhelIdm' are nil")
	}
	if source == nil || source.IpaDomain == nil {
		panic("'source' or 'source.IpaDomain' are nil")
	}
	target.RhelIdm.Locations = make(
		[]public.Location,
		len(source.IpaDomain.Locations),
	)
	for idx := range source.IpaDomain.Locations {
		target.RhelIdm.Locations[idx].Name = source.IpaDomain.Locations[idx].Name
		target.RhelIdm.Locations[idx].Description = source.IpaDomain.Locations[idx].Description
	}
}

func (p *domainPresenter) fillRhelIdmServers(
	target *public.Domain,
	source *model.Domain,
) {
	if target == nil || source == nil {
		return
	}
	if target.RhelIdm == nil || source.IpaDomain == nil {
		return
	}
	target.RhelIdm.Servers = make(
		[]public.DomainIpaServer,
		len(source.IpaDomain.Servers),
	)
	for i := range source.IpaDomain.Servers {
		target.RhelIdm.Servers[i].Fqdn =
			source.IpaDomain.Servers[i].FQDN
		var rhsmID *uuid.UUID = nil
		if source.IpaDomain.Servers[i].RHSMId != nil {
			rhsmID = &uuid.UUID{}
			*rhsmID = uuid.MustParse(*source.IpaDomain.Servers[i].RHSMId)
		}
		target.RhelIdm.Servers[i].SubscriptionManagerId = rhsmID
		target.RhelIdm.Servers[i].Location =
			source.IpaDomain.Servers[i].Location
		target.RhelIdm.Servers[i].CaServer =
			source.IpaDomain.Servers[i].CaServer
		target.RhelIdm.Servers[i].HccEnrollmentServer =
			source.IpaDomain.Servers[i].HCCEnrollmentServer
		target.RhelIdm.Servers[i].HccUpdateServer =
			source.IpaDomain.Servers[i].HCCUpdateServer
		target.RhelIdm.Servers[i].PkinitServer =
			source.IpaDomain.Servers[i].PKInitServer
	}
}

func (p *domainPresenter) fillRhelIdmCerts(
	output *public.Domain,
	domain *model.Domain,
) {
	if output == nil || domain == nil || output.RhelIdm == nil || domain.IpaDomain == nil {
		return
	}
	output.RhelIdm.CaCerts = make(
		[]public.Certificate,
		len(domain.IpaDomain.CaCerts),
	)
	for i := range domain.IpaDomain.CaCerts {
		output.RhelIdm.CaCerts[i].Nickname =
			domain.IpaDomain.CaCerts[i].Nickname
		output.RhelIdm.CaCerts[i].Issuer =
			domain.IpaDomain.CaCerts[i].Issuer
		output.RhelIdm.CaCerts[i].NotAfter =
			domain.IpaDomain.CaCerts[i].NotAfter
		output.RhelIdm.CaCerts[i].NotBefore =
			domain.IpaDomain.CaCerts[i].NotBefore
		output.RhelIdm.CaCerts[i].SerialNumber =
			domain.IpaDomain.CaCerts[i].SerialNumber
		output.RhelIdm.CaCerts[i].Subject =
			domain.IpaDomain.CaCerts[i].Subject
		output.RhelIdm.CaCerts[i].Pem =
			domain.IpaDomain.CaCerts[i].Pem
	}
}

func (p *domainPresenter) guardSharedDomain(
	domain *model.Domain,
) error {
	if domain == nil {
		return internal_errors.NilArgError("domain")
	}
	if domain.Type == nil {
		return internal_errors.NilArgError("domain.Type")
	}
	if *domain.Type == model.DomainTypeUndefined {
		return fmt.Errorf("'domain.Type' is invalid")
	}
	return nil
}

func (p *domainPresenter) sharedDomain(
	domain *model.Domain,
) (output *public.Domain, err error) {
	// Expect domain not nil and domain.Type filled
	if err = p.guardSharedDomain(domain); err != nil {
		return nil, err
	}

	// Domain common code
	output = &public.Domain{}
	p.sharedDomainFill(domain, output)

	switch *domain.Type {
	case model.DomainTypeIpa:
		// Specific rhel-idm domain code
		output.DomainType = model.DomainTypeIpaString
		output.RhelIdm = &public.DomainIpa{}
		err = p.sharedDomainFillRhelIdm(domain, output)
	default:
		err = fmt.Errorf("'domain.DomainType=%d' is invalid", *domain.Type)
	}
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (p *domainPresenter) sharedDomainFill(
	domain *model.Domain,
	output *public.Domain,
) {
	if domain.DomainUuid == (uuid.UUID{}) {
		output.DomainId = nil
	} else {
		output.DomainId = &domain.DomainUuid
	}
	if domain.AutoEnrollmentEnabled != nil {
		output.AutoEnrollmentEnabled = domain.AutoEnrollmentEnabled
	}
	if domain.DomainName != nil {
		output.DomainName = *domain.DomainName
	}
	if domain.Title != nil {
		output.Title = domain.Title
	}
	if domain.Description != nil {
		output.Description = domain.Description
	}
}

func (p *domainPresenter) sharedDomainFillRhelIdm(
	domain *model.Domain,
	output *public.Domain,
) (err error) {
	if domain.Type != nil && *domain.Type != model.DomainTypeIpa {
		return fmt.Errorf(
			"'domain.Type' is not '%s'",
			model.DomainTypeIpaString,
		)
	}
	if domain.IpaDomain == nil {
		return internal_errors.NilArgError("domain.IpaDomain")
	}
	if output == nil {
		panic("'output' is nil")
	}
	output.RhelIdm = &public.DomainIpa{}
	if domain.IpaDomain.RealmName != nil {
		output.RhelIdm.RealmName = *domain.IpaDomain.RealmName
	}

	if domain.IpaDomain.RealmDomains != nil {
		output.RhelIdm.RealmDomains = append(
			[]string{},
			domain.IpaDomain.RealmDomains...)
	} else {
		output.RhelIdm.RealmDomains = []string{}
	}

	p.fillRhelIdmCerts(output, domain)
	p.fillRhelIdmServers(output, domain)
	p.fillRhelIdmLocations(output, domain)

	return nil
}

func (p *domainPresenter) buildPaginationLink(offset int, limit int) string {
	if limit == 0 {
		limit = p.cfg.Application.PaginationDefaultLimit
	}
	if limit > p.cfg.Application.PaginationMaxLimit {
		limit = p.cfg.Application.PaginationMaxLimit
	}
	if offset < 0 {
		offset = 0
	}

	q := url.Values{}
	q.Add("limit", strconv.FormatInt(int64(limit), 10))
	q.Add("offset", strconv.FormatInt(int64(offset), 10))

	return fmt.Sprintf("%s/domains?%s", p.cfg.Application.PathPrefix, q.Encode())
}

func (p *domainPresenter) listFillLinks(output *public.ListDomainsResponse, count int64, offset int, limit int) {
	if output == nil {
		panic("'output' is nil")
	}
	if limit == 0 {
		panic("'limit' is zero")
	}

	// Calculate the offsets
	currentOffset := ((offset + limit - 1) / limit) * limit
	firstOffset := 0
	prevOffset := currentOffset - limit
	nextOffset := currentOffset + limit
	lastOffset := (int(count-1) / limit) * limit
	if currentOffset == firstOffset {
		prevOffset = firstOffset
	}
	if currentOffset == lastOffset {
		nextOffset = lastOffset
	}

	// Build the link
	output.Links.First = pointy.String(p.buildPaginationLink(firstOffset, limit))
	if firstOffset != currentOffset {
		output.Links.Previous = pointy.String(p.buildPaginationLink(prevOffset, limit))
	}
	if lastOffset != currentOffset {
		output.Links.Next = pointy.String(p.buildPaginationLink(nextOffset, limit))
	}
	output.Links.Last = pointy.String(p.buildPaginationLink(lastOffset, limit))
}

func (p *domainPresenter) listFillMeta(output *public.ListDomainsResponse, count int64, offset int, limit int) {
	if output == nil {
		panic("'output' is nil")
	}
	output.Meta.Count = count
	output.Meta.Offset = offset
	output.Meta.Limit = limit
}

func (p *domainPresenter) listFillItem(output *public.ListDomainsData, domain *model.Domain) {
	if output == nil {
		panic("'output' is nil")
	}
	if domain == nil {
		panic("'domain' is nil")
	}
	if domain.AutoEnrollmentEnabled == nil {
		output.AutoEnrollmentEnabled = false
	} else {
		output.AutoEnrollmentEnabled = *domain.AutoEnrollmentEnabled
	}
	if domain.DomainName != nil {
		output.DomainName = *domain.DomainName
	}
	output.DomainType = public.DomainType(model.DomainTypeString(*domain.Type))
	output.DomainId = domain.DomainUuid
	if domain.Title != nil {
		output.Title = *domain.Title
	}
	if domain.Description != nil {
		output.Description = *domain.Description
	}
}
