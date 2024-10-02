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

func (s *SuiteDomainUpdateAgent) buildUpdateAgentRequest(domainName string) *public.UpdateDomainAgentRequest {
	return builder_api.NewUpdateDomainAgent(domainName).
		WithSubscriptionManagerID(s.systemXRHID.Identity.System.CommonName).
		WithHCCUpdate(true).
		Build()
}

func (s *SuiteDomainUpdateAgent) TestUpdateDomain() {
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "domains", s.Domains[0].DomainId)

	domainName := s.Domains[0].DomainName
	requestWithChangedDomainName := s.buildUpdateAgentRequest(domainName)
	requestWithChangedDomainName.DomainName = "other.domain.test"

	requestWithChangedRealm := s.buildUpdateAgentRequest(domainName)
	requestWithChangedRealm.RhelIdm.RealmName = "DIFFERENT.REALM"

	requestWithBadSubscriptionManagerID := builder_api.NewUpdateDomainAgent(domainName).Build()

	okRequest := s.buildUpdateAgentRequest(domainName)

	expectedResponse := s.Domains[0]
	expectedResponse.RhelIdm = &okRequest.RhelIdm

	test_header := http.Header{
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
	}

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestPutDomainWithChangedDomainName",
			Given: TestCaseGiven{
				XRHIDProfile: XRHIDSystem,
				Method:       http.MethodPut,
				URL:          url,
				Header:       test_header,
				Body:         requestWithChangedDomainName,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusBadRequest,
				BodyFunc: WrapBodyFuncErrorResponse(func(t *testing.T, body *public.ErrorResponse) error {
					assert.Equal(t, builder_api.NewErrorResponse().
						Add(*builder_api.NewErrorInfo(http.StatusBadRequest).
							WithTitle("'domain_name' may not be changed").
							Build()).
						Build(), body)
					return nil
				}),
			},
		},
		{
			Name: "TestPutDomainWithChangedRealm",
			Given: TestCaseGiven{
				XRHIDProfile: XRHIDSystem,
				Method:       http.MethodPut,
				URL:          url,
				Header:       test_header,
				Body:         requestWithChangedRealm,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusBadRequest,
				BodyFunc: WrapBodyFuncErrorResponse(func(t *testing.T, body *public.ErrorResponse) error {
					assert.Equal(t, builder_api.NewErrorResponse().
						Add(*builder_api.NewErrorInfo(http.StatusBadRequest).
							WithTitle("'realm_name' may not be changed").
							Build()).
						Build(), body)
					return nil
				}),
			},
		},
		{
			Name: "TestPutDomainWithBadSubscriptionManagerID",
			Given: TestCaseGiven{
				XRHIDProfile: XRHIDSystem,
				Method:       http.MethodPut,
				URL:          url,
				Header:       test_header,
				Body:         requestWithBadSubscriptionManagerID,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusBadRequest,
				BodyFunc: WrapBodyFuncErrorResponse(func(t *testing.T, body *public.ErrorResponse) error {
					assert.Equal(t, builder_api.NewErrorResponse().
						Add(*builder_api.NewErrorInfo(http.StatusBadRequest).
							WithTitle("update server's 'Subscription Manager ID' not found in the authorized list of rhel-idm servers").
							Build()).
						Build(), body)
					return nil
				}),
			},
		},
		{
			Name: "TestPutDomain",
			Given: TestCaseGiven{
				XRHIDProfile: XRHIDSystem,
				Method:       http.MethodPut,
				URL:          url,
				Header:       test_header,
				Body:         okRequest,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					header.HeaderXRHID: nil,
				},
				BodyFunc: WrapBodyFuncDomainResponse(func(t *testing.T, body *public.Domain) error {
					test_assert.AssertDomain(t, expectedResponse, body)
					assert.Equal(t, s.Domains[0].DomainId, body.DomainId)
					return nil
				}),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
