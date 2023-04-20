package presenter

// TODO Too much code duplication
// TODO Investigate if some "inheritence" mechanism in
//      opanapi specification generate common structures
//      letting to reduce the boilerplate when transforming
//      internal types <--> api types

import (
	"fmt"

	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
)

// registerIpaCaCerts translate the list of certificates
// for the Register operation.
// ipa the Ipa instance from the business model.
// output the DomainResponseIpa schema representation
// to be filled.
// Return nil if the information is translated properly
// else return an error.
func (p domainPresenter) registerIpaCaCerts(domain *model.Domain, output *public.RegisterDomainResponse) error {
	ipa := domain.IpaDomain
	if ipa.CaCerts == nil {
		return fmt.Errorf("'ipa.CaCerts' is nil")
	}
	output.Ipa.CaCerts = make([]public.DomainIpaCert, len(ipa.CaCerts))
	for i := range ipa.CaCerts {
		err := p.FillCert(&output.Ipa.CaCerts[i], &ipa.CaCerts[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// registerIpaServers translate the list of servers
// for the Register operation.
// ipa the Ipa instance from the business model.
// output the DomainResponseIpa schema representation
// to be filled.
// Return nil if the information is translated properly
// else return an error.
func (p domainPresenter) registerIpaServers(domain *model.Domain, output *public.RegisterDomainResponse) error {
	ipa := domain.IpaDomain
	if ipa == nil {
		return fmt.Errorf("'IpaDomain' is nil")
	}
	if output == nil {
		return fmt.Errorf("'output' is nil")
	}
	if ipa.Servers == nil {
		return fmt.Errorf("'ipa.Servers' is nil")
	}
	output.Ipa.Servers = make([]public.DomainIpaServer, len(ipa.Servers))
	for i := range ipa.Servers {
		err := p.FillServer(&output.Ipa.Servers[i], &ipa.Servers[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (p domainPresenter) registerFillDomainData(
	domain *model.Domain,
	output *public.RegisterDomainResponse,
) {
	output.DomainUuid = domain.DomainUuid.String()
	if domain.AutoEnrollmentEnabled != nil {
		output.AutoEnrollmentEnabled = *domain.AutoEnrollmentEnabled
	}
	if domain.DomainName != nil {
		output.DomainName = *domain.DomainName
	}
	if domain.Title != nil {
		output.Title = *domain.Title
	}
	if domain.Description != nil {
		output.Description = *domain.Description
	}
}

func (p domainPresenter) registerIpa(
	domain *model.Domain,
	output *public.RegisterDomainResponse,
) (err error) {
	if domain.IpaDomain == nil {
		return fmt.Errorf("'domain.IpaDomain' is nil")
	}
	if output == nil {
		panic("'output' is nil")
	}
	if domain.IpaDomain.RealmName != nil {
		output.Ipa.RealmName = *domain.IpaDomain.RealmName
	}

	if domain.IpaDomain.RealmDomains != nil {
		output.Ipa.RealmDomains = append([]string{}, domain.IpaDomain.RealmDomains...)
	}

	if err = p.registerIpaCaCerts(domain, output); err != nil {
		return err
	}

	if err = p.registerIpaServers(domain, output); err != nil {
		return err
	}

	return nil
}
