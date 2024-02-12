package router

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/openapi"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
)

// skipperValidate is an alias to represent skipper for API validation middleware
// FIXME Once the openapi specification is propery defined, remove this skipper
// to validate every request
var skipperValidate echo_middleware.Skipper = skipperUserPredicate

type enforceRoute struct {
	Method string
	Path   string
}

var userEnforceRoutes = []enforceRoute{
	{"POST", "/api/idmsvc/v1/domains/token"},
	{"GET", "/api/idmsvc/v1/domains/:uuid"},
	{"GET", "/api/idmsvc/v1/domains"},
	{"PATCH", "/api/idmsvc/v1/domains/:uuid"},
	{"DELETE", "/api/idmsvc/v1/domains/:uuid"},
}

var systemEnforceRoutes = []enforceRoute{
	{"POST", "/api/idmsvc/v1/domains"},
	{"PUT", "/api/idmsvc/v1/domains/:uuid"},
	{"POST", "/api/idmsvc/v1/host-conf/:inventory_id/:fqdn"},
}

func getOpenapiPaths(c RouterConfig) func() []string {
	if c == (RouterConfig{}) {
		panic(fmt.Errorf("'c' is empty"))
	}
	if c.Version == "" {
		panic(fmt.Errorf("'c.Version' is empty"))
	}
	majorVersion := strings.Split(c.Version, ".")[0]
	fullVersion := c.Version
	cachedPaths := []string{
		fmt.Sprintf("%s/v%s/openapi.json", c.PublicPath, fullVersion),
		fmt.Sprintf("%s/v%s/openapi.json", c.PublicPath, majorVersion),
	}
	return func() []string {
		return cachedPaths
	}
}

func newGroupPublic(e *echo.Group, c RouterConfig) *echo.Group {
	if e == nil {
		panic("echo group is nil")
	}
	if c.Handlers == nil {
		panic("'handlers' is nil")
	}

	// Initialize middlewares
	var fakeIdentityMiddleware echo.MiddlewareFunc = middleware.DefaultNooperation
	if c.IsFakeEnabled {
		fakeIdentityMiddleware = middleware.FakeIdentityWithConfig(
			&middleware.FakeIdentityConfig{
				Skipper: skipperSystemPredicate,
			},
		)
	}

	systemIdentityMiddleware := middleware.EnforceIdentityWithConfig(
		&middleware.IdentityConfig{
			Skipper: skipperSystemPredicate,
			Predicates: []middleware.IdentityPredicateEntry{
				{
					Name:      "system-identity",
					Predicate: middleware.EnforceSystemPredicate,
				},
			},
		},
	)
	userIdentityMiddleware := middleware.EnforceIdentityWithConfig(
		&middleware.IdentityConfig{
			Skipper: skipperUserPredicate,
			Predicates: []middleware.IdentityPredicateEntry{
				{
					Name:      "user-identity",
					Predicate: middleware.EnforceUserPredicate,
				},
			},
		},
	)

	var rbacMiddleware echo.MiddlewareFunc
	if config.Get().Application.EnableRBAC {
		rbacMiddleware = middleware.RBACWithConfig(
			&middleware.RBACConfig{
				Skipper: nil,
				// TODO HMS-3522
				PermissionMap: middleware.RBACMap{},
			},
		)
	} else {
		rbacMiddleware = middleware.DefaultNooperation
	}
	metricsMiddleware := middleware.MetricsMiddlewareWithConfig(
		&middleware.MetricsConfig{
			Metrics: c.Metrics,
		},
	)
	requestIDMiddleware := echo_middleware.RequestIDWithConfig(
		echo_middleware.RequestIDConfig{
			TargetHeader: "X-Rh-Insights-Request-Id", // TODO Check this name is the expected
		},
	)
	var validateAPI echo.MiddlewareFunc = middleware.DefaultNooperation
	if c.EnableAPIValidator {
		middleware.InitOpenAPIFormats()
		validateAPI = middleware.RequestResponseValidatorWithConfig(
			// FIXME Get the values from the application config
			&middleware.RequestResponseValidatorConfig{
				Skipper:          nil,
				ValidateRequest:  true,
				ValidateResponse: false,
			},
		)
	}

	// Wire the middlewares
	e.Use(
		middleware.CreateContext(),
		fakeIdentityMiddleware,
		systemIdentityMiddleware,
		userIdentityMiddleware,
		rbacMiddleware,
		metricsMiddleware,
		echo_middleware.Secure(),
		// TODO Check if this is made by 3scale
		// middleware.CORSWithConfig(middleware.CORSConfig{}),
		requestIDMiddleware,
		validateAPI,
	)

	// Setup routes
	public.RegisterHandlersWithBaseURL(e, c.Handlers, "")
	openapi.RegisterHandlersWithBaseURL(e, c.Handlers, "")
	return e
}

// skipperUserPredicate applied when using EnforceUserPredicate.
// ctx is the request context.
// Return true if enforce identity is skipped, else false.
func skipperUserPredicate(ctx echo.Context) bool {
	var r enforceRoute
	path := ctx.Path()
	method := ctx.Request().Method
	// it is not expected a big number of routes, but if that were
	// the case into the future, it is more efficient to check
	// directly against a hashmap instead of traversing the slice
	for i := range userEnforceRoutes {
		r = userEnforceRoutes[i]
		if method == r.Method && path == r.Path {
			return false
		}

	}
	return true
}

// skipperSystemPredicate applied when using EnforceSystemPredicate.
// ctx is the request context.
// Return true if enforce identity is skipped, else false.
func skipperSystemPredicate(ctx echo.Context) bool {
	// Read the route path __pattern__ that matched this request
	var r enforceRoute
	path := ctx.Path()
	method := ctx.Request().Method
	// it is not expected a big number of routes, but if that were
	// the case into the future, it is more efficient to check
	// directly against a hashmap instead of traversing the slice
	for i := range systemEnforceRoutes {
		r = systemEnforceRoutes[i]
		if method == r.Method && path == r.Path {
			return false
		}
	}
	return true
}

// skipperOpenapi skip /api/idmsvc/v*/openapi.json path
func newSkipperOpenapi(c RouterConfig) echo_middleware.Skipper {
	paths := getOpenapiPaths(c)()
	return func(ctx echo.Context) bool {
		route := ctx.Path()
		for i := range paths {
			if paths[i] == route {
				return true
			}
		}
		return false
	}
}
