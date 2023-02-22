package router

import (
	"fmt"
	"strings"

	"github.com/hmsidm/internal/handler"
	"github.com/hmsidm/internal/metrics"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouterConfig struct {
	Handlers    handler.Application
	PublicPath  string
	PrivatePath string
	Version     string
	MetricsPath string
	Metrics     *metrics.Metrics
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
	if c.Version[0] != 'v' {
		return fmt.Errorf("Version should follow the pattern 'v<Major>.<Minor>'")
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
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Skipper: loggerSkipperWithPaths(
				c.MetricsPath,
				c.PrivatePath+"/readyz",
				c.PrivatePath+"/livez",
			),
		},
	))
	e.Use(middleware.Recover())
}

func NewRouterWithConfig(e *echo.Echo, c RouterConfig, metrics *metrics.Metrics) *echo.Echo {
	if err := checkRouterConfig(c); err != nil {
		panic(err.Error())
	}

	configCommonMiddlewares(e, c)

	newGroupPrivate(e.Group(c.PrivatePath), c)
	newGroupPublic(e.Group(c.PublicPath+"/v"+c.Version), c, metrics)
	newGroupPublic(e.Group(c.PublicPath+"/v"+getMajorVersion(c.Version)), c, metrics)
	return e
}

func NewRouterForMetrics(e *echo.Echo, c RouterConfig) *echo.Echo {
	if c.MetricsPath == "" {
		panic(fmt.Errorf("MetricsPath cannot be an empty string"))
	}

	configCommonMiddlewares(e, c)

	// Register handlers
	newGroupMetrics(e.Group(c.MetricsPath), c)
	return e
}
