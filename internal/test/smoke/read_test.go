package smoke

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	test_assert "github.com/podengo-project/idmsvc-backend/internal/test/assert"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"github.com/stretchr/testify/assert"
)

// SuiteReadDomain is the suite to validate the smoke test when read domain endpoint at GET /api/idmsvc/v1/domains/:domain_id
type SuiteReadDomain struct {
	SuiteBase
	Domains []*public.Domain
}

func (s *SuiteReadDomain) SetupTest() {
	s.SuiteBase.SetupTest()

	var (
		domainName string
		domain     *public.Domain
		err        error
		i          int
	)

	// Domain 1 in OrgID1
	i = 0
	s.Domains = []*public.Domain{}
	domainName = fmt.Sprintf("domain%d.test", i)
	domain, err = s.RegisterIpaDomain(builder_api.NewDomain(domainName).Build())
	if err != nil {
		s.FailNow("error creating ")
	}
	s.Domains = append(s.Domains, domain)
}

func (s *SuiteReadDomain) TearDownTest() {
	for i := range s.Domains {
		s.Domains[i] = nil
	}
	s.Domains = nil

	s.SuiteBase.TearDownTest()
}

func (s *SuiteReadDomain) TestReadDomain() {
	xrhidEncoded := header.EncodeXRHID(&s.UserXRHID)
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "domains", s.Domains[0].DomainId)
	domainName := builder_helper.GenRandDomainName(2)

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestReadDomain",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url,
				Header: http.Header{
					header.HeaderXRequestID: {"test_token"},
					header.HeaderXRHID:      {xrhidEncoded},
				},
				Body: builder_api.NewDomain(domainName).Build(),
			},
			Expected: TestCaseExpect{
				// FIXME It must be http.StatusCreated
				StatusCode: http.StatusOK,
				Header: http.Header{
					header.HeaderXRequestID: {"test_token"},
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
