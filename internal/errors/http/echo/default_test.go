package echo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func helperNewEcho(err error) *echo.Echo {
	e := echo.New()
	e.Use(middleware.ContextLogConfig(&middleware.LogConfig{}))
	e.HTTPErrorHandler = DefaultErrorHandler
	e.GET("/", func(_ echo.Context) error {
		return err
	})
	return e
}

func TestDefaultErrorHandler(t *testing.T) {
	type TestCaseExpected struct {
		Code int
		Body api_public.ErrorResponse
	}
	type TestCase struct {
		Name     string
		Given    error
		Expected TestCaseExpected
	}

	testCases := []TestCase{
		{
			Name:  "Raw error",
			Given: fmt.Errorf("Test error"),
			Expected: TestCaseExpected{
				Code: http.StatusInternalServerError,
				Body: *builder_api.NewErrorResponse().
					Add(*builder_api.NewErrorInfo(http.StatusInternalServerError).
						WithTitle(http.StatusText(http.StatusInternalServerError)).
						WithDetail("Test error").
						Build(),
					).Build(),
			},
		},
		{
			Name:  "Not nested echo.HTTPError",
			Given: echo.NewHTTPError(http.StatusNotFound, "Not Found"),
			Expected: TestCaseExpected{
				Code: http.StatusNotFound,
				Body: *builder_api.NewErrorResponse().
					Add(*builder_api.NewErrorInfo(http.StatusNotFound).
						WithTitle("Not Found").
						Build(),
					).Build(),
			},
		},
		{
			Name:  "Nested echo.HTTPError",
			Given: fmt.Errorf("primary error: %w", fmt.Errorf("secondary error")),
			Expected: TestCaseExpected{
				Code: http.StatusInternalServerError,
				Body: *builder_api.NewErrorResponse().
					Add(*builder_api.NewErrorInfo(http.StatusInternalServerError).
						WithTitle("Internal Server Error").
						WithDetail("primary error: secondary error").
						Build(),
					).
					Add(*builder_api.NewErrorInfo(http.StatusInternalServerError).
						WithTitle("Internal Server Error").
						WithDetail("secondary error").
						Build(),
					).
					Build(),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			e := helperNewEcho(testCase.Given)
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			e.ServeHTTP(resp, req)
			assert.Equal(t, testCase.Expected.Code, resp.Code)
			currentResponse := api_public.ErrorResponse{}
			err := json.Unmarshal(resp.Body.Bytes(), &currentResponse)
			require.NoError(t, err)
			assert.Equal(t, testCase.Expected.Body, currentResponse)
		})
	}
}
