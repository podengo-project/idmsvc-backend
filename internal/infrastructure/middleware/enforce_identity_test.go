package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	b64 "encoding/base64"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/openlyinc/pointy"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"encoding/json"
)

const testPath = "/test"

func helperCreatePredicate(username string) Predicate {
	return func(data *identity.Identity) error {
		if data == nil {
			return fmt.Errorf("data is nil")
		}
		if data.User.Username == username {
			return fmt.Errorf("username='%s' is not accepted", username)
		}
		return nil
	}
}

func helperNewEchoEnforceIdentity(middleware echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	}
	e.Use(middleware)
	e.Add("GET", testPath, h)

	return e
}

func helperEncodeIdentity(id identity.Identity) string {
	if bytes, err := json.Marshal(id); err == nil {
		sEnc := b64.StdEncoding.EncodeToString(bytes)
		return sEnc
	}
	return ""
}

// FIXME
func helperGenerateUserIdentity(orgId string, username string) identity.Identity {
	return identity.Identity{
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
	}
}

func helperGenerateCertificateIdentity(orgId string, subjectDN string, issuerDN string) identity.Identity {
	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/identities/jwt.json
	return identity.Identity{
		AccountNumber: "11111",
		OrgID:         orgId,
		Internal: identity.Internal{
			OrgID: orgId,
		},
		Type:     "Certificate",
		AuthType: "basic-auth",
		X509: identity.X509{
			SubjectDN: subjectDN,
			IssuerDN:  issuerDN,
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
				helperEncodeIdentity(
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
				helperEncodeIdentity(
					helperGenerateUserIdentity("12345", "testuser"),
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
			NewIdentityConfig(nil).
				Add(
					"test-predicate",
					helperCreatePredicate("test-fail-predicate"),
				),
		),
	)
	for _, testCase := range testCases {
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
			NewIdentityConfig(helperSkipper(false)),
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
			NewIdentityConfig(helperSkipper(true)),
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
