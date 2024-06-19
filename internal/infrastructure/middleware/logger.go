package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
)

const (
	successRequest = "success request"
)

type LogConfig struct {
	Skipper middleware.Skipper
}

func ContextLogConfig(cfg *LogConfig) echo.MiddlewareFunc {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper(c) {
				return next(c)
			}

			requestID := c.Request().Header.Get(header.HeaderXRequestID)

			// Let print the request-id for every call of the
			// request log
			logger := slog.Default().With(
				slog.String("request-id", requestID),
			)
			ctx := c.Request().Context()
			ctx = app_context.CtxWithLog(ctx, logger)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			// Splitted in two lines for a better debugging
			// experience
			err := next(c)
			return err
		}
	}
}
