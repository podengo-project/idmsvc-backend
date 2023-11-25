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
type SuiteListDomains struct {
	SuiteBase
	Domains []*public.Domain
}

// BodyFuncListResponse is the function that wrap
type BodyFuncListDomainsResponse func(t *testing.T, body *public.ListDomainsResponse) error

// WrapBodyFuncListResponse allow to implement custom body expectations for the specific type of the response.
// expected is the specific BodyFuncDomain for Domain type
// Returns a BodyFunc that wrap the generic expectation function.
func WrapBodyFuncListDomainsResponse(expected BodyFuncListDomainsResponse) BodyFunc {
	if expected == nil {
		return func(t *testing.T, body []byte) bool {
			return len(body) == 0
		}
	}
	return func(t *testing.T, body []byte) bool {
		// Unserialize the response to the expected type
		var data public.ListDomainsResponse
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

func (s *SuiteListDomains) SetupTest() {
	s.SuiteBase.SetupTest()

	s.Domains = []*public.Domain{}
	for i := 1; i < 50; i++ {
		domainName := fmt.Sprintf("domain%d.test", i)
		domain, err := s.RegisterIpaDomain(builder_api.NewDomain(domainName).Build())
		if err != nil {
			s.FailNow("error creating ")
		}
		s.Domains = append(s.Domains, domain)
	}
}

func (s *SuiteListDomains) TearDownTest() {
	for i := range s.Domains {
		s.Domains[i] = nil
	}
	s.Domains = nil

	s.SuiteBase.TearDownTest()
}

func (s *SuiteListDomains) TestListDomains() {
	t := s.T()
	xrhidEncoded := header.EncodeXRHID(&s.UserXRHID)
	req, err := http.NewRequest(http.MethodGet, s.DefaultPublicBaseURL()+"/domains", nil)
	require.NoError(t, err)
	q := req.URL.Query()
	q.Add("offset", "0")
	q.Add("limit", "10")
	url1 := req.URL.String() + "?" + q.Encode()
	q.Set("offset", "40")
	q.Set("limit", "10")
	url2 := req.URL.String() + "?" + q.Encode()
	domainName := builder_helper.GenRandDomainName(2)

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestListDomains: offset=0&limit=10",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url1,
				Header: http.Header{
					"X-Rh-Insights-Request-Id": {"test_token"},
					"X-Rh-Identity":            {xrhidEncoded},
				},
				Body: builder_api.NewDomain(domainName).Build(),
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
				BodyFunc: WrapBodyFuncListDomainsResponse(func(t *testing.T, body *public.ListDomainsResponse) error {
					require.NotNil(t, body)
					assert.Equal(t, 10, body.Meta.Limit)
					assert.Equal(t, 0, body.Meta.Offset)
					assert.Equal(t, int64(49), body.Meta.Count)
					return nil
				}),
			},
		},
		{
			Name: "TestListDomains: offset=40&limit=10",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url2,
				Header: http.Header{
					"X-Rh-Insights-Request-Id": {"test_token"},
					"X-Rh-Identity":            {xrhidEncoded},
				},
				Body: builder_api.NewDomain(domainName).Build(),
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
				BodyFunc: WrapBodyFuncListDomainsResponse(func(t *testing.T, body *public.ListDomainsResponse) error {
					require.NotNil(t, body)
					assert.Equal(t, 10, body.Meta.Limit)
					assert.Equal(t, 40, body.Meta.Offset)
					// FIXME Review the metadata to return
					// assert.Equal(t, int64(49), body.Meta.Count)
					assert.Equal(t, 9, len(body.Data))
					// TODO
					return nil
				}),
			},
		},
	}

	// Execute the test cases
	s.RunTestCases(testCases)
}
