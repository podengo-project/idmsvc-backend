package router

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
)

type RouterConfig struct {
	Handlers           handler.Application
	PublicPath         string
	PrivatePath        string
	Version            string
	MetricsPath        string
	IsFakeEnabled      bool
	EnableAPIValidator bool
	Metrics            *metrics.Metrics
}

func getMajorVersion(version string) string {
	if version == "" {
		return ""
	}
	return strings.Split(version, ".")[0]
}

func checkRouterConfig(c RouterConfig) error {
	if c.PublicPath == "" {
		return fmt.Errorf("PublicPath cannot be empty")
	}
	if c.PrivatePath == "" {
		return fmt.Errorf("PrivatePath cannot be empty")
	}
	if c.PublicPath == c.PrivatePath {
		return fmt.Errorf("PublicPath and PrivatePath cannot be equal")
	}
	if c.Version == "" {
		return fmt.Errorf("Version cannot be empty")
	}
	if c.Metrics == nil {
		return fmt.Errorf("Metrics is nil")
	}
	return nil
}

func loggerSkipperWithPaths(paths ...string) middleware.Skipper {
	return func(c echo.Context) bool {
		path := c.Path()
		for _, item := range paths {
			if item == path {
				return true
			}
		}
		return false
	}
}

func configCommonMiddlewares(e *echo.Echo, c RouterConfig) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		// Request logger values for middleware.RequestLoggerValues
		LogError:  true,
		LogMethod: true,
		LogStatus: true,
		LogURI:    true,

		// Forwards error to the global error handler, so it can decide
		// appropriate status code.
		HandleError: true,

		Skipper: loggerSkipperWithPaths(
			c.MetricsPath,
			c.PrivatePath+"/readyz",
			c.PrivatePath+"/livez",
		),

		LogValuesFunc: logger.MiddlewareLogValues,
	}))

	e.Use(middleware.Recover())
}

// NewRouterWithConfig fill the router configuration for the given echo instance,
// providing routes for the public endpoints, the private paths (includes the healthcheck),
// and the /metrics path
// e is the echo instance where to add the routes.
// c is the router configuration.
// metrics is the reference to the metrics storage.
// Return the echo instance set up; is something fails it panics.
func NewRouterWithConfig(e *echo.Echo, c RouterConfig) *echo.Echo {
	if e == nil {
		panic("'e' is nil")
	}
	if err := checkRouterConfig(c); err != nil {
		panic(err.Error())
	}

	configCommonMiddlewares(e, c)

	newGroupPrivate(e.Group(c.PrivatePath), c)
	newGroupPublic(e.Group(c.PublicPath+"/v"+c.Version), c)
	newGroupPublic(e.Group(c.PublicPath+"/v"+getMajorVersion(c.Version)), c)
	return e
}

// NewRouterForMetrics fill the routing information for /metrics endpoint.
// e is the echo instance
// c is the router configuration
// Return the echo instance configured for the metrics for success execution,
// else raise any panic.
func NewRouterForMetrics(e *echo.Echo, c RouterConfig) *echo.Echo {
	if e == nil {
		panic("'e' is nil")
	}
	if c.MetricsPath == "" {
		panic(fmt.Errorf("MetricsPath cannot be an empty string"))
	}

	configCommonMiddlewares(e, c)

	// Register handlers
	newGroupMetrics(e.Group(c.MetricsPath), c)
	return e
}
