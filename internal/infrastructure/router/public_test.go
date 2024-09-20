package router

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	mock_rbac "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/rbac/impl"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/rbac"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/podengo-project/idmsvc-backend/internal/test/mock/handler"
	client_rbac "github.com/podengo-project/idmsvc-backend/internal/usecase/client/rbac"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func helperNewContextForSkipper(route string, method string, path string, headers map[string]string) echo.Context {
	// See: https://echo.labstack.com/guide/testing/
	e := echo.New()
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(route)
	return c
}

func initRbacWrapper(t *testing.T, cfg *config.Config) rbac.Rbac {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svcRbac, mockRbac := mock_rbac.NewRbacMock(ctx, cfg)
	err := svcRbac.Start()
	require.NoError(t, err)
	defer svcRbac.Stop()
	mockRbac.WaitAddress(3 * time.Second)
	mockRbac.SetPermissions(mock_rbac.Profiles["domain-admin-profile"])
	rbacClient, err := client_rbac.NewClient("idmsvc", client_rbac.WithBaseURL(cfg.Clients.RbacBaseURL))
	require.NoError(t, err)
	assert.NotNil(t, rbacClient)
	rbac := client_rbac.New(cfg.Clients.RbacBaseURL, rbacClient)
	return rbac
}

func TestNewGroupPublicNotPanic(t *testing.T) {
	cfg := test.GetTestConfig()
	require.NotNil(t, cfg)

	reg := prometheus.NewRegistry()
	require.NotNil(t, reg)
	metrics := metrics.NewMetrics(reg)
	require.NotNil(t, metrics)
	appPrefix := trimVersionFromPathPrefix(cfg.Application.PathPrefix)
	version := "1.0"
	app := handler.NewApplication(t)
	require.NotNil(t, app)
	e := echo.New()
	require.NotNil(t, e)

	assert.NotPanics(t, func() {
		newGroupPublic(e.Group(appPrefix+"/"+version), cfg, app, metrics)
	})
	app.AssertExpectations(t)
}

func TestNewGroupPublic(t *testing.T) {
	const (
		appPrefix   = "/api"
		appName     = "/idmsvc"
		versionFull = "/v1.0"
	)
	type TestCaseExpected map[string]map[string]string
	empty := ""

	testCases := TestCaseExpected{
		"/private/readyz": {
			"GET": empty,
		},
		"/private/livez": {
			"GET": empty,
		},

		appPrefix + appName + versionFull + "/openapi.json": {
			"GET": empty,
		},

		appPrefix + appName + versionFull + "/domains": {
			"GET":  empty,
			"POST": empty,
		},

		appPrefix + appName + versionFull + "/domains/token": {
			"POST": empty,
		},

		appPrefix + appName + versionFull + "/domains/:uuid": {
			"GET":    empty,
			"PUT":    empty,
			"PATCH":  empty,
			"DELETE": empty,
		},

		appPrefix + appName + versionFull + "/host-conf/:inventory_id/:fqdn": {
			"POST": empty,
		},

		appPrefix + appName + versionFull + "/signing_keys": {
			"GET": empty,
		},

		// This routes are added when the group is created
		appPrefix + appName + versionFull + "/*": {
			"echo_route_not_found": empty,
		},
		appPrefix + appName + versionFull: {
			"echo_route_not_found": empty,
		},
	}

	reg := prometheus.NewRegistry()
	require.NotNil(t, reg)
	metrics := metrics.NewMetrics(reg)
	require.NotNil(t, metrics)
	cfg := test.GetTestConfig()
	require.NotNil(t, cfg)
	app := handler.NewApplication(t)
	require.NotNil(t, app)
	e := echo.New()
	require.NotNil(t, e)
	version := "1.0"
	pathPrefix := trimVersionFromPathPrefix(cfg.Application.PathPrefix)
	group := newGroupPublic(e.Group(pathPrefix+"/v"+version), cfg, app, metrics)
	require.NotNil(t, group)

	// Match Routes in expected
	for _, route := range e.Routes() {
		t.Logf("Method=%s Path=%s Name=%s", route.Method, route.Path, route.Name)

		methods, okPath := testCases[route.Path]
		require.Truef(t, okPath, "searching path=%s", route.Path)
		name, okMethod := methods[route.Method]
		require.Truef(t, okMethod, "searching method=%s for path=%s", route.Method, route.Path)
		if name != empty {
			assert.Equalf(t, name, route.Name, "handler for path=%s method=%s does not match: expected=%s current=%s", route.Path, route.Method, name, route.Name)
		}
	}
	app.AssertExpectations(t)

	// Same result when IsFakeEnabled
	e = echo.New()
	require.NotNil(t, e)
	group = newGroupPublic(
		e.Group(pathPrefix+"/v"+version),
		cfg,
		app,
		metrics)
	require.NotNil(t, group)
	for _, route := range e.Routes() {
		t.Logf("Method=%s Path=%s Name=%s", route.Method, route.Path, route.Name)

		methods, okPath := testCases[route.Path]
		assert.Truef(t, okPath, "path=%s not found into the expected ones", route.Path)

		_, okMethod := methods[route.Method]
		assert.Truef(t,
			okMethod,
			"method=%s not found into the expected ones for the path=%s",
			route.Method,
			route.Path)
	}
}

func TestSkipperUser(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    enforceRoute
		Expected bool
	}
	testCases := []TestCase{}
	for i := range userEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("No skip userEnforceRoutes[%d]: %v", i, userEnforceRoutes[i]),
			Given:    userEnforceRoutes[i],
			Expected: false,
		})
	}
	for i := range systemEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("Skip systemEnforceRoutes[%d]: %v", i, systemEnforceRoutes[i]),
			Given:    systemEnforceRoutes[i],
			Expected: true,
		})
	}
	for i := range mixedEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("Skip mixedEnforceRoutes[%d]: %v", i, mixedEnforceRoutes[i]),
			Given:    mixedEnforceRoutes[i],
			Expected: true,
		})
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		ctx := helperNewContextForSkipper(
			testCase.Given.Path,
			testCase.Given.Method,
			testCase.Given.Path,
			nil)
		result := skipperUserPredicate(ctx)
		assert.Equal(t, testCase.Expected, result)
	}
}

func TestSkipperSystem(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    enforceRoute
		Expected bool
	}
	testCases := []TestCase{}
	for i := range systemEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("No skip systemEnforceRoutes[%d]: %v", i, systemEnforceRoutes[i]),
			Given:    systemEnforceRoutes[i],
			Expected: false,
		})
	}
	for i := range userEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("Skip userEnforceRoutes[%d]: %v", i, userEnforceRoutes[i]),
			Given:    userEnforceRoutes[i],
			Expected: true,
		})
	}
	for i := range mixedEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("Skip mixedEnforceRoutes[%d]: %v", i, mixedEnforceRoutes[i]),
			Given:    mixedEnforceRoutes[i],
			Expected: true,
		})
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		ctx := helperNewContextForSkipper(
			testCase.Given.Path,
			testCase.Given.Method,
			testCase.Given.Path,
			nil,
		)
		result := skipperSystemPredicate(ctx)
		assert.Equal(t, testCase.Expected, result)
	}
}

func TestGetOpenapiPaths(t *testing.T) {

	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		getOpenapiPaths(nil, "")
	})

	cfg := test.GetTestConfig()
	assert.PanicsWithValue(t, "'version' is an empty string", func() {
		getOpenapiPaths(cfg, "")
	})

	cachedPaths := getOpenapiPaths(cfg, "1.4")
	assert.NotNil(t, cachedPaths)
	assert.Equal(t,
		[]string{
			"/api/idmsvc/v1.4/openapi.json",
			"/api/idmsvc/v1/openapi.json",
		},
		cachedPaths(),
	)
}

func TestNewSkipperOpenapi(t *testing.T) {
	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		newSkipperOpenapi(nil, "")
	})

	cfg := test.GetTestConfig()
	assert.PanicsWithValue(t, "'version' is an empty string", func() {
		newSkipperOpenapi(cfg, "")
	})

	version := "1.4"
	skipper := newSkipperOpenapi(cfg, version)
	assert.NotNil(t, skipper)

	prefixPath := trimVersionFromPathPrefix(cfg.Application.PathPrefix)
	path := fmt.Sprintf("%s/v%s/openapi.json", prefixPath, version)
	ctx := helperNewContextForSkipper(path, echo.GET, path, map[string]string{})
	assert.True(t, skipper(ctx))
	path = fmt.Sprintf("%s/v%s/openapi.json", prefixPath, strings.Split(version, ".")[0])
	ctx = helperNewContextForSkipper(path, echo.GET, path, map[string]string{})
	assert.True(t, skipper(ctx))
	path = fmt.Sprintf("%s/v%s/openapi2.json", prefixPath, strings.Split(version, ".")[0])
	ctx = helperNewContextForSkipper(path, echo.GET, path, map[string]string{})
	assert.False(t, skipper(ctx))
}

func TestNewRbacSkipper(t *testing.T) {
	var skipper echo_middleware.Skipper

	require.PanicsWithValue(t, "service is an empty string", func() {
		skipper = newRbacSkipper("")
	}, "newRbacSkipper panics on empty service")

	require.NotPanics(t, func() {
		skipper = newRbacSkipper("idmsvc")
	}, "no panics for idmsvc")
	skipper(echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/api/idmsvc/v1/openapi.json", http.NoBody), httptest.NewRecorder()))
}

func TestInitRbacMiddleware(t *testing.T) {
	var result echo.MiddlewareFunc
	assert.NotPanics(t, func() {
		result = initRbacMiddleware(&config.Config{
			Application: config.Application{
				EnableRBAC: false,
			},
		})
	}, "Return DefaultNooperation")
	assert.NotNil(t, result)

	assert.NotPanics(t, func() {
		result = initRbacMiddleware(&config.Config{
			Application: config.Application{
				EnableRBAC: true,
			},
			Clients: config.Clients{
				RbacBaseURL: "http://rbac:8000",
			},
		})
	}, "Initialize the rbac middleware")
	assert.NotNil(t, result)
}

func TestSkipperMixedPredicate(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    enforceRoute
		Expected bool
	}
	testCases := []TestCase{}
	for i := range userEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("Skip userEnforceRoutes[%d]: %v", i, userEnforceRoutes[i]),
			Given:    userEnforceRoutes[i],
			Expected: true,
		})
	}
	for i := range systemEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("Skip systemEnforceRoutes[%d]: %v", i, systemEnforceRoutes[i]),
			Given:    systemEnforceRoutes[i],
			Expected: true,
		})
	}
	for i := range mixedEnforceRoutes {
		testCases = append(testCases, TestCase{
			Name:     fmt.Sprintf("No skip mixedEnforceRoutes[%d]: %v", i, mixedEnforceRoutes[i]),
			Given:    mixedEnforceRoutes[i],
			Expected: false,
		})
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		ctx := helperNewContextForSkipper(
			testCase.Given.Path,
			testCase.Given.Method,
			testCase.Given.Path,
			nil)
		result := skipperMixedPredicate(ctx)
		assert.Equal(t, testCase.Expected, result)
	}
}

func TestGuardNewGroupPublic(t *testing.T) {
	cfg := test.GetTestConfig()
	require.NotNil(t, cfg)
	pathPrefix := trimVersionFromPathPrefix(cfg.Application.PathPrefix)
	version := "1.0"

	assert.PanicsWithValue(t, "'e' is nil", func() {
		guardNewGroupPublic(nil, nil, nil, nil)
	})

	e := echo.New()
	require.NotNil(t, e)
	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		guardNewGroupPublic(e.Group(pathPrefix+"/"+version), nil, nil, nil)
	})

	assert.PanicsWithValue(t, "'app' is nil", func() {
		guardNewGroupPublic(e.Group(pathPrefix+"/"+version), cfg, nil, nil)
	})

	app := handler.NewApplication(t)
	require.NotNil(t, app)
	assert.PanicsWithValue(t, "'metrics' is nil", func() {
		guardNewGroupPublic(e.Group(pathPrefix+"/"+version), cfg, app, nil)
	})

	reg := prometheus.NewRegistry()
	require.NotNil(t, reg)
	metrics := metrics.NewMetrics(reg)
	require.NotNil(t, metrics)
	assert.NotPanics(t, func() {
		guardNewGroupPublic(e.Group(pathPrefix+"/"+version), cfg, app, metrics)
	})
	app.AssertExpectations(t)
}
