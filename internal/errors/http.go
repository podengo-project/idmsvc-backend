package errors

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewHTTPErrorWithInternal creates a new HTTPError instance from a format
// string and an error object.
func NewHTTPErrorWithInternal(internal error, code int, format string, a ...any) error {
	var msg string
	if len(a) != 0 {
		msg = fmt.Sprintf(format, a...)
	} else {
		msg = format
	}
	return &echo.HTTPError{Code: code, Message: msg, Internal: internal}
}

// NewHTTPErrorF creates a new HTTPError instance from a format string.
func NewHTTPErrorF(code int, format string, a ...any) error {
	return NewHTTPErrorWithInternal(nil, code, format, a...)
}

// NilArgError creates a new HTTPError instance with "'name' cannot be nil"
// error message.
func NilArgError(name string) error {
	return NewHTTPErrorF(http.StatusInternalServerError, "'%s' cannot be nil", name)
}
