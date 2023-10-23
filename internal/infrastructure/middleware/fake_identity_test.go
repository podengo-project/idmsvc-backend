package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func helperNewEchoFakeIdentity(method string, path string, m echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	}
	e.Use(m)
	switch method {
	case http.MethodConnect:
	case http.MethodGet:
	case http.MethodHead:
	case http.MethodOptions:
	case http.MethodDelete:
	case http.MethodPatch:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodTrace:
	default:
		panic("'method' is invalid")
	}
	e.Add(method, path, h)

	return e
}

func assertHttpHeader(t *testing.T, expected http.Header, current http.Header) {
	if (expected == nil || len(expected) == 0) && (current == nil || len(current) == 0) {
		return
	}
	require.NotNil(t, expected)
	require.NotNil(t, current)
	assert.Equal(t, expected, current)
	// for ek, ev := range expected {
	// 	cv := current.Get(ek)
	// 	assert.Equal(t, ev, cv)
	// }
	// for k, cv := range current {
	// 	ev := expected.Get(k)
	// 	assert.Equal(t, cv, ev)
	// }
}

func TestFakeIdentityWithConfigPanic(t *testing.T) {
	assert.Panics(t, func() {
		FakeIdentityWithConfig(nil)
	})
}

func TestFakeIdentityWithConfig(t *testing.T) {
	const (
		testXRHIDString = "eyJpZGVudGl0eSI6eyJhY2NvdW50X251bWJlciI6IjIyNzUxIiwiYXV0aF90eXBlIjoiY2VydC1hdXRoIiwiaW50ZXJuYWwiOnsib3JnX2lkIjoiMTIzNDUifSwib3JnX2lkIjoiMTIzNDUiLCJ0eXBlIjoiVXNlciIsInVzZXIiOnsiZW1haWwiOiJhbm51YWxAY2hlbW9zaC5pbyIsImZpcnN0X25hbWUiOiJTYW5keSIsImlzX2FjdGl2ZSI6dHJ1ZSwiaXNfaW50ZXJuYWwiOnRydWUsImlzX29yZ19hZG1pbiI6dHJ1ZSwibGFzdF9uYW1lIjoiTGVkbmVyIiwibG9jYWxlIjoiYnMiLCJ1c2VyX2lkIjoidGVzdCIsInVzZXJuYW1lIjoidGVzdCJ9fX0K"
	)
	type TestCaseExpected struct {
		Header http.Header
	}
	type TestCase struct {
		Name     string
		Given    http.Header
		Expected http.Header
	}

	testCases := []TestCase{
		{
			Name:     "no fake header present",
			Given:    http.Header{},
			Expected: http.Header{},
		},
		{
			Name: "X-Rh-Fake-Identity only",
			Given: http.Header(map[string][]string{
				headerXRhFakeIdentity: []string{testXRHIDString},
			}),
			Expected: http.Header(map[string][]string{
				headerXRhIdentity: []string{testXRHIDString},
			}),
		},
	}

	// Get echo instance with the middleware and one predicate for test it
	e := helperNewEchoFakeIdentity(
		http.MethodGet,
		testPath,
		FakeIdentityWithConfig(
			&FakeIdentityConfig{},
		))
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		if testCase.Given != nil {
			for k, v := range testCase.Given {
				req.Header.Set(k, strings.Join(v, "; "))
			}
		}
		e.ServeHTTP(res, req)

		// Check expectations
		data, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "Ok", string(data))
		assert.Equal(t, http.StatusOK, res.Code)
		assertHttpHeader(t, testCase.Expected, req.Header)
	}
}

func TestFakeIdentitySkipper(t *testing.T) {
	const (
		testXRHIDString = "eyJpZGVudGl0eSI6eyJhY2NvdW50X251bWJlciI6IjIyNzUxIiwiYXV0aF90eXBlIjoiY2VydC1hdXRoIiwiaW50ZXJuYWwiOnsib3JnX2lkIjoiMTIzNDUifSwib3JnX2lkIjoiMTIzNDUiLCJ0eXBlIjoiVXNlciIsInVzZXIiOnsiZW1haWwiOiJhbm51YWxAY2hlbW9zaC5pbyIsImZpcnN0X25hbWUiOiJTYW5keSIsImlzX2FjdGl2ZSI6dHJ1ZSwiaXNfaW50ZXJuYWwiOnRydWUsImlzX29yZ19hZG1pbiI6dHJ1ZSwibGFzdF9uYW1lIjoiTGVkbmVyIiwibG9jYWxlIjoiYnMiLCJ1c2VyX2lkIjoidGVzdCIsInVzZXJuYW1lIjoidGVzdCJ9fX0K"
	)
	var (
		e    *echo.Echo
		res  *httptest.ResponseRecorder
		req  *http.Request
		data []byte
		err  error
	)

	// With no skipper, the middleware execute always
	e = helperNewEchoFakeIdentity(
		http.MethodGet,
		testPath,
		FakeIdentityWithConfig(
			&FakeIdentityConfig{},
		),
	)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header[headerXRhFakeIdentity] = []string{testXRHIDString}
	e.ServeHTTP(res, req)
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "Ok", string(data))
	assert.Equal(t, 1, len(req.Header))
	assert.Equal(t, testXRHIDString, req.Header.Get(headerXRhIdentity))
	assert.Equal(t, "", req.Header.Get(headerXRhFakeIdentity))
	assert.Equal(t, testXRHIDString, req.Header.Get(headerXRhIdentity))

	// With skipper returning false, the middleware excute always
	e = helperNewEchoFakeIdentity(
		http.MethodGet,
		testPath,
		FakeIdentityWithConfig(
			&FakeIdentityConfig{
				Skipper: helperSkipper(false),
			},
		),
	)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header[headerXRhFakeIdentity] = []string{testXRHIDString}
	e.ServeHTTP(res, req)
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "Ok", string(data))
	assert.Equal(t, 1, len(req.Header))
	assert.Equal(t, testXRHIDString, req.Header.Get(headerXRhIdentity))
	assert.Equal(t, "", req.Header.Get(headerXRhFakeIdentity))
	assert.Equal(t, testXRHIDString, req.Header.Get(headerXRhIdentity))

	// With skipper returning true, the middleware does not execute
	e = helperNewEchoFakeIdentity(
		http.MethodGet,
		testPath,
		FakeIdentityWithConfig(
			&FakeIdentityConfig{
				Skipper: helperSkipper(true),
			},
		),
	)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header[headerXRhFakeIdentity] = []string{testXRHIDString}
	e.ServeHTTP(res, req)
	data, err = ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "Ok", string(data))
	assert.Equal(t, 1, len(req.Header))
	assert.Equal(t, testXRHIDString, req.Header.Get(headerXRhFakeIdentity))
	assert.Equal(t, "", req.Header.Get(headerXRhIdentity))
}
