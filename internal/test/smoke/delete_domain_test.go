package smoke

import (
	"fmt"
	"net/http"

	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
)

// SuiteDeleteDomain is the suite to validate the smoke test when read domain endpoint at GET /api/idmsvc/v1/domains/:domain_id
type SuiteDeleteDomain struct {
	SuiteReadDomain
	Domains []*public.Domain
}

func (s *SuiteDeleteDomain) SetupTest() {
	s.SuiteReadDomain.SetupTest()

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

func (s *SuiteDeleteDomain) TearDownTest() {
	for i := range s.Domains {
		s.Domains[i] = nil
	}
	s.Domains = nil

	s.SuiteBase.TearDownTest()
}

func (s *SuiteDeleteDomain) TestReadDomain() {
	t := s.T()
	t.Skip("Skipping wrapped test to avoid duplication")
}

func (s *SuiteDeleteDomain) TestDeleteDomain() {
	xrhidEncoded := header.EncodeXRHID(&s.UserXRHID)
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "domains", s.Domains[0].DomainId)
	domainName := builder_helper.GenRandDomainName(2)

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestDeleteDomain",
			Given: TestCaseGiven{
				Method: http.MethodDelete,
				URL:    url,
				Header: http.Header{
					header.HeaderXRequestID: {"test_token"},
					header.HeaderXRHID:      {xrhidEncoded},
				},
				Body: builder_api.NewDomain(domainName).Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusNoContent,
				Header: http.Header{
					header.HeaderXRequestID: {"test_token"},
					header.HeaderXRHID:      nil,
				},
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
