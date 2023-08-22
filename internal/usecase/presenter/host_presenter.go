package presenter

import (
	"fmt"
	"strings"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/presenter"
)

type hostPresenter struct {
	cfg *config.Config
}

func NewHostPresenter(cfg *config.Config) presenter.HostPresenter {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	return &hostPresenter{cfg}
}

func (p *hostPresenter) fillRhelIdm(domain *model.Domain, response *public.HostConfResponse) error {
	// concatenate PEM certs
	if len(domain.IpaDomain.CaCerts) == 0 {
		return fmt.Errorf("domain '%s' has no CA certificates", *domain.DomainName)
	}
	var sb strings.Builder
	for _, ca_cert := range domain.IpaDomain.CaCerts {
		sb.WriteString(ca_cert.Pem)
		// ensure PEM blocks are separated by newline
		if !strings.HasSuffix(ca_cert.Pem, "\n") {
			sb.WriteString("\n")
		}
	}

	// create array of servers with HCC enrollment role
	var servers []public.HostConfIpaServer
	for _, ipa_server := range domain.IpaDomain.Servers {
		if ipa_server.HCCEnrollmentServer {
			servers = append(servers, public.HostConfIpaServer{Fqdn: ipa_server.FQDN, Location: ipa_server.Location})
		}
	}
	if len(servers) == 0 {
		return fmt.Errorf("domain '%s' has no enrollment servers", *domain.DomainName)
	}
	response.RhelIdm = public.HostConfIpa{
		Cabundle:          sb.String(),
		EnrollmentServers: servers,
		RealmName:         *domain.IpaDomain.RealmName,
	}
	return nil
}

func (p *hostPresenter) HostConf(domain *model.Domain) (*public.HostConfResponse, error) {
	var err error

	if domain == nil {
		return nil, fmt.Errorf("'domain' is nil")
	}
	if domain.Type == nil {
		return nil, fmt.Errorf("'domain.Type' is nil")
	}
	domainType := public.DomainType(model.DomainTypeString(*domain.Type))

	response := &public.HostConfResponse{
		AutoEnrollmentEnabled: *domain.AutoEnrollmentEnabled,
		DomainId:              domain.DomainUuid,
		DomainName:            *domain.DomainName,
		DomainType:            domainType,
	}

	switch *domain.Type {
	case model.DomainTypeIpa:
		err = p.fillRhelIdm(domain, response)
	default:
		err = fmt.Errorf("domain '%s' has unsupported domain type '%s'", *domain.DomainName, domainType)
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}
