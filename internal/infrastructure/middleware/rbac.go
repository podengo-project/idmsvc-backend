package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

// TODO HMS-3522 Add RBACMap and implement its methods:
//
// Add(route string, method string, permission string) RBACMap
// Get(route string, method string) string

type route string
type method string
type permission string

// RBACMap is a mapping for the permissions [path][method] => permission
type RBACMap map[route]map[method]permission

type RBACConfig struct {
	// Skipper function to skip for some request if necessary
	Skipper echo_middleware.Skipper
	// TODO HMS-3522 Add the mapping structure
	PermissionMap RBACMap
}

func RBACWithConfig(config *RBACConfig) echo.MiddlewareFunc {
	if config == nil {
		panic("'config' is nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO Implement HMS-3522
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
	}
}
