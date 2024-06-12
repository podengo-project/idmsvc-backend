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

// Specific expectation method that fit BodyFuncTokenResponse
func (s *SuiteTokenCreate) bodyExpectationTestToken(t *testing.T, body *public.DomainRegTokenResponse) error {
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

func (s *SuiteTokenCreate) TestToken() {
	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestToken",
			Given: TestCaseGiven{
				XRHIDProfile: XRHIDUser,
				Method:       http.MethodPost,
				URL:          s.DefaultPublicBaseURL() + "/domains/token",
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
					header.HeaderXRequestID: {"test_token"},
					header.HeaderXRHID:      nil,
				},
				BodyFunc: WrapBodyFuncTokenResponse(s.bodyExpectationTestToken),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
