package router

import (
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/infrastructure/middleware"
	"github.com/hmsidm/internal/metrics"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

var userEnforceRoutes = []string{
	"/api/hmsidm/v1/domains",
	"/api/hmsidm/v1/domains/:uuid",
}

var systemEnforceRoutes = []string{
	"/api/hmsidm/v1/domains/:uuid/register",
	"/api/hmsidm/v1/domains/:uuid/update",
}

func newGroupPublic(e *echo.Group, c RouterConfig, metrics *metrics.Metrics) *echo.Group {
	if e == nil {
		panic("echo group is nil")
	}
	if metrics == nil {
		panic("'metrics' is nil")
	}
	if c.Handlers == nil {
		panic("'handlers' is nil")
	}

	// Initialize middlewares
	systemIdentityMiddleware := middleware.EnforceIdentityWithConfig(
		&middleware.IdentityConfig{
			Skipper: skipperSystemPredicate,
			Predicates: map[string]middleware.IdentityPredicate{
				"system-identity": middleware.EnforceSystemPredicate,
			},
		},
	)
	userIdentityMiddleware := middleware.EnforceIdentityWithConfig(
		&middleware.IdentityConfig{
			Skipper: skipperUserPredicate,
			Predicates: map[string]middleware.IdentityPredicate{
				"user-identity": middleware.EnforceUserPredicate,
			},
		},
	)
	metricsMiddleware := middleware.MetricsMiddlewareWithConfig(
		&middleware.MetricsConfig{
			Metrics: metrics,
		},
	)
	requestIDMiddleware := echo_middleware.RequestIDWithConfig(
		echo_middleware.RequestIDConfig{
			TargetHeader: "X-Rh-Insights-Request-Id", // TODO Check this name is the expected
		},
	)

	// Wire the middlewares
	e.Use(
		middleware.CreateContext(),
		systemIdentityMiddleware,
		userIdentityMiddleware,
		metricsMiddleware,
		echo_middleware.Secure(),
		// TODO Check if this is made by 3scale
		// middleware.CORSWithConfig(middleware.CORSConfig{}),
		requestIDMiddleware,
		// FIXME Investigate why is failing when it is uncommented
		// middleware.NewApiServiceValidator(),
	)

	// Setup routes
	public.RegisterHandlersWithBaseURL(e, c.Handlers, "")
	return e
}

// skipperUserPredicate applied when using EnforceUserPredicate.
// ctx is the request context.
// Return true if enforce identity is skipped, else false.
func skipperUserPredicate(ctx echo.Context) bool {

	route := ctx.Path()
	// it is not expected a big number of routes, but if that were
	// the case into the future, it is more efficient to check
	// directly against a hashmap instead of traversing the slice
	for i := range userEnforceRoutes {
		if route == userEnforceRoutes[i] {
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
	route := ctx.Path()
	// it is not expected a big number of routes, but if that were
	// the case into the future, it is more efficient to check
	// directly against a hashmap instead of traversing the slice
	for i := range systemEnforceRoutes {
		if route == systemEnforceRoutes[i] {
			return false
		}
	}
	return true
}
