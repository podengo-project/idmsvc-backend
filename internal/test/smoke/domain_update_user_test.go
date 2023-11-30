package smoke

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteDomainUpdateUser is the suite to validate the smoke test when a user update the domain endpoint at PATCH /api/idmsvc/v1/domains/:domain_id
type SuiteDomainUpdateUser struct {
	SuiteReadDomain
}

func (s *SuiteDomainUpdateUser) SetupTest() {
	s.SuiteReadDomain.SetupTest()
}

func (s *SuiteDomainUpdateUser) TearDownTest() {
	s.SuiteReadDomain.TearDownTest()
}

func (s *SuiteDomainUpdateUser) TestReadDomain() {
	// This avoid the test from the wrapped suite is executed, so it runs only
	// once
	t := s.T()
	t.Skip("skipping parent duplicated test")
}

func (s *SuiteDomainUpdateUser) TestPatchDomain() {
	xrhidEncoded := header.EncodeXRHID(&s.UserXRHID)
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "domains", s.Domains[0].DomainId.String())
	patchedDomain := builder_api.NewUpdateDomainUserJSONRequestBody().Build()

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestPatchDomain",
			Given: TestCaseGiven{
				Method: http.MethodPatch,
				URL:    url,
				Header: http.Header{
					"X-Rh-Insights-Request-Id": {"test_domain_patch"},
					"X-Rh-Identity":            {xrhidEncoded},
				},
				Body: patchedDomain,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"X-Rh-Insights-Request-Id": {"test_domain_patch"},
					"X-Rh-Identity":            nil,
				},
				BodyFunc: WrapBodyFuncDomainResponse(func(t *testing.T, body *public.Domain) error {
					require.NotNil(t, body)
					if patchedDomain.Title != nil {
						assert.Equal(t, patchedDomain.Title, body.Title)
					}
					if patchedDomain.Description != nil {
						assert.Equal(t, patchedDomain.Description, body.Description)
					}
					if patchedDomain.AutoEnrollmentEnabled != nil {
						assert.Equal(t, patchedDomain.AutoEnrollmentEnabled, body.AutoEnrollmentEnabled)
					}
					return nil
				}),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
