package middleware

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"regexp"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	public_api "github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/errors"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
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

// InitOpenAPIFormats configure the admited formats in the openapi
// specification. This function must be called before receive any
// request. Suggested to call before instantiate the middleware.
func InitOpenAPIFormats() {
	// TODO Review all the regular expressions
	openapi3.DefineStringFormat("url", `^https?:\/\/.*$`)
	openapi3.DefineStringFormat("realm", `^(([A-Z0-9][A-Z0-9\-]*[A-Z0-9])|[A-Z0-9]+\.)*([A-Z]+|xn\-\-[A-Z0-9]+)\.?$`)

	// FIXME Search the regular expressions for the below formats
	openapi3.DefineStringFormatCallback("cert-issuer", func(value string) error {
		return checkFormatIssuer(value)
	})
	openapi3.DefineStringFormatCallback("cert-pem", func(value string) error {
		return checkCertificateFormat(value)
	})
	openapi3.DefineStringFormatCallback("cert-subject", func(value string) error {
		return checkFormatSubject(value)
	})
	openapi3.DefineStringFormat("domain-description", `^[\n\x20-\x7E]*$`)
	openapi3.DefineStringFormat("domain-title", `^[a-zA-Z0-9\s]+$`)
	openapi3.DefineStringFormatCallback("ipa-realm-domains", func(value string) error {
		return checkFormatRealmDomains(value)
	})
	openapi3.DefineStringFormat("ipa-server-location", `^[a-zA-Z0-9\s]+$`)
}

func helperCheckRegEx(regex string, fieldName string, fieldValue string) error {
	// https://regex101.com/
	match, err := regexp.MatchString(regex, fieldValue)
	if err != nil {
		return fmt.Errorf("error compiling regular expression: %w", err)
	}
	if !match {
		return fmt.Errorf("'%s'='%s' format not matching", fieldName, fieldValue)
	}
	return nil
}

// checkCertificateFormat check the pem certificate string represented
// by value that can be parsed.
// Return an error or nil for success parsed data.
func checkCertificateFormat(value string) error {
	caCertBlock, _ := pem.Decode([]byte(value))
	if caCertBlock == nil || caCertBlock.Type != "CERTIFICATE" {
		return fmt.Errorf("Failed to decode CA certificate")
	}
	_, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("Failed to parse CA certificate: %w", err)
	}
	return nil
}

// checkFormatIssuer check the subject and issuer format in a certificate
// is valid.
// Return an error or nil for success parsed data.
func checkFormatIssuer(value string) error {
	// https://regex101.com/
	issuerRegEx := `^((?:[A-Z]+=[A-Za-z0-9\.\-\s]+)(?:[ ]*,[ ]*[A-Z]+=[A-Za-z0-9.\-\s]+)*)$`
	return helperCheckRegEx(issuerRegEx, "issuer", value)
}

// checkFormatSubject check the subject and issuer format in a certificate
// is valid.
// Return an error or nil for success parsed data.
func checkFormatSubject(value string) error {
	// https://regex101.com/
	subjectRegEx := `^((?:[A-Z]+=[A-Za-z0-9\.\-\s]+)(?:[ ]*,[ ]*[A-Z]+=[A-Za-z0-9.\-\s]+)*)$`
	return helperCheckRegEx(subjectRegEx, "subject", value)
}

func checkFormatRealmDomains(value string) error {
	// TODO Translate value in a slice, and all the items should validate for a domain
	return nil
}
