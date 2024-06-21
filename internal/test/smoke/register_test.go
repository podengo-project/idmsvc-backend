package smoke

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"github.com/stretchr/testify/assert"
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
	s.As(XRHIDUser)
	if s.token, err = s.CreateToken(); err != nil {
		s.FailNow("Error creating a token for registering a rhel-idm domain", "%s", err.Error())
	}
}

func (s *SuiteRegisterDomain) TearDownTest() {
	s.SuiteBase.TearDownTest()
}

func (s *SuiteRegisterDomain) TestRegisterDomain() {
	xrhidEncoded := header.EncodeXRHID(&s.systemXRHID)
	url := s.DefaultPublicBaseURL() + "/domains"
	domainName := builder_helper.GenRandDomainName(2)
	bodyRequest := builder_api.
		NewDomain(domainName).
		Build()

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestRegisterDomain rhel-idm",
			Given: TestCaseGiven{
				XRHIDProfile: XRHIDSystem,
				Method:       http.MethodPost,
				URL:          url,
				Header: http.Header{
					header.HeaderXRequestID:              {"test_token"},
					header.HeaderXRHID:                   {xrhidEncoded},
					header.HeaderXRHIDMRegistrationToken: {s.token.DomainToken},
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
				Body: bodyRequest,
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusCreated,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					header.HeaderXRequestID: {"test_token"},
					header.HeaderXRHID:      nil,
				},
				BodyFunc: WrapBodyFuncDomainResponse(func(t *testing.T, body *public.Domain) error {
					require.NotNil(t, body)
					require.NotNil(t, body.DomainId)
					assert.Equal(t, bodyRequest.DomainName, body.DomainName)
					assert.Equal(t, bodyRequest.DomainType, body.DomainType)

					if bodyRequest.Title != nil {
						require.NotNil(t, body.Title)
						assert.Equal(t, *bodyRequest.Title, *body.Title)
					} else {
						assert.Nil(t, body.Title)
					}

					require.NotNil(t, body.Description)
					assert.Equal(t, "", *body.Description)

					require.NotNil(t, body.AutoEnrollmentEnabled)
					assert.False(t, *body.AutoEnrollmentEnabled)

					// Check rhel-idm
					if bodyRequest != nil {
						require.NotNil(t, body.RhelIdm)
						assert.Equal(t, bodyRequest.RhelIdm.RealmName, body.RhelIdm.RealmName)
						assert.Equal(t, bodyRequest.RhelIdm.RealmDomains, body.RhelIdm.RealmDomains)
						if bodyRequest.RhelIdm.AutomountLocations != nil && *bodyRequest.RhelIdm.AutomountLocations != nil && len(*bodyRequest.RhelIdm.AutomountLocations) > 0 {
							assert.Equal(t, *bodyRequest.RhelIdm.AutomountLocations, body.RhelIdm.AutomountLocations)
						}
						if bodyRequest.RhelIdm.Locations != nil && len(bodyRequest.RhelIdm.Locations) > 0 {
							assert.Equal(t, bodyRequest.RhelIdm.Locations, body.RhelIdm.Locations)
						} else {
							assert.Condition(t, func() (success bool) {
								return body.RhelIdm.Locations == nil || len(body.RhelIdm.Locations) == 0
							})
						}
						assert.Equal(t, bodyRequest.RhelIdm.CaCerts, body.RhelIdm.CaCerts)
						assert.Equal(t, bodyRequest.RhelIdm.Servers, body.RhelIdm.Servers)
					}

					return nil
				}),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
