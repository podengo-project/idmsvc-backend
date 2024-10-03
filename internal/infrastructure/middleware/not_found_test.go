package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func helperNewEchoNotFound(m echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	}
	e.Use(
		middleware.RemoveTrailingSlash(),
		ContextLogConfig(&LogConfig{}),
		CreateContext(),
		m,
	)

	e.Add(http.MethodPost, "/test-1", h)
	e.Add(http.MethodPost, "/test-2/:param1", h)
	e.Add(http.MethodPost, "/test-3/:param1/:param2", h)

	g := e.Group("/group-1")
	g.Add(http.MethodPost, "/test-1", h)
	g.Add(http.MethodPost, "/test-2/:param1", h)
	g.Add(http.MethodPost, "/test-3/:param1/:param2", h)

	g = e.Group("/group-2", Nooperation())
	g.Add(http.MethodPost, "/test-1", h)
	g.Add(http.MethodPost, "/test-2/:param1", h)
	g.Add(http.MethodPost, "/test-3/:param1/:param2", h)

	return e
}

func TestNotFoundMiddleware(t *testing.T) {
	type TestCase struct {
		Name               string
		Path               string
		ExpectedStatusCode int
	}

	testCases := []TestCase{
		// No group
		{
			Name:               "Not found for /not-found",
			Path:               "/not-found",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "Not found for /",
			Path:               "/",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "Not found for /test-2",
			Path:               "/test-2",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "Not found for /test-3/value-1",
			Path:               "/test-3/value-1",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "OK for /test-1",
			Path:               "/test-1",
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "OK for /test-2/:param1",
			Path:               "/test-2/value-1",
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "OK for /test-3/:param1/:param2",
			Path:               "/test-3/value-1/value-2",
			ExpectedStatusCode: http.StatusOK,
		},
		// With group-1
		{
			Name:               "Not found for /group-1/not-found",
			Path:               "/group-1/not-found",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "Not found for /group-1",
			Path:               "/group-1",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "Not found for /group-1/test-2",
			Path:               "/group-1/test-2",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "Not found for /group-1/test-3/value-1",
			Path:               "/group-1/test-3/value-1",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "Not found for /group-1/test-3/value-1/",
			Path:               "/group-1/test-3/value-1/",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "OK for /group-1/test-1",
			Path:               "/group-1/test-1",
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "OK for /group-1/test-2/:param1",
			Path:               "/group-1/test-2/value-1",
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "OK for /group-1/test-3/:param1/:param2",
			Path:               "/group-1/test-3/value-1/value-2",
			ExpectedStatusCode: http.StatusOK,
		},
		// With group-2
		{
			// This replay Path == "/group-2/*"
			Name:               "Not found for /group-2/not-found",
			Path:               "/group-2/not-found",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "Not found for /group-2",
			Path:               "/group-2",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			// This replay Path ==  "/group-2/*"
			Name:               "Not found for /group-2/test-2",
			Path:               "/group-2/test-2",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			// This replay Path ==  "/group-2/*"
			Name:               "Not found for /group-2/test-3/value-1",
			Path:               "/group-2/test-3/value-1",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			// This replay Path ==  "/group-2/*"
			Name:               "Not found for /group-2/test-3/value-1/",
			Path:               "/group-2/test-3/value-1/",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:               "OK for /group-2/test-1",
			Path:               "/group-2/test-1",
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "OK for /group-2/test-2/:param1",
			Path:               "/group-2/test-2/value-1",
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "OK for /group-2/test-3/:param1/:param2",
			Path:               "/group-2/test-3/value-1/value-2",
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)

		// Given
		e := helperNewEchoNotFound(NotFound())

		// When
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, testCase.Path, nil)
		e.ServeHTTP(res, req)

		// Then
		assert.Equal(t, testCase.ExpectedStatusCode, res.Result().StatusCode)
	}
}
