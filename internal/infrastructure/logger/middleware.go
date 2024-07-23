package logger

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
)

// This requires the following values to be set in
// middleware.RequestLoggerWithConfig:
//
// LogError:  true,
// LogMethod: true,
// LogStatus: true,
// LogURI:    true,
func MiddlewareLogValues(c echo.Context, v middleware.RequestLoggerValues) error {
	var logLevel slog.Level
	var logAttr []slog.Attr = make([]slog.Attr, 5)

	ctx := c.Request().Context()
	log := app_context.LogFromCtx(ctx)

	req := c.Request()
	res := c.Response()

	request_id := req.Header.Get(header.HeaderXRequestID)
	if request_id == "" {
		request_id = res.Header().Get(header.HeaderXRequestID)
	}

	logAttr = append(logAttr,
		slog.String("method", v.Method),
		slog.String("uri", v.URI),
		slog.Int("status", v.Status),
	)
	if v.Error == nil {
		logLevel = slog.LevelInfo
	} else {
		logLevel = slog.LevelError
		logAttr = append(logAttr, slog.String("err", v.Error.Error()))
	}

	log.LogAttrs(
		ctx,
		logLevel,
		"http_request",
		logAttr...,
	)

	return nil
}
