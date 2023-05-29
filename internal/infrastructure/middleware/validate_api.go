package middleware

import (
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	public_api "github.com/hmsidm/internal/api/public"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type xrhiAlwaysTrue struct{}

func (x xrhiAlwaysTrue) ValidateXRhIdentity(xrhi *identity.XRHID) error {
	// TODO Implement behavior here
	return nil
}

func NewApiServiceValidator(Skipper echo_middleware.Skipper) echo.MiddlewareFunc {
	swagger, err := public_api.GetSwagger()
	if err != nil {
		panic(err)
	}
	return middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			ExcludeResponseBody:   false,
			ExcludeRequestBody:    false,
			IncludeResponseStatus: true,
			MultiError:            false,
			AuthenticationFunc:    NewAuthenticator(xrhiAlwaysTrue{}),
		},
		Skipper: Skipper,
	})
}
