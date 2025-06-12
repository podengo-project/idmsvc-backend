package middleware

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
)

const testPath = "/test"

func helperCreatePredicate(principal string) IdentityPredicate {
	return func(data *identity.XRHID) error {
		if data == nil {
			return internal_errors.NilArgError("data")
		}
		if data.Identity.Type == "User" {
			if data.Identity.User == nil {
				return fmt.Errorf("'Itentity.User' is nil")
			}
			if data.Identity.User.Username == principal {
				return fmt.Errorf("principal='%s' is not accepted", principal)
			}
		}
		if data.Identity.Type == "System" {
			if data.Identity.System == nil {
				return fmt.Errorf("'Identity.System' is nil")
			}
			if data.Identity.System.CommonName == principal {
				return fmt.Errorf("principal='%s' is not accepted", principal)
			}
		}
		return nil
	}
}

func helperNewEchoEnforceIdentity(m echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	}
	e.Use(ContextLogConfig(&LogConfig{}))
	e.Use(CreateContext())
	e.Use(ParseXRHIDMiddlewareWithConfig(&ParseXRHIDMiddlewareConfig{}))
	e.Use(m)
	e.Add("GET", testPath, h)

	return e
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

	badUserXRHID := api.NewUserXRHID().
		WithOrgID("12345").
		WithUsername("test-fail-predicate").
		Build()
	userXRHID := api.NewUserXRHID().
		WithOrgID("12345").
		WithUsername("testuser").
		Build()
	systemXRHID := api.NewSystemXRHID().
		WithOrgID("12345").
		WithCommonName("testuser").
		Build()

	testCases := []TestCase{
		{
			Name:  header.HeaderXRHID + " header not present",
			Given: nil,
			Expected: TestCaseExpected{
				Code: http.StatusUnauthorized,
				Body: "{\"message\":\"Unauthorized\"}\n",
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
					&badUserXRHID,
				),
			),
			Expected: TestCaseExpected{
				Code: http.StatusUnauthorized,
				Body: "{\"message\":\"Unauthorized\"}\n",
			},
		},
		{
			Name: header.HeaderXRHID + " user pass predicates",
			Given: pointy.String(
				header.EncodeXRHID(
					&userXRHID,
				),
			),
			Expected: TestCaseExpected{
				Code: http.StatusOK,
				Body: "Ok",
			},
		},
		{
			Name: header.HeaderXRHID + " system pass predicates",
			Given: pointy.String(
				header.EncodeXRHID(
					&systemXRHID,
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
		ContextLogConfig(&LogConfig{}),
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

	userXRHID := api.NewUserXRHID().
		WithOrgID("12345").
		WithUsername("test-fail-predicate").
		Build()

	// When skipper return false, return 403 Unauthorized fail due to failed predicate
	e = helperNewEchoEnforceIdentity(
		EnforceIdentityWithConfig(
			&IdentityConfig{
				Skipper: helperSkipper(false),
				Predicates: []IdentityPredicateEntry{
					{
						Name:      "test-fail-predicate",
						Predicate: helperCreatePredicate("test-fail-predicate"),
					},
				},
			},
		),
	)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(header.HeaderXRHID, header.EncodeXRHID(&userXRHID))
	e.ServeHTTP(res, req)
	// Check expectations
	data, err = io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.Code)
	assert.Equal(t, "{\"message\":\"Unauthorized\"}\n", string(data))

	// When skipper return true, the predicate is not run and request is authorised
	e = helperNewEchoEnforceIdentity(
		EnforceIdentityWithConfig(
			&IdentityConfig{
				Skipper: helperSkipper(true),
				Predicates: []IdentityPredicateEntry{
					{
						Name:      "test-fail-predicate",
						Predicate: helperCreatePredicate("test-fail-predicate"),
					},
				},
			},
		),
	)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add(header.HeaderXRHID, header.EncodeXRHID(&userXRHID))
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
			Name: "Identity with User nil",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
					User: nil,
				},
			},
			Expected: fmt.Errorf("'Identity.User' is nil"),
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
			Name: "'Identity.AuthType' is not 'cert-auth'",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "System",
				},
			},
			Expected: fmt.Errorf("'Identity.AuthType' is not 'cert-auth'"),
		},
		{
			Name: "'Identity.System' is not nil",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					AuthType: "cert-auth",
					Type:     "System",
					System:   nil,
				},
			},
			Expected: fmt.Errorf("'Identity.System' is nil"),
		},
		{
			Name: "'CertType' is not 'system'",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					AuthType: "cert-auth",
					Type:     "System",
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
					AuthType: "cert-auth",
					Type:     "System",
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
					AuthType: "cert-auth",
					Type:     "System",
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
	xrhid := &identity.XRHID{}
	*xrhid = api.NewUserXRHID().
		WithAccountNumber("12345").
		WithOrgID("12345").
		WithUsername("test").
		WithUserID("12345").
		WithActive(true).
		WithInternal(true).
		WithOrgAdmin(true).
		WithLocale("en").
		Build()
	xrhidRaw := header.EncodeXRHID(xrhid)
	for i := 0; i < 1000; i++ {
		order["first"] = time.Time{}
		order["second"] = time.Time{}
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Add(header.HeaderXRHID, xrhidRaw)
		e.ServeHTTP(res, req)

		// Check expectations
		require.Condition(t, func() (success bool) {
			return order["first"].Compare(order["second"]) <= 0
		})
	}
}

// identityAlwaysTrue is a predicate that always return nil
// as if everything was ok.
func identityAlwaysTrue(*identity.XRHID) error {
	return nil
}

// newIdentityAlwaysFalse returns a predicate that always return
// a specified error.
// Return an IdentityPredicate.
func newIdentityAlwaysFalse(err error) IdentityPredicate {
	return func(*identity.XRHID) error {
		return err
	}
}

func TestNewEnforceOr(t *testing.T) {
	type TestCaseGiven struct {
		First  IdentityPredicate
		Second IdentityPredicate
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected error
	}
	testCases := []TestCase{
		{
			Name: "Fail first predicate, but return nil because second succeed",
			Given: TestCaseGiven{
				First:  newIdentityAlwaysFalse(errors.New("Failing first predicate")),
				Second: identityAlwaysTrue,
			},
			Expected: nil,
		},
		{
			Name: "Fail second predicate, but return nil because first succeed",
			Given: TestCaseGiven{
				First:  identityAlwaysTrue,
				Second: newIdentityAlwaysFalse(errors.New("Failing second predicate")),
			},
			Expected: nil,
		},
		{
			Name: "Failing both but return first",
			Given: TestCaseGiven{
				First:  newIdentityAlwaysFalse(errors.New("Failing both: first predicate")),
				Second: newIdentityAlwaysFalse(errors.New("Failing both: second predicate")),
			},
			Expected: fmt.Errorf("Failing both: first predicate\nFailing both: second predicate"),
		},
		{
			Name: "Success",
			Given: TestCaseGiven{
				First:  identityAlwaysTrue,
				Second: identityAlwaysTrue,
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		predicate := NewEnforceOr(testCase.Given.First, testCase.Given.Second)
		err := predicate(nil)
		if testCase.Expected != nil {
			require.Error(t, err)
			assert.EqualError(t, err, testCase.Expected.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestEnforceServiceAccountPredicate(t *testing.T) {
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
			Name: "'Identity.Type' is not 'ServiceAccount'",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "System",
				},
			},
			Expected: fmt.Errorf("'Identity.Type' must be 'ServiceAccount'"),
		},
		{
			Name: "'Identity.AuthType' must be 'jwt-auth'",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type:     "ServiceAccount",
					AuthType: "cert",
				},
			},
			Expected: fmt.Errorf("'Identity.AuthType' is not 'jwt-auth'"),
		},
		{
			Name: "'Identity.ServiceAccount' is nil",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type:           "ServiceAccount",
					AuthType:       "jwt-auth",
					ServiceAccount: nil,
				},
			},
			Expected: fmt.Errorf("'Identity.ServiceAccount' is nil"),
		},
		{
			Name: "'Identity.ServiceAccount.ClientId' is empty",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type:     "ServiceAccount",
					AuthType: "jwt-auth",
					ServiceAccount: &identity.ServiceAccount{
						ClientId: "",
					},
				},
			},
			Expected: fmt.Errorf("'Identity.ServiceAccount.ClientId' is empty"),
		},
		{
			Name: "'Identity.ServiceAccount.Username' is empty",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type:     "ServiceAccount",
					AuthType: "jwt-auth",
					ServiceAccount: &identity.ServiceAccount{
						ClientId: uuid.NewString(),
						Username: "",
					},
				},
			},
			Expected: fmt.Errorf("'Identity.ServiceAccount.Username' is empty"),
		},
		{
			Name: "Success ServiceAccount predicate",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type:     "ServiceAccount",
					AuthType: "jwt-auth",
					ServiceAccount: &identity.ServiceAccount{
						ClientId: uuid.NewString(),
						Username: "test-user",
					},
				},
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		err := EnforceServiceAccountPredicate(testCase.Given)
		if testCase.Expected != nil {
			assert.EqualError(t, err, testCase.Expected.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
