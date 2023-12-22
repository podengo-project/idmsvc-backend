package smoke

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"github.com/stretchr/testify/require"
)

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteRegisterDomain struct {
	SuiteBase
	token *public.DomainRegTokenResponse
}

// BodyFuncDomainResponse is the function that wrap
type BodyFuncDomainResponse func(t *testing.T, body *public.Domain) error

// WrapBodyFuncDomainResponse allow to implement custom body expectations for the specific type of the response.
// expected is the specific BodyFuncDomain for Domain type
// Returns a BodyFunc that wrap the generic expectation function.
func WrapBodyFuncDomainResponse(expected BodyFuncDomainResponse) BodyFunc {
	if expected == nil {
		return func(t *testing.T, body []byte) bool {
			return len(body) == 0
		}
	}
	return func(t *testing.T, body []byte) bool {
		// Unserialize the response to the expected type
		var data public.Domain
		if err := json.Unmarshal(body, &data); err != nil {
			require.Fail(t, fmt.Sprintf("Error unmarshalling body:\n"+
				"error: %q",
				err.Error(),
			))
			return false
		}

		// Run body expectetion on the unserialized data
		if err := expected(t, &data); err != nil {
			require.Fail(t, fmt.Sprintf("Error in body response:\n"+
				"error: %q",
				err.Error(),
			))
			return false
		}

		return true
	}
}

func (s *SuiteRegisterDomain) SetupTest() {
	var err error

	s.SuiteBase.SetupTest()

	// Get a token for the registration
	if s.token, err = s.CreateToken(); err != nil {
		s.FailNow("Error creating a token for registering a rhel-idm domain", "%s", err.Error())
	}
}

func (s *SuiteRegisterDomain) TearDownTest() {
	s.SuiteBase.TearDownTest()
}

// Specific expectation method that fit BodyFuncTokenResponse
func (s *SuiteRegisterDomain) bodyExpectationTestToken(t *testing.T, body *public.DomainRegTokenResponse) error {
	if body.DomainToken == "" {
		return fmt.Errorf("'domain_token' is empty")
	}

	if body.DomainType != "rhel-idm" {
		return fmt.Errorf("'domain_type' is not rhel-idm")
	}

	if body.DomainId == (uuid.UUID{}) {
		return fmt.Errorf("'domain_id' is empty")
	}

	if body.Expiration <= int(time.Now().Unix()) {
		return fmt.Errorf("'expiration' is in the past")
	}

	return nil
}

func (s *SuiteRegisterDomain) TestRegisterDomain() {
	xrhidEncoded := header.EncodeXRHID(&s.SystemXRHID)
	url := s.DefaultPublicBaseURL() + "/domains"
	domainName := builder_helper.GenRandomDomainName(2)

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestRegisterDomain rhel-idm",
			Given: TestCaseGiven{
				Method: http.MethodPost,
				URL:    url,
				Header: http.Header{
					"X-Rh-Insights-Request-Id":    {"test_token"},
					"X-Rh-Identity":               {xrhidEncoded},
					"X-Rh-Idm-Registration-Token": {s.token.DomainToken},
					"X-Rh-Idm-Version": {
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
				Body: builder_api.
					NewDomain(domainName).
					Build(),
			},
			Expected: TestCaseExpect{
				// FIXME It must be http.StatusCreated
				StatusCode: http.StatusOK,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					"X-Rh-Insights-Request-Id": {"test_token"},
					"X-Rh-Identity":            nil,
					// TODO Check format for X-Rh-Idm-Version
				},
				BodyFunc: WrapBodyFuncDomainResponse(func(t *testing.T, body *public.Domain) error {
					// TODO Add the checks here
					return nil
				}),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)

	// s.RunTestCase(&testCases[0])
}
