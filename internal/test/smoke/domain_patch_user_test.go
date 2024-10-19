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
	"github.com/stretchr/testify/suite"
)

// SuiteDomainUpdateUser is the suite to validate the smoke test when a user update the domain endpoint at PATCH /api/idmsvc/v1/domains/:domain_id
type SuiteDomainUpdateUser struct {
	SuiteBaseWithDomain
}

func (s *SuiteDomainUpdateUser) SetupTest() {
	s.SuiteBaseWithDomain.SetupTest()
}

func (s *SuiteDomainUpdateUser) TearDownTest() {
	s.SuiteBaseWithDomain.TearDownTest()
}

func (s *SuiteDomainUpdateUser) TestPatchDomain() {
	url := fmt.Sprintf("%s/%s/%s", s.DefaultPublicBaseURL(), "domains", s.Domains[0].DomainId.String())
	patchedDomain := builder_api.NewUpdateDomainUserRequest().Build()
	xrhids := []XRHIDProfile{XRHIDUser, XRHIDServiceAccount}

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestPatchDomain",
			Given: TestCaseGiven{
				Method: http.MethodPatch,
				URL:    url,
				Header: http.Header{
					header.HeaderXRequestID: {"test_domain_patch"},
				},
				Body: patchedDomain,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					header.HeaderXRHID: nil,
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
	for _, xrhid := range xrhids {
		for i := range testCases {
			testCases[i].Given.XRHIDProfile = xrhid
		}
		s.RunTestCases(testCases)
	}
}

func TestSuiteDomainUpdateUser(t *testing.T) {
	suite.Run(t, new(SuiteDomainUpdateUser))
}
