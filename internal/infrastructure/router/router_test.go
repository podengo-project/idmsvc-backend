package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	api_metrics "github.com/podengo-project/idmsvc-backend/internal/test/mock/api/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/test/mock/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestGetMajorVersion(t *testing.T) {
	assert.Equal(t, "", getMajorVersion(""))
	assert.Equal(t, "1", getMajorVersion("1.0"))
	assert.Equal(t, "1", getMajorVersion("1.0.3"))
	assert.Equal(t, "1", getMajorVersion("1."))
	assert.Equal(t, "a", getMajorVersion("a.b.c"))
}

func TestLoggerSkipperWithPaths(t *testing.T) {
	var skipper middleware.Skipper

	// Empty path does not panic
	assert.NotPanics(t, func() {
		skipper = loggerSkipperWithPaths()
	})
	assert.NotNil(t, skipper)

	// Only one path does not panic
	assert.NotPanics(t, func() {
		skipper = loggerSkipperWithPaths("/test")
	})
	assert.NotNil(t, skipper)

	// Check several paths
	assert.NotPanics(t, func() {
		skipper = loggerSkipperWithPaths("/test", "/anothertest")
	})
	assert.NotNil(t, skipper)

	// Check skipped paths
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/test")
	assert.True(t, skipper(ctx))

	req = httptest.NewRequest(http.MethodGet, "/anothertest", nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetPath("/anothertest")
	assert.True(t, skipper(ctx))

	// Check no skipped paths
	req = httptest.NewRequest(http.MethodGet, "/noskipped", nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetPath("/noskipped")
	assert.False(t, skipper(ctx))
}

func TestConfigCommonMiddlewares(t *testing.T) {
	e := echo.New()
	cfg := config.Get()
	configCommonMiddlewares(e, cfg)
}

func TestGuardNewRouterWithConfig(t *testing.T) {
	assert.PanicsWithValue(t, "'e' is nil", func() {
		guardNewRouterWithConfig(nil, nil, nil, nil)
	})

	e := echo.New()
	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		guardNewRouterWithConfig(e, nil, nil, nil)
	})

	cfg := test.GetTestConfig()
	assert.PanicsWithValue(t, "'app' is nil", func() {
		guardNewRouterWithConfig(e, cfg, nil, nil)
	})

	app := handler.NewApplication(t)
	assert.PanicsWithValue(t, "'metrics' is nil", func() {
		guardNewRouterWithConfig(e, cfg, app, nil)
	})
}

func TestNewRouterWithConfig(t *testing.T) {
	e := echo.New()
	cfg := test.GetTestConfig()
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg)
	app := handler.NewApplication(t)

	assert.NotPanics(t, func() {
		e = NewRouterWithConfig(e, cfg, app, metrics)
	})
	app.AssertExpectations(t)
}

func TestGuardNewRouterForMetrics(t *testing.T) {
	assert.PanicsWithValue(t, "'e' is nil", func() {
		guardNewRouterForMetrics(nil, nil, nil)
	})

	e := echo.New()
	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		guardNewRouterForMetrics(e, nil, nil)
	})

	cfg := test.GetTestConfig()
	assert.PanicsWithValue(t, "'apiMetrics' is nil", func() {
		guardNewRouterForMetrics(e, cfg, nil)
	})

	oldMetricsPath := cfg.Metrics.Path
	cfg.Metrics.Path = ""
	assert.PanicsWithValue(t, "'MetricsPath' cannot be an empty string", func() {
		guardNewRouterForMetrics(e, cfg, nil)
	})

	cfg.Metrics.Path = oldMetricsPath
	apiMetrics := api_metrics.NewServerInterface(t)
	assert.NotPanics(t, func() {
		guardNewRouterForMetrics(e, cfg, apiMetrics)
	})
	apiMetrics.AssertExpectations(t)
}

func TestNewRouterForMetrics(t *testing.T) {
	e := echo.New()
	cfg := test.GetTestConfig()
	app := handler.NewApplication(t)
	assert.NotPanics(t, func() {
		e = NewRouterForMetrics(e, cfg, app)
	}, "NewRouterForMetrics success")
	app.AssertExpectations(t)
}

func TestTrimVersionFromPathPrefix(t *testing.T) {
	assert.Equal(t, "", trimVersionFromPathPrefix(""))
	assert.Equal(t, "/api/idmsvc", trimVersionFromPathPrefix("/api/idmsvc"))
	assert.Equal(t, "/api/idmsvc", trimVersionFromPathPrefix("/api/idmsvc/"))
	assert.Equal(t, "/api/idmsvc", trimVersionFromPathPrefix("/api/idmsvc/v1"))
	assert.Equal(t, "/api/idmsvc", trimVersionFromPathPrefix("/api/idmsvc/v1.0"))
	assert.Equal(t, "/api/idmsvc", trimVersionFromPathPrefix("/api/idmsvc/v1.0/"))
}
