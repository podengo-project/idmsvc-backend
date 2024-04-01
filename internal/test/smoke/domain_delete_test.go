package smoke

import (
	"fmt"
	"net/http"

	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
)

// SuiteDeleteDomain is the suite to validate the smoke test when read domain endpoint at GET /api/idmsvc/v1/domains/:domain_id
type SuiteDeleteDomain struct {
	SuiteBaseWithDomain
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
