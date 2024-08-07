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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteTokenCreate struct {
	SuiteBase
}

// BodyFuncTokenResponse is the function that wrap
type BodyFuncTokenResponse func(t *testing.T, expect *public.DomainRegTokenResponse) error

// WrapBodyFuncTokenResponse allow to implement custom body expectations for the specific type of the response.
// expected is the specific BodyFuncTokenResponse for DomainRegTokenResponse type
// Returns a BodyFunc that wrap the generic expectation function.
func WrapBodyFuncTokenResponse(expected BodyFuncTokenResponse) BodyFunc {
	// To allow a generic interface for any body response type
	// I have to use `body []byte`; I cannot use `any` because
	// the response type is particular for the endpoint.
	// That means the input to the function is not in a golang
	// structure; to let the tests to be defined with less boilerplate,
	// every response type would implement a wrapper function like
	// this, which unmarshall the bytes, and call to the more specific
	// custom body function.
	if expected == nil {
		return func(t *testing.T, body []byte) bool {
			return len(body) == 0
		}
	}
	return func(t *testing.T, body []byte) bool {
		// Unserialize the response to the expected type
		var data public.DomainRegTokenResponse
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

func (s *SuiteTokenCreate) TestToken() {
	xrhidSlice := []XRHIDProfile{XRHIDUser, XRHIDServiceAccount}

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestToken",
			Given: TestCaseGiven{
				Method: http.MethodPost,
				URL:    s.DefaultPublicBaseURL() + "/domains/token",
				Header: http.Header{
					header.HeaderXRequestID: {"test_token"},
				},
				Body: public.DomainRegTokenRequest{
					DomainType: "rhel-idm",
				},
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					header.HeaderXRHID: nil,
				},
				BodyFunc: WrapBodyFuncTokenResponse(func(t *testing.T, body *public.DomainRegTokenResponse) error {
					assert.NotEmpty(t, body.DomainToken)
					assert.Equal(t, public.DomainType("rhel-idm"), body.DomainType)
					assert.NotEqual(t, uuid.UUID{}, body.DomainId)
					assert.True(t, int(time.Now().Unix()) < body.Expiration)
					return nil
				}),
			},
		},
	}

	// Run for users and service accounts
	s.As(RBACAdmin)
	for _, xrhid := range xrhidSlice {
		for i := range testCases {
			testCases[i].Given.XRHIDProfile = xrhid
		}
		// Execute the test cases
		s.RunTestCases(testCases)
	}
}
