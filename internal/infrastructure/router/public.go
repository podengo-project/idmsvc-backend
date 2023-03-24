package router

import (
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/infrastructure/middleware"
	"github.com/hmsidm/internal/metrics"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

func newGroupPublic(e *echo.Group, c RouterConfig, metrics *metrics.Metrics) *echo.Group {
	if e == nil {
		panic("no echo group was specified")
	}
	if metrics == nil {
		panic("no metrics was specified")
	}
	if c.Handlers == nil {
		panic("handlers not specified in the router configuration")
	}

	// Set up middlewares
	e.Use(middleware.CreateContext())
	e.Use(middleware.EnforceIdentityWithConfig(middleware.NewIdentityConfig().
		SetSkipper(middleware.SkipperUserPredicate).
		AddPredicate("user-predicate", middleware.EnforceUserPredicate)))
	e.Use(middleware.EnforceIdentityWithConfig(middleware.NewIdentityConfig().
		SetSkipper(middleware.SkipperSystemPredicate).
		AddPredicate("system-predicate", middleware.EnforceSystemPredicate)))
	e.Use(middleware.MetricsMiddlewareWithConfig(&middleware.MetricsConfig{
		Metrics: metrics,
	}))
	e.Use(echo_middleware.Secure())
	// TODO Check if this is made by 3scale
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
	e.Use(echo_middleware.RequestIDWithConfig(echo_middleware.RequestIDConfig{
		TargetHeader: "X-Rh-Insights-Request-Id", // TODO Check this name is the expected
	}))
	// FIXME Investigate why is failing when it is uncommented
	// e.Use(middleware.NewApiServiceValidator())

	// Setup routes
	public.RegisterHandlersWithBaseURL(e, c.Handlers, "")
	return e
}
