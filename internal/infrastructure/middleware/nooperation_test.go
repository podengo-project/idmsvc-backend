package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func helperNewEchoNooperation(method string, path string, m echo.MiddlewareFunc) *echo.Echo {
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

func TestNooperation(t *testing.T) {
	testMethod := http.MethodGet
	testPath := "/test"
	m := Nooperation()
	e := helperNewEchoNooperation(testMethod, testPath, m)

	assert.NotPanics(t, func() {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, testPath, nil)
		e.ServeHTTP(res, req)

		// Check expectations
		data, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "Ok", string(data))
		assert.Equal(t, http.StatusOK, res.Code)
	})
}
