package middleware

import (
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	public_api "github.com/hmsidm/internal/api/public"
	"github.com/labstack/echo/v4"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type xrhiAlwaysTrue struct{}

func (x xrhiAlwaysTrue) ValidateXRhIdentity(xrhi *identity.Identity) error {
	// TODO Implement behavior here
	return nil
}

func NewApiServiceValidator() func(echo.HandlerFunc) echo.HandlerFunc {
	swagger, err := public_api.GetSwagger()
	if err != nil {
		panic(err)
	}
	return middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: NewAuthenticator(xrhiAlwaysTrue{}),
		},
	})
}
