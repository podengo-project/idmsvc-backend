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

// sharedDomainFillRhelIdmCaCerts translate the list of certificates
// for the Register operation.
// ipa the Ipa instance from the business model.
// output the DomainResponseIpa schema representation
// to be filled.
// Return nil if the information is translated properly
// else return an error.

func (p domainPresenter) getCheckRhelIdm(
	domain *model.Domain,
) error {
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

func (p domainPresenter) getChecks(
	domain *model.Domain,
) error {
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

func (p domainPresenter) fillRhelIdmServers(
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
		target.RhelIdm.Servers[i].SubscriptionManagerId =
			source.IpaDomain.Servers[i].RHSMId
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

func (p domainPresenter) fillRhelIdmCerts(
	output *public.Domain,
	domain *model.Domain,
) {
	if output == nil || domain == nil || output.RhelIdm == nil || domain.IpaDomain == nil {
		return
	}
	output.RhelIdm.CaCerts = make(
		[]public.DomainIpaCert,
		len(domain.IpaDomain.CaCerts),
	)
	for i := range domain.IpaDomain.CaCerts {
		output.RhelIdm.CaCerts[i].Nickname =
			domain.IpaDomain.CaCerts[i].Nickname
		output.RhelIdm.CaCerts[i].Issuer =
			domain.IpaDomain.CaCerts[i].Issuer
		output.RhelIdm.CaCerts[i].NotValidAfter =
			domain.IpaDomain.CaCerts[i].NotValidAfter
		output.RhelIdm.CaCerts[i].NotValidBefore =
			domain.IpaDomain.CaCerts[i].NotValidBefore
		output.RhelIdm.CaCerts[i].SerialNumber =
			domain.IpaDomain.CaCerts[i].SerialNumber
		output.RhelIdm.CaCerts[i].Subject =
			domain.IpaDomain.CaCerts[i].Subject
		output.RhelIdm.CaCerts[i].Pem =
			domain.IpaDomain.CaCerts[i].Pem
	}
}

func (p domainPresenter) createRhelIdmCheckDomain(
	domain *model.Domain,
) error {
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

func (p domainPresenter) createRhelIdmFillRealmDomains(
	output *public.Domain,
	domain *model.Domain,
) {
	if domain.IpaDomain.RealmDomains == nil {
		output.RhelIdm.RealmDomains = []string{}
	} else {
		output.RhelIdm.RealmDomains = domain.IpaDomain.RealmDomains
	}
}

func (p domainPresenter) createRhelIdm(
	output *public.Domain,
	domain *model.Domain,
) error {
	if err := p.createRhelIdmCheckDomain(domain); err != nil {
		return err
	}
	output.RhelIdm = &public.DomainIpa{}
	output.RhelIdm.RealmName = *domain.IpaDomain.RealmName
	p.fillRhelIdmCerts(output, domain)
	p.fillRhelIdmServers(output, domain)
	p.createRhelIdmFillRealmDomains(output, domain)
	return nil
}

func (p domainPresenter) guardSharedDomain(
	domain *model.Domain,
) error {
	if domain == nil {
		return fmt.Errorf("'domain' is nil")
	}
	if domain.Type == nil {
		return fmt.Errorf("'domain.Type' is nil")
	}
	if *domain.Type == model.DomainTypeUndefined {
		return fmt.Errorf("'domain.Type' is invalid")
	}
	return nil
}

func (p domainPresenter) sharedDomain(
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
		output.Type = model.DomainTypeIpaString
		output.RhelIdm = &public.DomainIpa{}
		err = p.sharedDomainFillRhelIdm(domain, output)
	default:
		err = fmt.Errorf("'domain.Type=%d' is invalid", *domain.Type)
	}
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (p domainPresenter) sharedDomainFill(
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

func (p domainPresenter) sharedDomainFillRhelIdm(
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
		return fmt.Errorf("'domain.IpaDomain' is nil")
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

	// if err = p.sharedDomainFillRhelIdmCaCerts(
	// 	domain, output,
	// ); err != nil {
	// 	return err
	// }
	p.fillRhelIdmCerts(output, domain)

	// if err = p.sharedDomainFillRhelIdmServers(
	// 	domain, output,
	// ); err != nil {
	// 	return err
	// }
	p.fillRhelIdmServers(output, domain)

	return nil
}
