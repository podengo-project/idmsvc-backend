package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

// NewHandlerTester will help to build a temporary echo service which
// will mock the response we want to test the inventory client component.
// method indicate the http method to mock, such as "GET", "PUT", "DELETE"
// path the url path we need to mock.
// requestParams hashmap for the parameters value.
// body the payload string that we want to receive for the above request.
// Return an echo framework initialized and ready to start.
func NewHandlerTester(t *testing.T,
	method string,
	path string,
	routerPath string,
	requestParams map[string]string,
	requestPayload string,
	requestHeaders map[string]string,
	responseStatus int,
	responseBody string,
	responseHeaders map[string]string,
	middleware ...echo.MiddlewareFunc,
) (*echo.Echo, echo.Context, *httptest.ResponseRecorder) {
	// See: https://echo.labstack.com/guide/testing/

	// Create echo instance
	e := echo.New()

	// Add middleware if any
	e.Use(middleware...)

	// Prepare the request
	var reqReader io.Reader
	if requestPayload != "" {
		reqReader = strings.NewReader(requestPayload)
	}
	req := httptest.NewRequest(method, path, reqReader)
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for key, value := range requestHeaders {
		req.Header.Add(key, value)
	}

	// Prepare the response recorder
	rec := httptest.NewRecorder()

	// Add the mock handler
	e.Add(method, path, createMockHandler(responseStatus, responseBody, responseHeaders))

	// Create a new context with the request and response recorder
	ctx := e.NewContext(req, rec)
	if routerPath != "" {
		ctx.SetPath(routerPath)
	} else {
		ctx.SetPath(path)
	}

	// Prepare the request parameters if any
	setRequestParams(ctx, requestParams)

	// Start Server to listen in a go routine
	server := http.Server{}
	go func() {
		e.HideBanner = true
		e.HidePort = true
		_ = e.StartServer(&server)
	}()
	// Await until the listener is created, so it
	// can start to receive requests
	for {
		if e.Listener != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	return e, ctx, rec
}

func createMockHandler(responseStatus int, responseBody string, responseHeaders map[string]string) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		ctx.Response().Status = responseStatus
		for key, value := range responseHeaders {
			ctx.Response().Header().Add(key, value)
		}
		if responseStatus >= 200 && responseStatus < 300 {
			if responseBody != "" {
				return ctx.JSONBlob(responseStatus, []byte(responseBody))
			}
			return ctx.NoContent(responseStatus)
		} else {
			err := echo.NewHTTPError(responseStatus, responseBody)
			return err
		}
	}
}

func setRequestParams(ctx echo.Context, requestParams map[string]string) {
	if ctx == nil {
		panic("'ctx' cannot be nil")
	}
	if requestParams != nil {
		keys := make([]string, 0, len(requestParams))
		for k := range requestParams {
			keys = append(keys, k)
		}
		values := make([]string, 0, len(requestParams))
		for _, v := range requestParams {
			values = append(values, v)
		}
		ctx.SetParamNames(keys...)
		ctx.SetParamValues(values...)
	}
}
