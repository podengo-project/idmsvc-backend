package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPath = "/test"

func helperCreatePredicate(username string) IdentityPredicate {
	return func(data *identity.XRHID) error {
		if data == nil {
			return internal_errors.NilArgError("data")
		}
		if data.Identity.User.Username == username {
			return fmt.Errorf("username='%s' is not accepted", username)
		}
		return nil
	}
}

func helperNewEchoEnforceIdentity(m echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	}
	e.Use(CreateContext())
	e.Use(m)
	e.Add("GET", testPath, h)

	return e
}

// FIXME
func helperGenerateUserIdentity(orgId string, username string) *identity.XRHID {
	return &identity.XRHID{
		Identity: identity.Identity{
			AccountNumber: "12345",
			OrgID:         orgId,
			Internal: identity.Internal{
				OrgID: orgId,
			},
			Type: "User",
			User: &identity.User{
				Username: username,
				UserID:   "12345",
				Active:   true,
				Internal: true,
				OrgAdmin: true,
				Locale:   "en",
			},
		},
	}
}

func helperGenerateSystemIdentity(orgId string, commonName string) *identity.XRHID {
	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/identities/cert.json
	return &identity.XRHID{
		Identity: identity.Identity{
			OrgID:         orgId,
			AccountNumber: "11111",
			AuthType:      "cert-auth",
			Type:          "System",
			Internal: identity.Internal{
				OrgID: orgId,
			},
			System: &identity.System{
				CommonName: commonName,
				CertType:   "system",
			},
		},
	}
}

func helperSkipper(data bool) echo_middleware.Skipper {
	return func(c echo.Context) bool {
		return data
	}
}

func TestEnforceIdentityWithConfigPanic(t *testing.T) {
	assert.Panics(t, func() {
		EnforceIdentityWithConfig(nil)
	})
}

func TestPredicateIdentityAlwaysTrue(t *testing.T) {
	assert.Nil(t, IdentityAlwaysTrue(nil))
}

func TestEnforceIdentity(t *testing.T) {
	// TODO Double check if the http response code are
	//      the expected, or if it has to be changed to 403 or 4
	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/identities/cert.json
	type TestCaseExpected struct {
		Code int
		Body string
	}
	type TestCase struct {
		Name     string
		Given    *string
		Expected TestCaseExpected
	}

	testCases := []TestCase{
		{
			Name:  header.HeaderXRHID + " header not present",
			Given: nil,
			Expected: TestCaseExpected{
				Code: http.StatusBadRequest,
				Body: "{\"message\":\"Bad Request\"}\n",
			},
		},
		{
			Name:  header.HeaderXRHID + " bad base64 coding",
			Given: pointy.String("bad base64 coding"),
			Expected: TestCaseExpected{
				Code: http.StatusBadRequest,
				Body: "{\"message\":\"Bad Request\"}\n",
			},
		},
		{
			Name:  header.HeaderXRHID + " bad json encoding",
			Given: pointy.String("ewo="),
			Expected: TestCaseExpected{
				Code: http.StatusBadRequest,
				Body: "{\"message\":\"Bad Request\"}\n",
			},
		},
		{
			Name: header.HeaderXRHID + " fail predicates",
			Given: pointy.String(
				header.EncodeXRHID(
					helperGenerateUserIdentity("12345", "test-fail-predicate"),
				),
			),
			Expected: TestCaseExpected{
				Code: http.StatusUnauthorized,
				Body: "{\"message\":\"Unauthorized\"}\n",
			},
		},
		{
			Name: header.HeaderXRHID + " pass predicates",
			Given: pointy.String(
				header.EncodeXRHID(
					helperGenerateUserIdentity("12345", "testuser"),
				),
			),
			Expected: TestCaseExpected{
				Code: http.StatusOK,
				Body: "Ok",
			},
		},
		{
			Name: header.HeaderXRHID + " pass predicates",
			Given: pointy.String(
				header.EncodeXRHID(
					helperGenerateSystemIdentity("12345", "testuser"),
				),
			),
			Expected: TestCaseExpected{
				Code: http.StatusOK,
				Body: "Ok",
			},
		},
	}

	// Get echo instance with the middleware and one predicate for test it
	e := helperNewEchoEnforceIdentity(
		EnforceIdentityWithConfig(
			&IdentityConfig{
				Predicates: []IdentityPredicateEntry{
					{
						Name:      "test-fail-predicate",
						Predicate: helperCreatePredicate("test-fail-predicate"),
					},
				},
			},
		))
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		if testCase.Given != nil {
			req.Header.Add(header.HeaderXRHID, *testCase.Given)
		}
		e.ServeHTTP(res, req)

		// Check expectations
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, testCase.Expected.Code, res.Code)
		assert.Equal(t, testCase.Expected.Body, string(data))
	}
}

func TestEnforceIdentityNoDomainContext(t *testing.T) {
	e := echo.New()
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	}
	e.Use(
		EnforceIdentityWithConfig(
			&IdentityConfig{
				Predicates: []IdentityPredicateEntry{
					{
						Name:      "test-fail-predicate",
						Predicate: helperCreatePredicate("test-fail-predicate"),
					},
				},
			},
		))
	e.Add("GET", testPath, h)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	iden := header.EncodeXRHID(&identity.XRHID{})
	req.Header.Add(header.HeaderXRHID, iden)
	e.ServeHTTP(res, req)

	// Check expectations
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t,
		fmt.Sprintf(
			`{"message":"%s"}
`,
			http.StatusText(http.StatusInternalServerError),
		),
		string(data),
	)
}

func TestEnforceIdentitySkipper(t *testing.T) {
	var (
		e    *echo.Echo
		res  *httptest.ResponseRecorder
		req  *http.Request
		data []byte
		err  error
	)

	// When skipper return false, as no x-rh-identity provided, will return unauthorized
	e = helperNewEchoEnforceIdentity(
		EnforceIdentityWithConfig(
			&IdentityConfig{
				Skipper: helperSkipper(false),
			},
		),
	)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	e.ServeHTTP(res, req)
	// Check expectations
	data, err = io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "{\"message\":\"Bad Request\"}\n", string(data))

	// When skipper return true the middleware does not process the header or the predicates
	e = helperNewEchoEnforceIdentity(
		EnforceIdentityWithConfig(
			&IdentityConfig{
				Skipper: helperSkipper(true),
			},
		),
	)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	e.ServeHTTP(res, req)
	// Check expectations
	data, err = io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "Ok", string(data))
}

func TestEnforceUserPredicate(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    *identity.XRHID
		Expected error
	}
	testCases := []TestCase{
		{
			Name:     "nil argument",
			Given:    nil,
			Expected: fmt.Errorf("'data' cannot be nil"),
		},
		{
			Name: "Identity type is not 'User'",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "System",
				},
			},
			Expected: fmt.Errorf("'Identity.Type=System' is not 'User'"),
		},
		{
			Name: "Identity with disabled user",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
					User: &identity.User{
						Active: false,
					},
				},
			},
			Expected: fmt.Errorf("'Identity.User.Active' is not true"),
		},
		{
			Name: "'UserName' is empty",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
					User: &identity.User{
						Active:   true,
						UserID:   "jdoe",
						Username: "",
					},
				},
			},
			Expected: fmt.Errorf("'Identity.User.Username' cannot be empty"),
		},
		{
			Name: "Success case",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
					User: &identity.User{
						Active:   true,
						UserID:   "jdoe",
						Username: "jdoe",
					},
				},
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		err := EnforceUserPredicate(testCase.Given)
		if testCase.Expected != nil {
			require.NotNil(t, err)
			assert.EqualError(t, err, testCase.Expected.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestEnforceSystemPredicate(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    *identity.XRHID
		Expected error
	}
	testCases := []TestCase{
		{
			Name:     "nil argument",
			Given:    nil,
			Expected: fmt.Errorf("'data' cannot be nil"),
		},
		{
			Name: "'Identity' type is not 'System'",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
				},
			},
			Expected: fmt.Errorf("'Identity.Type' must be 'System'"),
		},
		{
			Name: "'CertType' is not 'system'",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "System",
					System: &identity.System{
						CertType: "anothevalue",
					},
				},
			},
			Expected: fmt.Errorf("'Identity.System.CertType' is not 'system'"),
		},
		{
			Name: "'CommonName' is empty",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "System",
					System: &identity.System{
						CertType:   "system",
						CommonName: "",
					},
				},
			},
			Expected: fmt.Errorf("'Identity.System.CommonName' is empty"),
		},
		{
			Name: "Success case",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "System",
					System: &identity.System{
						CertType:   "system",
						CommonName: "10fbb716-ca5d-11ed-b384-482ae3863d30",
					},
				},
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		err := EnforceSystemPredicate(testCase.Given)
		if testCase.Expected != nil {
			require.NotNil(t, err)
			assert.EqualError(t, err, testCase.Expected.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestEnforceIdentityOrder(t *testing.T) {
	var order map[string]time.Time = map[string]time.Time{
		"first":  {},
		"second": {},
	}

	// Get echo instance with the middleware and one predicate for test it
	e := helperNewEchoEnforceIdentity(
		EnforceIdentityWithConfig(
			&IdentityConfig{
				Predicates: []IdentityPredicateEntry{
					{
						Name: "first",
						Predicate: func(data *identity.XRHID) error {
							order["first"] = time.Now().UTC()
							return nil
						},
					},
					{
						Name: "second",
						Predicate: func(data *identity.XRHID) error {
							order["second"] = time.Now().UTC()
							return nil
						},
					},
				},
			},
		))
	// xrhid := `{"identity":{"org_id":"12345","internal":{"org_id":"12345"},"user":{"username":"sapheaded","email":"hooked@bought.biz","first_name":"Leslie","last_name":"Jacobs","is_active":false,"is_org_admin":false,"is_internal":false,"locale":"km","user_id":"jeweljeweler"},"system":{},"associate":{"Role":null,"email":"","givenName":"","rhatUUID":"","surname":""},"x509":{"subject_dn":"","issuer_dn":""},"service_account":{"client_id":"","username":""},"type":"User","auth_type":"basic-auth"},"entitlements":null}`
	xrhid := header.EncodeXRHID(helperGenerateUserIdentity("12345", "test"))
	for i := 0; i < 1000; i++ {
		order["first"] = time.Time{}
		order["second"] = time.Time{}
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Add(header.HeaderXRHID, xrhid)
		e.ServeHTTP(res, req)

		// Check expectations
		require.Condition(t, func() (success bool) {
			return order["first"].Compare(order["second"]) < 0
		})
	}
}

func TestNewEnforceOr(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    *identity.XRHID
		Expected error
	}
	testCases := []TestCase{
		{
			Name: "'Identity' type is not 'System'",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
				},
			},
			Expected: fmt.Errorf("'Identity.Type' must be 'System'"),
		},
		{
			Name: "Identity type is not 'User'",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "System",
				},
			},
			Expected: fmt.Errorf("'Identity.System.CertType' is not 'system'"),
		},
		{
			Name: "Success case for system",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "System",
					System: &identity.System{
						CertType:   "system",
						CommonName: "10fbb716-ca5d-11ed-b384-482ae3863d30",
					},
				},
			},
			Expected: nil,
		},
		{
			Name: "Success case for user",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
					User: &identity.User{
						Active:   true,
						UserID:   "jdoe",
						Username: "jdoe",
					},
				},
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		predicate := NewEnforceOr(EnforceSystemPredicate, EnforceUserPredicate)
		err := predicate(testCase.Given)
		if testCase.Expected != nil {
			require.Error(t, err)
			assert.EqualError(t, err, testCase.Expected.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
