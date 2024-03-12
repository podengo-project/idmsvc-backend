package rbac_data

import (
	echo_middleware "github.com/labstack/echo/v4/middleware"
	rbac_client "github.com/podengo-project/idmsvc-backend/internal/interface/client/rbac"
)

type (
	Route      string
	Method     string
	Permission string
)

type RBACConfig struct {
	// Skipper function to skip for some request if necessary
	Skipper echo_middleware.Skipper
	// Prefix for the permission map
	Prefix string
	// PermissionMap has the mapping between {route,method}=>permission
	PermissionMap RBACMap
	// Client for rbac access
	Client rbac_client.Rbac
}
