package middleware

import (
	"context"
	"log/slog"
	"net/http"

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

func LoggerWithConfig(cfg *LogConfig) echo.MiddlewareFunc {
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

			var logFunc func(ctx context.Context, msg string, args ...any)
			requestID := c.Request().Header.Get(header.HeaderXRequestID)
			method := c.Request().Method
			path := c.Request().RequestURI

			// Let print the request-id for every call of the
			// request log
			s := slog.Default().With(
				slog.String("request-id", requestID),
			)
			ctx := c.Request().Context()
			ctx = app_context.CtxWithLog(ctx, s)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			err := next(c)
			status := c.Response().Status
			msg := http.StatusText(status)
			if err != nil {
				msg = err.Error()
				logFunc = s.ErrorContext
			} else if status >= 400 {
				logFunc = s.ErrorContext
			} else {
				logFunc = s.InfoContext
			}

			logFunc(ctx,
				msg,
				"method", method,
				"path", path,
				"status", status,
			)
			return err
		}
	}
}
