package router

import (
	"github.com/hmsidm/internal/api/private"
	"github.com/labstack/echo/v4"
)

func newGroupPrivate(e *echo.Group, c RouterConfig) *echo.Group {
	private.RegisterHandlers(e, c.Handlers)
	return e
}
