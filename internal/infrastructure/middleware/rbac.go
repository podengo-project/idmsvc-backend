package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	rbac_data "github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware/rbac-data"
	rbac_client "github.com/podengo-project/idmsvc-backend/internal/interface/client/rbac"
)

// RBACConfig hold the skipper, route prefix, the rbac permissions
// mapping for each authorized public route, and the client to
// reach out the rbac micro-service.
type RBACConfig struct {
	// Skipper function to skip for some request if necessary
	Skipper echo_middleware.Skipper
	// Prefix for the permission map
	Prefix string
	// PermissionMap has the mapping between {route,method}=>permission
	PermissionMap rbac_data.RBACMap
	// Client for rbac access
	Client rbac_client.Rbac
}

// RBACWithConfig create a middleware for authorizing requests by using
// the intgration with rbac micro-service
// rbacConfig provide the skipper, prefix, permission map and client
// for the configuration.
// Return the initialized middleware or panic if some guard condition
// is matched.
func RBACWithConfig(rbacConfig *RBACConfig) echo.MiddlewareFunc {
	if rbacConfig == nil {
		panic("'rbacConfig' is nil")
	}
	if rbacConfig.Prefix == "" {
		panic("'Prefix' is an empty string")
	}
	if rbacConfig.Client == nil {
		panic("'Client' is nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				err        error
				permission rbac_data.RBACPermission
				xrhid      string
				isAllowed  bool
			)
			logger := app_context.LogFromCtx(c.Request().Context())

			// Process skippers
			if rbacConfig.Skipper != nil && rbacConfig.Skipper(c) {
				logger.Debug("Skipping rbac for " + c.Request().RequestURI)
				return next(c)
			}

			// Get permission for the current route
			path := c.Path()
			method := c.Request().Method
			if permission, err = rbacConfig.PermissionMap.GetPermission(rbacConfig.Prefix, path, method); err != nil {
				return err
			}

			// Get X-Rh-Identity header
			// This if statement is only possible if no enforce middleware
			// is executed for the public API.
			if xrhid = c.Request().Header.Get(header.HeaderXRHID); xrhid == "" {
				return echo.NewHTTPError(http.StatusBadRequest, header.HeaderXRHID+" is missed")
			}

			// Get User permissions
			context := c.Request().Context()
			if isAllowed, err = rbacConfig.Client.IsAllowed(context, xrhid, string(permission)); !isAllowed {
				if err != nil {
					return err
				}
				logger.Error("unauthorized", "permission", permission)
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			logger.Debug("Authorized", "path", c.Request().URL.Path, "method", c.Request().Method)
			return next(c)
		}
	}
}
