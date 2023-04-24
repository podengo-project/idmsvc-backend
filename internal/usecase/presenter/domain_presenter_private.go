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
func (p domainPresenter) registerIpaCaCerts(
	domain *model.Domain,
	output *public.Domain,
) error {
	ipa := domain.IpaDomain
	if ipa.CaCerts == nil {
		return fmt.Errorf("'ipa.CaCerts' is nil")
	}
	// output.RhelIdm.CaCerts = make([]public.DomainIpaCert, len(ipa.CaCerts))
	// for i := range ipa.CaCerts {
	// 	err := p.FillCert(&output.RhelIdm.CaCerts[i], &ipa.CaCerts[i])
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	if err := p.fillRhelIdmCerts(output, domain); err != nil {
		return err
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
func (p domainPresenter) registerIpaServers(
	domain *model.Domain,
	output *public.Domain,
) error {
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
	if err := p.fillRhelIdmServers(output, domain); err != nil {
		return err
	}

	return nil
}

func (p domainPresenter) registerFillDomainData(
	domain *model.Domain,
	output *public.Domain,
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
	output *public.Domain,
) (err error) {
	if domain.IpaDomain == nil {
		return fmt.Errorf("'domain.IpaDomain' is nil")
	}
	if output == nil {
		panic("'output' is nil")
	}
	if domain.IpaDomain.RealmName != nil {
		output.RhelIdm.RealmName = *domain.IpaDomain.RealmName
	}

	if domain.IpaDomain.RealmDomains != nil {
		output.RhelIdm.RealmDomains = append([]string{}, domain.IpaDomain.RealmDomains...)
	}

	if err = p.registerIpaCaCerts(domain, output); err != nil {
		return err
	}

	if err = p.registerIpaServers(domain, output); err != nil {
		return err
	}

	return nil
}

func (p domainPresenter) getCheckRhelIdm(domain *model.Domain) error {
	if domain.IpaDomain == nil {
		return fmt.Errorf("'IpaDomain' is nil")
	}
	if domain.IpaDomain.CaCerts == nil {
		return fmt.Errorf("'CaCerts' is nil")
	}
	if domain.IpaDomain.RealmName == nil {
		return fmt.Errorf("'RealmName' is nil")
	}
	if domain.IpaDomain.Servers == nil {
		return fmt.Errorf("'Servers' is nil")
	}
	return nil
}

func (p domainPresenter) getChecks(domain *model.Domain) error {
	if domain == nil {
		return fmt.Errorf("'domain' is nil")
	}
	if domain.AutoEnrollmentEnabled == nil {
		return fmt.Errorf("'AutoenrollmentEnabled' is nil")
	}
	if domain.Type == nil {
		return fmt.Errorf("'DomainType' is nil")
	}
	switch *domain.Type {
	case model.DomainTypeIpa:
		return p.getCheckRhelIdm(domain)
	default:
		return fmt.Errorf("'Type' is invalid")
	}
}

func (p domainPresenter) fillRhelIdmServers(output *public.Domain, domain *model.Domain) error {
	if domain.IpaDomain == nil {
		return nil
	}
	output.RhelIdm.Servers = make([]public.DomainIpaServer, len(domain.IpaDomain.Servers))
	for i := range domain.IpaDomain.Servers {
		output.RhelIdm.Servers[i].Fqdn = domain.IpaDomain.Servers[i].FQDN
		output.RhelIdm.Servers[i].SubscriptionManagerId = domain.IpaDomain.Servers[i].RHSMId
		output.RhelIdm.Servers[i].CaServer = domain.IpaDomain.Servers[i].CaServer
		output.RhelIdm.Servers[i].HccEnrollmentServer = domain.IpaDomain.Servers[i].HCCEnrollmentServer
		output.RhelIdm.Servers[i].HccUpdateServer = domain.IpaDomain.Servers[i].HCCUpdateServer
		output.RhelIdm.Servers[i].PkinitServer = domain.IpaDomain.Servers[i].PKInitServer
	}
	return nil
}

func (p domainPresenter) fillRhelIdmCerts(output *public.Domain, domain *model.Domain) error {
	if domain.IpaDomain == nil {
		return nil
	}
	output.RhelIdm.CaCerts = make([]public.DomainIpaCert, len(domain.IpaDomain.CaCerts))
	for i := range domain.IpaDomain.CaCerts {
		output.RhelIdm.CaCerts[i].Nickname = domain.IpaDomain.CaCerts[i].Nickname
		output.RhelIdm.CaCerts[i].Issuer = domain.IpaDomain.CaCerts[i].Issuer
		output.RhelIdm.CaCerts[i].NotValidAfter = domain.IpaDomain.CaCerts[i].NotValidAfter
		output.RhelIdm.CaCerts[i].NotValidBefore = domain.IpaDomain.CaCerts[i].NotValidBefore
		output.RhelIdm.CaCerts[i].SerialNumber = domain.IpaDomain.CaCerts[i].SerialNumber
		output.RhelIdm.CaCerts[i].Subject = domain.IpaDomain.CaCerts[i].Subject
		output.RhelIdm.CaCerts[i].Pem = domain.IpaDomain.CaCerts[i].Pem
	}
	return nil
}

func (p domainPresenter) createRhelIdmCheckDomain(domain *model.Domain) error {
	if domain.IpaDomain.RealmName == nil {
		return fmt.Errorf("'RealmName' is nil")
	}
	if domain.IpaDomain.CaCerts == nil {
		return fmt.Errorf("'CaCerts' is nil")
	}
	if domain.IpaDomain.Servers == nil {
		return fmt.Errorf("'Servers' is nil")
	}
	return nil
}

func (p domainPresenter) createRhelIdmFillRealmDomains(output *public.Domain, domain *model.Domain) {
	if domain.IpaDomain.RealmDomains == nil {
		output.RhelIdm.RealmDomains = []string{}
	} else {
		output.RhelIdm.RealmDomains = domain.IpaDomain.RealmDomains
	}
}

func (p domainPresenter) createRhelIdm(output *public.Domain, domain *model.Domain) error {
	if err := p.createRhelIdmCheckDomain(domain); err != nil {
		return err
	}
	output.RhelIdm = &public.DomainIpa{}
	output.RhelIdm.RealmName = *domain.IpaDomain.RealmName
	if err := p.fillRhelIdmCerts(output, domain); err != nil {
		return err
	}
	if err := p.fillRhelIdmServers(output, domain); err != nil {
		return err
	}
	p.createRhelIdmFillRealmDomains(output, domain)
	return nil
}
