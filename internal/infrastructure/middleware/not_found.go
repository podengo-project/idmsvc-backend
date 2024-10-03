package middleware

import (
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
)

// NotFound middleware is a workaround to avoid the middleware chain
// is executed when the path is not found, because at the moment of
// writting this, the not found paths are matching a default handler.
//
// This workaround evoke an early return by calling the NotFoundHandler
func NotFound() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			if path == "" || strings.HasSuffix(path, "/*") {
				logger := context.LogFromCtx(c.Request().Context())
				logger.Error("not found", slog.String("path", path))
				return echo.NotFoundHandler(c)
			}
			return next(c)
		}
	}
}
