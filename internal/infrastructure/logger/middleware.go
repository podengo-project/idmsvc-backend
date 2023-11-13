package logger

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
)

const (
	HeaderXRHRequestID = "X-Rh-Insights-Request-Id"
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

	req := c.Request()
	res := c.Response()

	request_id := req.Header.Get(HeaderXRHRequestID)
	if request_id == "" {
		request_id = res.Header().Get(HeaderXRHRequestID)
	}

	logAttr = append(logAttr,
		slog.String("request_id", request_id),
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

	slog.LogAttrs(
		context.Background(),
		logLevel,
		"http_request",
		logAttr...,
	)

	return nil
}
