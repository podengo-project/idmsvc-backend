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
	SuiteBaseWithDomain
}

func (s *SuiteReadDomain) SetupTest() {
	s.SuiteBaseWithDomain.SetupTest()
}

func (s *SuiteReadDomain) TearDownTest() {
	s.SuiteBaseWithDomain.TearDownTest()
}

func (s *SuiteReadDomain) TestReadDomain() {
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "domains", s.Domains[0].DomainId.String())
	domainName := builder_helper.GenRandDomainName(2)
	xrhids := []XRHIDProfile{XRHIDUser, XRHIDServiceAccount, XRHIDSystem}

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestReadDomain",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url,
				Header: http.Header{
					header.HeaderXRequestID: {"test_read"},
				},
				Body: builder_api.NewDomain(domainName).Build(),
			},
			Expected: TestCaseExpect{
				// FIXME It must be http.StatusCreated
				StatusCode: http.StatusOK,
				Header: http.Header{
					header.HeaderXRHID: nil,
				},
				BodyFunc: WrapBodyFuncDomainResponse(func(t *testing.T, body *public.Domain) error {
					test_assert.AssertDomain(t, s.Domains[0], body)
					assert.Equal(t, s.Domains[0].DomainId, body.DomainId)
					return nil
				}),
			},
		},
	}

	for _, xrhid := range xrhids {
		for i := range testCases {
			testCases[i].Given.XRHIDProfile = xrhid
		}
		// Execute the test cases
		s.RunTestCases(testCases)
	}
}
