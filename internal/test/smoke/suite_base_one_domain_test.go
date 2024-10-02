package smoke

import (
	"fmt"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
)

// SuiteBaseWithDomain is the suite to validate the smoke test when read domain endpoint at GET /api/idmsvc/v1/domains/:domain_id
type SuiteBaseWithDomain struct {
	SuiteBase
	Domains []*public.Domain
}

func (s *SuiteBaseWithDomain) SetupTest() {
	s.SuiteBase.SetupTest()

	var (
		domainName string
		token      *public.DomainRegToken
		domain     *public.Domain
		err        error
		i          int
	)

	// Domain 1 in OrgID1
	i = 0
	s.Domains = []*public.Domain{}
	// As an admin user
	s.As(RBACAdmin, XRHIDUser)
	if token, err = s.CreateToken(); err != nil {
		s.FailNow("error creating token")
	}
	domainName = fmt.Sprintf("domain%d.test", i)
	domainRequest := builder_api.NewDomain(domainName).Build()
	setFirstServerRHSMId(s.T(), domainRequest, s.systemXRHID)
	setFirstAsUpdateServer(domainRequest)
	s.As(XRHIDSystem)
	if domain, err = s.RegisterIpaDomain(token.DomainToken, domainRequest); err != nil {
		s.FailNow("error creating rhel-idm domain")
	}
	s.Domains = append(s.Domains, domain)
}

func (s *SuiteBaseWithDomain) TearDownTest() {
	for i := range s.Domains {
		s.Domains[i] = nil
	}
	s.Domains = nil
	s.SuiteBase.TearDownTest()
}
