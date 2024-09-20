package router

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	api_metrics "github.com/podengo-project/idmsvc-backend/internal/api/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	app_middleware "github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
)

func getMajorVersion(version string) string {
	if version == "" {
		return ""
	}
	return strings.Split(version, ".")[0]
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

func configCommonMiddlewares(e *echo.Echo, cfg *config.Config) {
	privatePath := "/private"
	metricsPath := cfg.Metrics.Path
	e.Pre(middleware.RemoveTrailingSlash())

	skipperPaths := []string{
		privatePath + "/readyz",
		privatePath + "/livez",
		metricsPath,
	}

	e.Use(app_middleware.ContextLogConfig(&app_middleware.LogConfig{
		Skipper: loggerSkipperWithPaths(skipperPaths...),
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		// Request logger values for middleware.RequestLoggerValues
		LogError:  true,
		LogMethod: true,
		LogStatus: true,
		LogURI:    true,

		// We need to set HandleError to false, to avoid double execution, first by the
		// logger middleware, second by the echo framework internals.
		// - https://github.com/labstack/echo/blob/v4.12.0/middleware/request_logger.go#L287
		// - https://github.com/labstack/echo/blob/v4.12.0/echo.go#L674
		HandleError: false,

		Skipper: loggerSkipperWithPaths(skipperPaths...),

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
func NewRouterWithConfig(e *echo.Echo, cfg *config.Config, app handler.Application, metrics *metrics.Metrics) *echo.Echo {
	guardNewRouterWithConfig(e, cfg, app, metrics)
	// TODO Add version to the configuration, an set it from config.example.yaml
	// or clowder.yaml deployment descriptor
	version := "1.0"
	privatePath := "/private"
	publicPath := trimVersionFromPathPrefix(cfg.Application.PathPrefix)

	configCommonMiddlewares(e, cfg)

	newGroupPrivate(e.Group(privatePath), app)
	newGroupPublic(e.Group(publicPath+"/v"+version), cfg, app, metrics)
	newGroupPublic(e.Group(publicPath+"/v"+getMajorVersion(version)), cfg, app, metrics)
	return e
}

func guardNewRouterWithConfig(e *echo.Echo, cfg *config.Config, app handler.Application, metrics *metrics.Metrics) {
	if e == nil {
		panic("'e' is nil")
	}
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if app == nil {
		panic("'app' is nil")
	}
	if metrics == nil {
		panic("'metrics' is nil")
	}
}

func trimVersionFromPathPrefix(pathPrefix string) string {
	if pathPrefix == "" {
		return pathPrefix
	}
	pathPrefix = strings.TrimSuffix(pathPrefix, "/")
	pathPrefixItems := strings.Split(pathPrefix, "/")
	lenItems := len(pathPrefixItems)
	if len(pathPrefixItems) > 0 {
		item := pathPrefixItems[lenItems-1]
		if item[0] == 'v' {
			pathPrefix = strings.Join(pathPrefixItems[:lenItems-1], "/")
			return pathPrefix
		}
	}
	return pathPrefix
}

func guardNewRouterForMetrics(e *echo.Echo, cfg *config.Config, apiMetrics api_metrics.ServerInterface) {
	if e == nil {
		panic("'e' is nil")
	}
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if cfg.Metrics.Path == "" {
		panic("'MetricsPath' cannot be an empty string")
	}
	if apiMetrics == nil {
		panic("'apiMetrics' is nil")
	}
}

// NewRouterForMetrics fill the routing information for /metrics endpoint.
// e is the echo instance
// cfg is the router configuration
// Return the echo instance configured for the metrics for success execution,
// else raise any panic.
func NewRouterForMetrics(e *echo.Echo, cfg *config.Config, app api_metrics.ServerInterface) *echo.Echo {
	guardNewRouterForMetrics(e, cfg, app)
	configCommonMiddlewares(e, cfg)
	return newGroupMetrics(e, cfg, app)
}
