package smoke

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	test_assert "github.com/podengo-project/idmsvc-backend/internal/test/assert"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/stretchr/testify/assert"
)

// SuiteDomainUpdateAgent is the suite to validate the smoke test when read domain endpoint at PUT /api/idmsvc/v1/domains/:domain_id
type SuiteDomainUpdateAgent struct {
	SuiteBaseWithDomain
}

func (s *SuiteDomainUpdateAgent) SetupTest() {
	s.SuiteBase.SetupTest()

	var (
		domainName string
		domain     *public.Domain
		err        error
		i          int
		token      *public.DomainRegToken
	)

	i = 0
	s.Domains = []*public.Domain{}
	domainName = fmt.Sprintf("domain%d.test", i)
	newDomain := builder_api.NewDomain(domainName).Build()
	*newDomain.RhelIdm.Servers[0].SubscriptionManagerId = uuid.MustParse(s.systemXRHID.Identity.System.CommonName)
	newDomain.RhelIdm.Servers[0].HccUpdateServer = true
	s.As(XRHIDUser)
	if token, err = s.CreateToken(); err != nil {
		s.FailNow("error creating token")
	}
	s.As(XRHIDSystem)
	domain, err = s.RegisterIpaDomain(token.DomainToken, newDomain)
	if err != nil {
		s.FailNow("error registering rhel-idm domain")
	}
	s.Domains = append(s.Domains, domain)
}

func (s *SuiteDomainUpdateAgent) TearDownTest() {
	for i := range s.Domains {
		s.Domains[i] = nil
	}
	s.Domains = nil

	s.SuiteBase.TearDownTest()
}

func (s *SuiteDomainUpdateAgent) TestUpdateDomain() {
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "domains", s.Domains[0].DomainId)
	domainName := s.Domains[0].DomainName
	updatedDomain := builder_api.NewUpdateDomainAgent(domainName).WithSubscriptionManagerID(s.systemXRHID.Identity.System.CommonName).Build()
	expectedDomain := s.Domains[0]
	expectedDomain.RhelIdm = &updatedDomain.RhelIdm

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestReadDomain",
			Given: TestCaseGiven{
				XRHIDProfile: XRHIDSystem,
				Method:       http.MethodPut,
				URL:          url,
				Header: http.Header{
					header.HeaderXRequestID: {"test_domain_update"},
					header.HeaderXRHIDMVersion: {
						header.EncodeXRHIDMVersion(
							header.NewXRHIDMVersion(
								"v1.0.0",
								"4.19.0",
								"redhat-9.3",
								"9.3",
							),
						),
					},
				},
				Body: updatedDomain,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					header.HeaderXRequestID: {"test_domain_update"},
					header.HeaderXRHID:      nil,
				},
				BodyFunc: WrapBodyFuncDomainResponse(func(t *testing.T, body *public.Domain) error {
					test_assert.AssertDomain(t, s.Domains[0], body)
					assert.Equal(t, s.Domains[0].DomainId, body.DomainId)
					return nil
				}),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
