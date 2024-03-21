package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler/impl"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/podengo-project/idmsvc-backend/internal/test/mock/interface/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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

func TestNewGroupPublicPanics(t *testing.T) {
	var (
		err error
		db  *gorm.DB
	)
	// FIXME Refactor in the future; it is too complex
	const appPrefix = "/api"
	const appName = "/idmsvc"

	reg := prometheus.NewRegistry()
	require.NotNil(t, reg)
	metrics := metrics.NewMetrics(reg)
	require.NotNil(t, metrics)
	inventory := client.NewHostInventory(t)

	cfg := test.GetTestConfig()
	_, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, db)

	// FIXME Refactor and encapsulate routerConfig in a factory function
	routerConfig := RouterConfig{
		PublicPath: appPrefix + appName,
		Handlers:   impl.NewHandler(cfg, db, metrics, inventory),
		Metrics:    metrics,
	}
	routerWrongConfig := RouterConfig{
		PublicPath: appPrefix + appName,
		Handlers:   nil,
	}
	e := echo.New()
	require.NotNil(t, e)

	assert.Panics(t, func() {
		newGroupPublic(nil, routerConfig)
	})
	assert.Panics(t, func() {
		newGroupPublic(e.Group(routerConfig.PublicPath+"/"+routerConfig.Version), routerWrongConfig)
	})
	assert.NotPanics(t, func() {
		newGroupPublic(e.Group(routerConfig.PublicPath+"/"+routerConfig.Version), routerConfig)
	})
}

func TestNewGroupPublic(t *testing.T) {
	var (
		err error
		db  *gorm.DB
	)
	const (
		appPrefix   = "/api"
		appName     = "/idmsvc"
		versionFull = "/v1.0"
	)
	type TestCaseExpected map[string]map[string]string

	testCases := TestCaseExpected{
		"/private/readyz": {
			"GET": "github.com/podengo-project/idmsvc-backend/internal/api/handler.ping",
		},
		"/private/livez": {
			"GET": "github.com/podengo-project/idmsvc-backend/internal/api/public.ping",
		},

		appPrefix + appName + versionFull + "/openapi.json": {
			"GET": "github.com/podengo-project/idmsvc-backend/internal/api/openapi.(*ServerInterfaceWrapper).GetOpenapi-fm",
		},

		appPrefix + appName + versionFull + "/domains": {
			"GET":  "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).ListDomains-fm",
			"POST": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).RegisterDomain-fm",
		},

		appPrefix + appName + versionFull + "/domains/token": {
			"POST": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).CreateDomainToken-fm",
		},

		appPrefix + appName + versionFull + "/domains/:uuid": {
			"GET":    "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).ReadDomain-fm",
			"PUT":    "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).UpdateDomainAgent-fm",
			"PATCH":  "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).UpdateDomainUser-fm",
			"DELETE": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).DeleteDomain-fm",
		},

		appPrefix + appName + versionFull + "/host-conf/:inventory_id/:fqdn": {
			"POST": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).HostConf-fm",
		},

		appPrefix + appName + versionFull + "/signing_keys": {
			"GET": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).GetSigningKeys-fm",
		},

		// This routes are added when the group is created
		appPrefix + appName + versionFull + "/*": {
			"echo_route_not_found": "github.com/labstack/echo/v4.glob..func1",
		},
		appPrefix + appName + versionFull: {
			"echo_route_not_found": "github.com/labstack/echo/v4.glob..func1",
		},
	}

	reg := prometheus.NewRegistry()
	require.NotNil(t, reg)
	metrics := metrics.NewMetrics(reg)
	require.NotNil(t, metrics)
	inventory := client.NewHostInventory(t)
	require.NotNil(t, inventory)

	cfg := test.GetTestConfig()
	_, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, db)

	// FIXME Refactor and encapsulate routerConfig in a factory function
	routerConfig := RouterConfig{
		PublicPath: appPrefix + appName,
		Version:    "1.0",
		Handlers:   impl.NewHandler(cfg, db, metrics, inventory),
		Metrics:    metrics,
	}

	e := echo.New()
	require.NotNil(t, e)
	group := newGroupPublic(e.Group(routerConfig.PublicPath+"/v"+routerConfig.Version), routerConfig)
	require.NotNil(t, group)

	// Match Routes in expected
	for _, route := range e.Routes() {
		t.Logf("Method=%s Path=%s Name=%s", route.Method, route.Path, route.Name)

		methods, okPath := testCases[route.Path]
		if !okPath {
			t.Logf("Path=%s not found", route.Path)
		}
		assert.Truef(t, okPath, "path=%s not found into the expected ones", route.Path)

		name, okMethod := methods[route.Method]
		if !okMethod {
			t.Logf("Method=%s not found for path=%s", route.Method, route.Path)
		}
		assert.Truef(t, okMethod, "method=%s not found into the expected ones for the path=%s", route.Method, route.Path)
		assert.Equalf(t, name, route.Name, "handler for path=%s method=%s does not match", route.Path, route.Method)
	}

	// Same result when IsFakeEnabled
	e = echo.New()
	require.NotNil(t, e)
	routerConfig.IsFakeEnabled = true
	group = newGroupPublic(
		e.Group(routerConfig.PublicPath+"/v"+routerConfig.Version),
		routerConfig)
	require.NotNil(t, group)
	for _, route := range e.Routes() {
		t.Logf("Method=%s Path=%s Name=%s", route.Method, route.Path, route.Name)

		methods, okPath := testCases[route.Path]
		assert.Truef(t, okPath, "path=%s not found into the expected ones", route.Path)

		name, okMethod := methods[route.Method]
		assert.Truef(t,
			okMethod,
			"method=%s not found into the expected ones for the path=%s",
			route.Method,
			route.Path)
		assert.Equalf(t,
			name,
			route.Name,
			"handler for path=%s method=%s does not match",
			route.Path,
			route.Method)
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
	assert.PanicsWithError(t, "'c' is empty", func() {
		getOpenapiPaths(RouterConfig{})
	})

	assert.PanicsWithError(t, "'c.Version' is empty", func() {
		getOpenapiPaths(RouterConfig{
			Version:    "",
			PublicPath: "/api/idmsvc",
		})
	})

	cachedPaths := getOpenapiPaths(RouterConfig{
		Version:    "1.4",
		PublicPath: "/api/idmsvc",
	})
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
	assert.PanicsWithError(t, "'c' is empty", func() {
		newSkipperOpenapi(RouterConfig{})
	})

	c := RouterConfig{
		PublicPath: "/api/idmsvc",
		Version:    "1.4",
	}
	skipper := newSkipperOpenapi(c)
	assert.NotNil(t, skipper)

	path := fmt.Sprintf("%s/v%s/openapi.json", c.PublicPath, c.Version)
	ctx := helperNewContextForSkipper(path, echo.GET, path, map[string]string{})
	assert.True(t, skipper(ctx))
	path = fmt.Sprintf("%s/v%s/openapi.json", c.PublicPath, strings.Split(c.Version, ".")[0])
	ctx = helperNewContextForSkipper(path, echo.GET, path, map[string]string{})
	assert.True(t, skipper(ctx))
	path = fmt.Sprintf("%s/v%s/openapi2.json", c.PublicPath, strings.Split(c.Version, ".")[0])
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
