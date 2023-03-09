package router

import (
	"testing"

	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/handler/impl"
	"github.com/hmsidm/internal/metrics"
	"github.com/hmsidm/internal/test"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewGroupPublicPanics(t *testing.T) {
	var (
		err error
		db  *gorm.DB
	)
	// FIXME Refactor in the future; it is too complex
	const appPrefix = "/api"
	const appName = "/hmsidm"

	reg := prometheus.NewRegistry()
	require.NotNil(t, reg)
	metrics := metrics.NewMetrics(reg)
	require.NotNil(t, metrics)

	cfg := config.Get()
	_, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, db)

	// FIXME Refactor and encapsulate routerConfig in a factory function
	routerConfig := RouterConfig{
		PublicPath: appPrefix + appName,
		Handlers:   impl.NewHandler(cfg, db, metrics),
	}
	routerWrongConfig := RouterConfig{
		PublicPath: appPrefix + appName,
		Handlers:   nil,
	}
	e := echo.New()
	require.NotNil(t, e)

	assert.Panics(t, func() {
		newGroupPublic(nil, routerConfig, metrics)
	})
	assert.Panics(t, func() {
		newGroupPublic(e.Group(routerConfig.PublicPath+"/"+routerConfig.Version), routerWrongConfig, metrics)
	})
	assert.Panics(t, func() {
		newGroupPublic(e.Group(routerConfig.PublicPath+"/"+routerConfig.Version), routerConfig, nil)
	})

}

func TestNewGroupPublic(t *testing.T) {
	var (
		err error
		db  *gorm.DB
	)
	const (
		appPrefix   = "/api"
		appName     = "/hmsidm"
		versionFull = "/v1.0"
	)
	type TestCaseExpected map[string]map[string]string

	testCases := TestCaseExpected{
		"/private/readyz": {
			"GET": "github.com/hmsidm/internal/api/handler.ping",
		},
		"/private/livez": {
			"GET": "github.com/hmsidm/internal/api/public.ping",
		},

		appPrefix + appName + versionFull + "/openapi.json": {
			"GET": "github.com/hmsidm/internal/public.openapi",
		},

		appPrefix + appName + versionFull + "/domains": {
			"GET":  "github.com/hmsidm/internal/api/public.(*ServerInterfaceWrapper).ListDomains-fm",
			"POST": "github.com/hmsidm/internal/api/public.(*ServerInterfaceWrapper).CreateDomain-fm",
		},

		appPrefix + appName + versionFull + "/domains/:uuid": {
			"GET":    "github.com/hmsidm/internal/api/public.(*ServerInterfaceWrapper).ReadDomain-fm",
			"PUT":    "github.com/hmsidm/internal/api/public.(*ServerInterfaceWrapper).UpdateDomain-fm",
			"PATCH":  "github.com/hmsidm/internal/api/public.(*ServerInterfaceWrapper).PartialUpdateDomain-fm",
			"DELETE": "github.com/hmsidm/internal/api/public.(*ServerInterfaceWrapper).DeleteDomain-fm",
		},

		appPrefix + appName + versionFull + "/domains/:uuid/ipa": {
			"PUT": "github.com/hmsidm/internal/api/public.(*ServerInterfaceWrapper).RegisterIpaDomain-fm",
		},

		appPrefix + appName + versionFull + "/host-conf/:fqdn": {
			"POST": "github.com/hmsidm/internal/api/public.(*ServerInterfaceWrapper).HostConf-fm",
		},

		appPrefix + appName + versionFull + "/check-host/:subscription_manager_id/:fqdn": {
			"POST": "github.com/hmsidm/internal/api/public.(*ServerInterfaceWrapper).CheckHost-fm",
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

	cfg := config.Get()
	_, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, db)

	// FIXME Refactor and encapsulate routerConfig in a factory function
	routerConfig := RouterConfig{
		PublicPath: appPrefix + appName,
		Version:    "1.0",
		Handlers:   impl.NewHandler(cfg, db, metrics),
	}

	e := echo.New()
	require.NotNil(t, e)
	group := newGroupPublic(e.Group(routerConfig.PublicPath+"/v"+routerConfig.Version), routerConfig, metrics)
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
}
