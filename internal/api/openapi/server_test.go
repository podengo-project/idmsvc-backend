package openapi

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	mock_api_openapi "github.com/podengo-project/idmsvc-backend/internal/test/mock/api/openapi"
)

func TestRegisterHandlers(t *testing.T) {
	e := echo.New()
	si := ServerInterfaceWrapper{}
	assert.NotPanics(t, func() {
		RegisterHandlers(e, si.Handler)
	})
}

func TestRegisterHandlersWithBaseURL(t *testing.T) {
	const baseURL = "/api/idmsvc/v1"
	e := echo.New()
	si := ServerInterfaceWrapper{}
	assert.NotPanics(t, func() {
		RegisterHandlersWithBaseURL(e, si.Handler, baseURL)
	})
}

func TestGetOpenapi(t *testing.T) {
	// Prepare the request
	var reqReader io.Reader
	e := echo.New()
	path := "/api/idmsvc/v1/openapi.json"
	req := httptest.NewRequest(echo.GET, path, reqReader)
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Prepare the response recorder
	rec := httptest.NewRecorder()

	// Create a new context with the request and response recorder
	ctx := e.NewContext(req, rec)
	ctx.SetPath(path)

	m := &mock_api_openapi.ServerInterface{}
	h := ServerInterfaceWrapper{Handler: m}
	m.On("GetOpenapi", ctx).Return(nil)
	err := h.GetOpenapi(ctx)
	assert.NoError(t, err)
	m.AssertExpectations(t)

	m = &mock_api_openapi.ServerInterface{}
	h = ServerInterfaceWrapper{Handler: m}
	m.On("GetOpenapi", ctx).Return(fmt.Errorf("some error"))
	err = h.GetOpenapi(ctx)
	assert.EqualError(t, err, "some error")
	m.AssertExpectations(t)
}
