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
	"go.openly.dev/pointy"
)

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteListDomains struct {
	SuiteBaseWithDomain
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
	var (
		token  *public.DomainRegToken
		domain *public.Domain
		err    error
	)
	s.SuiteBase.SetupTest()

	s.Domains = []*public.Domain{}
	for i := 1; i < 50; i++ {
		domainName := fmt.Sprintf("domain%d.test", i)
		s.As(XRHIDUser)
		if token, err = s.CreateToken(); err != nil {
			s.FailNow("error creating token")
		}
		s.As(XRHIDSystem)
		domainRequest := builder_api.NewDomain(domainName).Build()
		setFirstAsUpdateServer(domainRequest)
		setFirstServerRHSMId(s.T(), domainRequest, s.systemXRHID)
		domain, err = s.RegisterIpaDomain(token.DomainToken, domainRequest)
		if err != nil {
			s.FailNow("error registering domain")
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

func (s *SuiteListDomains) existDomain(id string) bool {
	for j := range s.Domains {
		if s.Domains[j].DomainId.String() == id {
			return true
		}
	}
	return false
}

func (s *SuiteListDomains) assertInDomains(t *testing.T, data []public.ListDomainsData, msgAndArgs ...any) bool {
	if data == nil {
		return true
	}
	if len(data) == 0 {
		return true
	}

	for i := range data {
		DomainIDString := data[i].DomainId.String()
		if !s.existDomain(DomainIDString) {
			return assert.Fail(t, fmt.Sprintf("Not in slice: DomainID=%s\n", DomainIDString), msgAndArgs...)
		}
	}

	return true
}

func (s *SuiteListDomains) TestListDomains() {
	t := s.T()

	xrhids := []XRHIDProfile{XRHIDUser, XRHIDServiceAccount}

	req, err := http.NewRequest(http.MethodGet, s.DefaultPublicBaseURL()+"/domains", nil)
	require.NoError(t, err)
	q := req.URL.Query()
	q.Add("offset", "0")
	q.Add("limit", "10")
	url1 := req.URL.String() + "?" + q.Encode()
	q.Set("offset", "40")
	q.Set("limit", "10")
	url2 := req.URL.String() + "?" + q.Encode()
	q.Set("offset", "20")
	q.Set("limit", "10")
	url3 := req.URL.String() + "?" + q.Encode()
	q.Del("offset")
	q.Del("limit")
	url4 := req.URL.String() + "?" + q.Encode()
	domainName := builder_helper.GenRandDomainName(2)

	// Prepare the tests
	testCases := []TestCase{
		{
			Name: "TestListDomains: offset=0&limit=10 case",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url1,
				Header: http.Header{
					header.HeaderXRequestID: {"test_domains_list_1"},
				},
				Body: builder_api.NewDomain(domainName).Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					header.HeaderXRHID: nil,
				},
				BodyFunc: WrapBodyFuncListDomainsResponse(func(t *testing.T, body *public.ListDomainsResponse) error {
					require.NotNil(t, body)

					// Check Meta
					assert.Equal(t, public.PaginationMeta{Count: int64(len(s.Domains)), Limit: 10, Offset: 0}, body.Meta)

					// Check links
					assert.Equal(t, public.PaginationLinks{
						First:    pointy.String("/api/idmsvc/v1/domains?limit=10&offset=0"),
						Previous: nil,
						Next:     pointy.String("/api/idmsvc/v1/domains?limit=10&offset=10"),
						Last:     pointy.String("/api/idmsvc/v1/domains?limit=10&offset=40"),
					}, body.Links)

					// Check items
					assert.Equal(t, 10, len(body.Data))
					s.assertInDomains(t, body.Data)

					return nil
				}),
			},
		},
		{
			Name: "TestListDomains: offset=40&limit=10 case",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url2,
				Header: http.Header{
					header.HeaderXRequestID: {"test_domains_list_2"},
				},
				Body: builder_api.NewDomain(domainName).Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					header.HeaderXRHID: nil,
					// TODO Check format for X-Rh-Idm-Version
				},
				BodyFunc: WrapBodyFuncListDomainsResponse(func(t *testing.T, body *public.ListDomainsResponse) error {
					require.NotNil(t, body)

					// Check Meta
					assert.Equal(t, public.PaginationMeta{Count: int64(len(s.Domains)), Limit: 10, Offset: 40}, body.Meta)

					// Check links
					assert.Equal(t, public.PaginationLinks{
						First:    pointy.String("/api/idmsvc/v1/domains?limit=10&offset=0"),
						Previous: pointy.String("/api/idmsvc/v1/domains?limit=10&offset=30"),
						Next:     nil,
						Last:     pointy.String("/api/idmsvc/v1/domains?limit=10&offset=40"),
					}, body.Links)

					// Check items
					assert.Equal(t, 9, len(body.Data))
					s.assertInDomains(t, body.Data)

					return nil
				}),
			},
		},
		{
			Name: "TestListDomains: offset=20&limit=10 case",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url3,
				Header: http.Header{
					header.HeaderXRequestID: {"test_domains_list_3"},
				},
				Body: builder_api.NewDomain(domainName).Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					header.HeaderXRHID: nil,
					// TODO Check format for X-Rh-Idm-Version
				},
				BodyFunc: WrapBodyFuncListDomainsResponse(func(t *testing.T, body *public.ListDomainsResponse) error {
					require.NotNil(t, body)

					// Check Meta
					assert.Equal(t, public.PaginationMeta{Count: int64(len(s.Domains)), Limit: 10, Offset: 20}, body.Meta)

					// Check links
					assert.Equal(t, public.PaginationLinks{
						First:    pointy.String("/api/idmsvc/v1/domains?limit=10&offset=0"),
						Previous: pointy.String("/api/idmsvc/v1/domains?limit=10&offset=10"),
						Next:     pointy.String("/api/idmsvc/v1/domains?limit=10&offset=30"),
						Last:     pointy.String("/api/idmsvc/v1/domains?limit=10&offset=40"),
					}, body.Links)

					// Check items
					assert.Equal(t, 10, len(body.Data))
					s.assertInDomains(t, body.Data)

					return nil
				}),
			},
		},
		{
			Name: "TestListDomains: no params",
			Given: TestCaseGiven{
				Method: http.MethodGet,
				URL:    url4,
				Header: http.Header{
					header.HeaderXRequestID: {"test_domains_list_4"},
				},
				Body: builder_api.NewDomain(domainName).Build(),
			},
			Expected: TestCaseExpect{
				StatusCode: http.StatusOK,
				Header: http.Header{
					// FIXME Avoid hardcode the key name of the header
					header.HeaderXRHID: nil,
					// TODO Check format for X-Rh-Idm-Version
				},
				BodyFunc: WrapBodyFuncListDomainsResponse(func(t *testing.T, body *public.ListDomainsResponse) error {
					require.NotNil(t, body)

					// Check Meta
					assert.Equal(t, public.PaginationMeta{Count: int64(len(s.Domains)), Limit: 10, Offset: 0}, body.Meta)

					// Check links
					assert.Equal(t, public.PaginationLinks{
						First:    pointy.String("/api/idmsvc/v1/domains?limit=10&offset=0"),
						Previous: nil,
						Next:     pointy.String("/api/idmsvc/v1/domains?limit=10&offset=10"),
						Last:     pointy.String("/api/idmsvc/v1/domains?limit=10&offset=40"),
					}, body.Links)

					// Check items
					assert.Equal(t, 10, len(body.Data))
					s.assertInDomains(t, body.Data)

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
