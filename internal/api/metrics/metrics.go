// Package metrics returns a hand marshalled /metrics response. It follows
// the same pattern as the generated code from the openapi code generator.
package metrics

import "github.com/labstack/echo/v4"

// ServerInterface provides the endpoint to retrieve the metrics for this service.
type ServerInterface interface {
	// Return the metrics
	// (GET /metrics)
	GetMetrics(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetMetrics converts echo context to params.
func (w *ServerInterfaceWrapper) GetMetrics(ctx echo.Context) error {
	// Invoke the callback with all the unmarshalled arguments
	err := w.Handler.GetMetrics(ctx)
	return err
}

// EchoRouter is a simple interface which specifies additional echo.Route
// functions which are present on both echo.Echo and echo.Group, since we want
// to allow using either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "/metrics")
}

// RegisterHandlersWithBaseURL handlers, and prepends BaseURL to the paths, so
// that the paths can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}
	router.GET(baseURL, wrapper.GetMetrics)
}
