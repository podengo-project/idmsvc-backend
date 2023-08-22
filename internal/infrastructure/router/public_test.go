package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
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

	cfg := config.Get()
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
			"GET": "github.com/podengo-project/idmsvc-backend/internal/public.openapi",
		},

		appPrefix + appName + versionFull + "/domains": {
			"GET":  "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).ListDomains-fm",
			"POST": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).CreateDomain-fm",
		},

		appPrefix + appName + versionFull + "/domains/:uuid": {
			"GET":    "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).ReadDomain-fm",
			"PUT":    "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).UpdateDomain-fm",
			"PATCH":  "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).PartialUpdateDomain-fm",
			"DELETE": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).DeleteDomain-fm",
		},

		appPrefix + appName + versionFull + "/domains/:uuid/register": {
			"PUT": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).RegisterDomain-fm",
		},

		appPrefix + appName + versionFull + "/domains/:uuid/update": {
			"PUT": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).UpdateDomain-fm",
		},

		appPrefix + appName + versionFull + "/host-conf/:inventory_id/:fqdn": {
			"POST": "github.com/podengo-project/idmsvc-backend/internal/api/public.(*ServerInterfaceWrapper).HostConf-fm",
		},

		// This routes are added when the group is created
		appPrefix + appName + versionFull + "/*": {
			"GET":      "github.com/labstack/echo/v4.glob..func1",
			"POST":     "github.com/labstack/echo/v4.glob..func1",
			"PUT":      "github.com/labstack/echo/v4.glob..func1",
			"PATCH":    "github.com/labstack/echo/v4.glob..func1",
			"DELETE":   "github.com/labstack/echo/v4.glob..func1",
			"HEAD":     "github.com/labstack/echo/v4.glob..func1",
			"OPTIONS":  "github.com/labstack/echo/v4.glob..func1",
			"REPORT":   "github.com/labstack/echo/v4.glob..func1",
			"PROPFIND": "github.com/labstack/echo/v4.glob..func1",
			"TRACE":    "github.com/labstack/echo/v4.glob..func1",
			"CONNECT":  "github.com/labstack/echo/v4.glob..func1",
		},
		appPrefix + appName + versionFull: {
			"GET":      "github.com/labstack/echo/v4.glob..func1",
			"POST":     "github.com/labstack/echo/v4.glob..func1",
			"PUT":      "github.com/labstack/echo/v4.glob..func1",
			"PATCH":    "github.com/labstack/echo/v4.glob..func1",
			"DELETE":   "github.com/labstack/echo/v4.glob..func1",
			"HEAD":     "github.com/labstack/echo/v4.glob..func1",
			"OPTIONS":  "github.com/labstack/echo/v4.glob..func1",
			"REPORT":   "github.com/labstack/echo/v4.glob..func1",
			"PROPFIND": "github.com/labstack/echo/v4.glob..func1",
			"TRACE":    "github.com/labstack/echo/v4.glob..func1",
			"CONNECT":  "github.com/labstack/echo/v4.glob..func1",
		},
	}

	reg := prometheus.NewRegistry()
	require.NotNil(t, reg)
	metrics := metrics.NewMetrics(reg)
	require.NotNil(t, metrics)
	inventory := client.NewHostInventory(t)
	require.NotNil(t, inventory)

	cfg := config.Get()
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
		assert.Truef(t, okPath, "path=%s not found into the expected ones", route.Path)

		name, okMethod := methods[route.Method]
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
		t.Log(testCase.Name)
		ctx := helperNewContextForSkipper(
			testCase.Given,
			http.MethodGet,
			testCase.Given,
			nil)
		result := skipperUserPredicate(ctx)
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
		t.Log(testCase.Name)
		ctx := helperNewContextForSkipper(testCase.Given, http.MethodGet, testCase.Given, nil)
		result := skipperSystemPredicate(ctx)
		assert.Equal(t, testCase.Expected, result)
	}
}
