package router

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	mock_rbac "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/rbac/impl"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	client_inventory "github.com/podengo-project/idmsvc-backend/internal/test/mock/interface/client/inventory"
	client_rbac "github.com/podengo-project/idmsvc-backend/internal/usecase/client/rbac"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	handler_impl "github.com/podengo-project/idmsvc-backend/internal/handler/impl"
)

func TestGetMajorVersion(t *testing.T) {
	assert.Equal(t, "", getMajorVersion(""))
	assert.Equal(t, "1", getMajorVersion("1.0"))
	assert.Equal(t, "1", getMajorVersion("1.0.3"))
	assert.Equal(t, "1", getMajorVersion("1."))
	assert.Equal(t, "a", getMajorVersion("a.b.c"))
}

func TestCheckRouterConfig(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    RouterConfig
		Expected error
	}
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg)
	testCases := []TestCase{
		{
			Name: "PublicPath is empty",
			Given: RouterConfig{
				PublicPath: "",
			},
			Expected: fmt.Errorf("PublicPath cannot be empty"),
		},
		{
			Name: "PrivatePath is empty",
			Given: RouterConfig{
				PublicPath:  "/api/idmsvc/v1",
				PrivatePath: "",
			},
			Expected: fmt.Errorf("PrivatePath cannot be empty"),
		},
		{
			Name: "PublicPath and PrivatePath are equal",
			Given: RouterConfig{
				PublicPath:  "/api/idmsvc/v1",
				PrivatePath: "/api/idmsvc/v1",
				Version:     "",
			},
			Expected: fmt.Errorf("PublicPath and PrivatePath cannot be equal"),
		},
		{
			Name: "Version is empty",
			Given: RouterConfig{
				PublicPath:  "/api/idmsvc/v1",
				PrivatePath: "/private",
				Version:     "",
			},
			Expected: fmt.Errorf("Version cannot be empty"),
		},
		{
			Name: "Metrics is nil",
			Given: RouterConfig{
				PublicPath:  "/api/idmsvc/v1",
				PrivatePath: "/private",
				Version:     "1.0",
			},
			Expected: fmt.Errorf("Metrics is nil"),
		},
		{
			Name: "Success scenario",
			Given: RouterConfig{
				PublicPath:  "/api/idmsvc/v1",
				PrivatePath: "/private",
				Version:     "1.0",
				Metrics:     metrics,
			},
			Expected: nil,
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		err := checkRouterConfig(testCase.Given)
		if testCase.Expected != nil {
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
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
	rcfg := RouterConfig{
		MetricsPath: "/metrics",
		PrivatePath: "",
	}
	configCommonMiddlewares(e, rcfg)
}

func TestNewRouterWithConfig(t *testing.T) {
	assert.Panics(t, func() {
		NewRouterWithConfig(nil, RouterConfig{})
	}, "'e' is nil")

	e := echo.New()
	assert.Panics(t, func() {
		NewRouterWithConfig(e, RouterConfig{})
	})

	cfg := test.GetTestConfig()
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg)
	_, db, _ := test.NewSqlMock(&gorm.Session{SkipHooks: true})
	inventory := client_inventory.NewHostInventory(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svcRbac, mockRbac := mock_rbac.NewRbacMock(ctx, cfg)
	err := svcRbac.Start()
	defer svcRbac.Stop()
	require.NoError(t, err)
	require.NoError(t, mockRbac.WaitAddress(3*time.Second))
	mockRbac.SetPermissions(mock_rbac.Profiles["domain-admin-profile"])
	rbacClient, err := client_rbac.NewClient("idmsvc", client_rbac.WithBaseURL(cfg.Clients.RbacBaseURL))
	if err != nil {
		panic(err)
	}
	rbac := client_rbac.New(cfg.Clients.RbacBaseURL, rbacClient)
	// Create application handlers
	app := handler_impl.NewHandler(cfg, db, metrics, inventory, rbac)

	goodConfig := RouterConfig{
		Version:     "1.0",
		PublicPath:  "/api/idmsvc",
		PrivatePath: "/private",
		Handlers:    app,
		Metrics:     metrics,
	}
	badConfig := RouterConfig{}

	assert.Panics(t, func() {
		_ = NewRouterWithConfig(e, badConfig)
	})

	assert.NotPanics(t, func() {
		e = NewRouterWithConfig(e, goodConfig)
	})
}

func TestNewRouterForMetrics(t *testing.T) {
	assert.Panics(t, func() {
		NewRouterForMetrics(nil, RouterConfig{})
	})

	e := echo.New()
	assert.Panics(t, func() {
		NewRouterForMetrics(e, RouterConfig{})
	})

	cfg := test.GetTestConfig()
	reg := prometheus.NewRegistry()
	metrics := metrics.NewMetrics(reg)
	_, db, _ := test.NewSqlMock(&gorm.Session{SkipHooks: true})
	inventory := client_inventory.NewHostInventory(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svcRbac, mockRbac := mock_rbac.NewRbacMock(ctx, cfg)
	err := svcRbac.Start()
	defer svcRbac.Stop()
	require.NoError(t, err)
	mockRbac.WaitAddress(3 * time.Second)
	mockRbac.SetPermissions(mock_rbac.Profiles["domain-admin-profile"])
	rbacClient, err := client_rbac.NewClient("idmsvc", client_rbac.WithBaseURL(cfg.Clients.RbacBaseURL))
	if err != nil {
		panic(err)
	}
	rbac := client_rbac.New(cfg.Clients.RbacBaseURL, rbacClient)
	// Create application handlers
	app := handler_impl.NewHandler(cfg, db, metrics, inventory, rbac)

	goodConfig := RouterConfig{
		Version:     "1.0",
		PublicPath:  "/api/idmsvc",
		PrivatePath: "/private",
		MetricsPath: "/metrics",
		Handlers:    app,
		Metrics:     metrics,
	}
	badConfig := RouterConfig{
		Version:     "1.0",
		PublicPath:  "/api/idmsvc",
		PrivatePath: "/private",
		Handlers:    app,
		Metrics:     metrics,
	}

	assert.Panics(t, func() {
		_ = NewRouterForMetrics(e, badConfig)
	}, "'e' is nil")

	assert.NotPanics(t, func() {
		e = NewRouterForMetrics(e, goodConfig)
	}, "MetricsPath cannot be an empty string")
}
