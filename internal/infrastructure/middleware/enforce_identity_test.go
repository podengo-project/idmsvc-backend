package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hmsidm/internal/api/header"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/openlyinc/pointy"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPath = "/test"

func helperCreatePredicate(username string) IdentityPredicate {
	return func(data *identity.XRHID) error {
		if data == nil {
			return fmt.Errorf("'data' is nil")
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
			User: identity.User{
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
			System: identity.System{
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
			Name:  "x-rh-identity header not present",
			Given: nil,
			Expected: TestCaseExpected{
				Code: http.StatusUnauthorized,
				Body: "{\"message\":\"Unauthorized\"}\n",
			},
		},
		{
			Name:  "x-rh-identity bad base64 coding",
			Given: pointy.String("bad base64 coding"),
			Expected: TestCaseExpected{
				Code: http.StatusUnauthorized,
				Body: "{\"message\":\"Unauthorized\"}\n",
			},
		},
		{
			Name:  "x-rh-identity bad json encoding",
			Given: pointy.String("ewo="),
			Expected: TestCaseExpected{
				Code: http.StatusUnauthorized,
				Body: "{\"message\":\"Unauthorized\"}\n",
			},
		},
		{
			Name: "x-rh-identity fail predicates",
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
			Name: "x-rh-identity pass predicates",
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
			Name: "x-rh-identity pass predicates",
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
			NewIdentityConfig().
				SetSkipper(nil).
				AddPredicate(
					"test-predicate",
					helperCreatePredicate("test-fail-predicate"),
				),
		),
	)
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		if testCase.Given != nil {
			req.Header.Add("X-Rh-Identity", *testCase.Given)
		}
		e.ServeHTTP(res, req)

		// Check expectations
		data, err := ioutil.ReadAll(res.Body)
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
			NewIdentityConfig().
				SetSkipper(nil).
				AddPredicate(
					"test-predicate",
					helperCreatePredicate("test-fail-predicate"),
				),
		),
	)
	e.Add("GET", testPath, h)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	iden := header.EncodeXRHID(&identity.XRHID{})
	req.Header.Add("X-Rh-Identity", iden)
	e.ServeHTTP(res, req)

	// Check expectations
	data, err := ioutil.ReadAll(res.Body)
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
			NewIdentityConfig().
				SetSkipper(helperSkipper(false)),
		),
	)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	e.ServeHTTP(res, req)
	// Check expectations
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.Code)
	assert.Equal(t, "{\"message\":\"Unauthorized\"}\n", string(data))

	// When skipper return true the middleware does not process the header or the predicates
	e = helperNewEchoEnforceIdentity(
		EnforceIdentityWithConfig(
			NewIdentityConfig().
				SetSkipper(helperSkipper(true)),
		),
	)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	e.ServeHTTP(res, req)
	// Check expectations
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "Ok", string(data))
}

func helperNewContextForSkipper(route string, method string, path string, headers map[string]string) echo.Context {
	// See: https://echo.labstack.com/guide/testing/
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(route)
	return c
}

func TestSkipperUSer(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    string
		Expected bool
	}
	testCases := []TestCase{}
	for i := range userEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("No skip userEnforceRoutes[%d]", i),
			Given:    userEnforceRoutes[i],
			Expected: false,
		})
	}
	for i := range systemEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("Skip systemEnforceRoutes[%d]", i),
			Given:    systemEnforceRoutes[i],
			Expected: true,
		})
	}
	for _, testCase := range testCases {
		ctx := helperNewContextForSkipper(testCase.Given, http.MethodGet, testCase.Given, nil)
		result := SkipperUserPredicate(ctx)
		assert.Equal(t, testCase.Expected, result)
	}
}

func TestSkipperSystem(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    string
		Expected bool
	}
	testCases := []TestCase{}
	for i := range systemEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("No skip systemEnforceRoutes[%d]", i),
			Given:    systemEnforceRoutes[i],
			Expected: false,
		})
	}
	for i := range userEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("Skip userEnforceRoutes[%d]", i),
			Given:    userEnforceRoutes[i],
			Expected: true,
		})
	}
	for _, testCase := range testCases {
		ctx := helperNewContextForSkipper(testCase.Given, http.MethodGet, testCase.Given, nil)
		result := SkipperSystemPredicate(ctx)
		assert.Equal(t, testCase.Expected, result)
	}
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
			Expected: fmt.Errorf("'Identity.Type' is not 'User'"),
		},
		{
			Name: "Identity with disabled user",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
					User: identity.User{
						Active: false,
					},
				},
			},
			Expected: fmt.Errorf("'Identity.User.Active' is not true"),
		},
		{
			Name: "'UserID' is empty",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
					User: identity.User{
						Active: true,
						UserID: "",
					},
				},
			},
			Expected: fmt.Errorf("'Identity.User.UserID' cannot be empty"),
		},
		{
			Name: "'UserName' is empty",
			Given: &identity.XRHID{
				Identity: identity.Identity{
					Type: "User",
					User: identity.User{
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
					User: identity.User{
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
					System: identity.System{
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
					System: identity.System{
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
					System: identity.System{
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
