package router

import (
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/private"
)

func newGroupPrivate(e *echo.Group, c RouterConfig) *echo.Group {
	private.RegisterHandlers(e, c.Handlers)
	return e
}
