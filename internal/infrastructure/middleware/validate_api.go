package middleware

import (
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	public_api "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/errors"
	internal_errors "github.com/hmsidm/internal/errors"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

// See: https://github.com/stellirin/go-validator/blob/main/config.go
type xrhiAlwaysTrue struct{}

func (x xrhiAlwaysTrue) ValidateXRhIdentity(xrhi *identity.XRHID) error {
	// TODO Implement behavior here
	// TODO Bear in mind the Identity validation is made at
	//      UserEnforce and SystemEnforce identities
	return nil
}

// NewApiServiceValidator create an API validator middleware.
// Skipper represent the logic to bypass the middleware execution.
// Return the echo middleware or panic on error.
func NewApiServiceValidator(Skipper echo_middleware.Skipper) echo.MiddlewareFunc {
	swagger, err := public_api.GetSwagger()
	if err != nil {
		panic(internal_errors.NewLocationError(err))
	}
	return middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			ExcludeResponseBody:   false,
			ExcludeRequestBody:    false,
			IncludeResponseStatus: true,
			MultiError:            false,
			AuthenticationFunc:    openapi3filter.NoopAuthenticationFunc,
		},
		Skipper: Skipper,
		ErrorHandler: func(c echo.Context, err *echo.HTTPError) error {
			return errors.NewLocationErrorWithLevel(err, 1)
		},
	})
}
