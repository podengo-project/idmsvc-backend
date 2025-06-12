// package echo
//
// This package implements the default error handler to be
// used for the api service to match the openapi specification.
//
// See: https://echo.labstack.com/docs/error-handling#custom-http-error-handler
package echo

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"go.openly.dev/pointy"
)

func DefaultErrorHandler(err error, c echo.Context) {
	logger := app_context.LogFromCtx(c.Request().Context())
	if c.Response().Committed {
		logger.Debug("response already commited",
			slog.String("err", err.Error()),
		)
		return
	}
	errList := make([]api_public.ErrorInfo, 0, 10)

	firstCode := 0
	response := api_public.ErrorResponse{}
	for {
		logger.Debug("parsing error", slog.String("msg", err.Error()))
		errInfo, code := parseError(err)
		errList = append(errList, errInfo)
		if firstCode == 0 {
			firstCode = code
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
	}
	response.Errors = &errList
	if err = c.JSON(firstCode, &response); err != nil {
		logger.Error("failed to send JSON error http response",
			slog.String("msg", err.Error()),
		)
	}
}

// parseError translate an error into api_public.ErrorInfo and
// return the first http status code of the error list.
func parseError(err error) (api_public.ErrorInfo, int) {
	switch value := err.(type) { //nolint:all
	case *echo.HTTPError:
		return api_public.ErrorInfo{
			Status: strconv.Itoa(value.Code),
			Title:  fmt.Sprintf("%v", value.Message),
		}, value.Code
	default:
		return api_public.ErrorInfo{
			Status: strconv.Itoa(http.StatusInternalServerError),
			Title:  http.StatusText(http.StatusInternalServerError),
			Detail: pointy.String(err.Error()),
		}, http.StatusInternalServerError
	}
}
