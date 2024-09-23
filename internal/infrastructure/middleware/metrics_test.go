package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const URLPrefix = "/api/" + config.DefaultAppName

func TestCreateMetricsMiddleware(t *testing.T) {
	var (
		m          *metrics.Metrics
		middleware echo.MiddlewareFunc
	)
	m = metrics.NewMetrics(prometheus.NewRegistry())
	middleware = CreateMetricsMiddleware(m)

	assert.NotNil(t, middleware)
}

func TestMetricsMiddlewareWithConfigCreation(t *testing.T) {
	var (
		reg              *prometheus.Registry
		config           *MetricsConfig
		wasSkipped       bool
		wasSkipperCalled bool
	)

	config = &MetricsConfig{
		Metrics: nil,
		Skipper: nil,
	}
	assert.Panics(t, func() {
		MetricsMiddlewareWithConfig(config)
	})

	reg = prometheus.NewRegistry()
	config = &MetricsConfig{
		Metrics: metrics.NewMetrics(reg),
		Skipper: func(c echo.Context) bool {
			wasSkipperCalled = true
			// We don't use c.Path() because when we
			// invoke ServeHTTP the Path attribute is
			// empty, however the c.Request().RequestURI
			// contain the request path
			wasSkipped = c.Request().RequestURI == "/ping"
			return wasSkipped
		},
	}

	require.NotPanics(t, func() {
		MetricsMiddlewareWithConfig(config)
	})

	assert.NotPanics(t, func() {
		MetricsMiddlewareWithConfig(nil)
	})

	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	}

	e := echo.New()
	e.Use(ContextLogConfig(&LogConfig{}))
	m := MetricsMiddlewareWithConfig(config)
	e.Use(m)
	path := "/api/idmsvc/v1/domains"
	e.Add(http.MethodGet, path, h)

	wasSkipperCalled = false
	wasSkipped = false
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	e.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "Ok", resp.Body.String())
	assert.True(t, wasSkipperCalled)
	assert.False(t, wasSkipped)

	// Check skipper
	wasSkipperCalled = false
	wasSkipped = false
	path = "/ping"
	resp = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, path, nil)
	e.ServeHTTP(resp, req)
	assert.True(t, wasSkipperCalled)
	assert.True(t, wasSkipped)
}
