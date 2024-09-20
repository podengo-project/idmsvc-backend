package router

import (
	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/private"
)

func guardNewGroupPrivate(e *echo.Group, apiPrivate private.ServerInterface) {
	if e == nil {
		panic("'e' is nil")
	}
	if apiPrivate == nil {
		panic("'apiPrivate' is nil")
	}
}

func newGroupPrivate(e *echo.Group, apiPrivate private.ServerInterface) *echo.Group {
	guardNewGroupPrivate(e, apiPrivate)
	private.RegisterHandlers(e, apiPrivate)
	return e
}
