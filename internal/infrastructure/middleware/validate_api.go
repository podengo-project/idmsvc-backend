package middleware

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	public_api "github.com/podengo-project/idmsvc-backend/internal/api/public"
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
	openapi3.DefineStringFormatCallback("domain-description", checkUtf8MultiLine)
	openapi3.DefineStringFormatCallback("domain-title", checkUtf8SingleLine)
	openapi3.DefineStringFormatCallback("ipa-realm-domains", func(value string) error {
		return checkFormatRealmDomains(value)
	})
	openapi3.DefineStringFormat("ipa-server-location", `^[a-zA-Z0-9\s]+$`)
	openapi3.DefineStringFormat("pagination-ref", `^(\/\w+){4}\?(offset|limit)=(0|[1-9]\d*)\&(offset|limit)=(0|[1-9]\d*)$`)
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

// Check that input is a valid UTF-8 string with no control characters
// (not even CR/LF).
func checkUtf8SingleLine(s string) error {
	if !utf8.ValidString(s) {
		return fmt.Errorf("not a valid utf-8 string")
	}
	for _ /* index */, r /* rune */ := range s {
		if unicode.IsControl(r) {
			return fmt.Errorf("invalid code point: %U", r)
		}
	}
	return nil
}

// Check that input is a valid UTF-8 string with no control characters,
// except spacing chars including CR/LF are allowed.
func checkUtf8MultiLine(s string) error {
	if !utf8.ValidString(s) {
		return fmt.Errorf("not a valid utf-8 string")
	}
	for _ /* index */, r /* rune */ := range s {
		if unicode.IsControl(r) && !unicode.IsSpace(r) {
			return fmt.Errorf("invalid code point: %U", r)
		}
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

// In order to validate a response, we need to have access to the bytes of
// the response. The following code allows us to get access to it.
type ResponseRecorder struct {
	buffer   *bytes.Buffer
	status   int
	original http.ResponseWriter
}

// Implements WriteHeader of http.ResponseWriter
func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
}

// Implements Header of http.ResponseWriter
func (r *ResponseRecorder) Header() http.Header {
	return r.original.Header()
}

// Implements Write of http.ResponseWriter
func (r *ResponseRecorder) Write(p []byte) (n int, err error) {
	n, err = r.buffer.Write(p)
	if err != nil {
		return n, err
	}
	return len(p), nil
}

type (
	// RequestResponseValidatorConfig defines the config for RequestResponseValidator middleware.
	RequestResponseValidatorConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper          echo_middleware.Skipper
		ValidateRequest  bool
		ValidateResponse bool
	}
)

// DefaultRequestResponseValidatorConfig is the default RequestResponseValidator
// middleware config.
var DefaultRequestResponseValidatorConfig = RequestResponseValidatorConfig{
	Skipper:          echo_middleware.DefaultSkipper,
	ValidateRequest:  true,
	ValidateResponse: false,
}

// RequestResponseValidator returns a middleware which validates the HTTP response
func RequestResponseValidator() echo.MiddlewareFunc {
	return RequestResponseValidatorWithConfig(&DefaultRequestResponseValidatorConfig)
}

func RequestResponseValidatorWithConfig(config *RequestResponseValidatorConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestResponseValidatorConfig.Skipper
	}

	// Setup routing
	swagger, err := public_api.GetSwagger()
	if err != nil {
		panic(internal_errors.NewLocationError(err))
	}

	// TODO Can we use our own router?  Something to investigate later.
	router, err := gorillamux.NewRouter(swagger)
	if err != nil {
		panic(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// Nothing to validate, just leave.
			if !config.ValidateRequest && !config.ValidateResponse {
				return next(c)
			}

			req := c.Request()
			route, pathParams, err := router.FindRoute(req)
			if err != nil {
				// No route found in the OpenAPI for this request.
				// Returning 404 (or other error status) here is
				// possible but WRONG because:
				//
				// - There are some special case routes not represented
				//   in the OpenAPI spec, e.g. /api/v1/openapi.json.
				//   If we response 404 (or 5xx), we would need to define
				//   a skipper to skip these special cases.
				//
				// - Is it not the concern of this middleware to route
				//   the request (or response 404 if no route found).
				//   That behaviour already exists in the Echo framework.
				//   This middleware need only concern itself about whether
				//   to validate the request.
				//
				// So, if we don't find a round in the OpenAPI spec, just
				// call the next middleware/handler.
				//
				return next(c)
			}

			options := openapi3filter.Options{
				ExcludeResponseBody:   false,
				ExcludeRequestBody:    false,
				IncludeResponseStatus: true,
				MultiError:            false,
				AuthenticationFunc:    openapi3filter.NoopAuthenticationFunc,
			}

			requestValidationInput := &openapi3filter.RequestValidationInput{
				Request:    req,
				PathParams: pathParams,
				Route:      route,
				Options:    &options,
			}

			ctx := c.Request().Context()
			if config.ValidateRequest {
				err = openapi3filter.ValidateRequest(
					ctx,
					requestValidationInput,
				)
				if err != nil {
					c.Response().Header().Set(echo.HeaderContentType, "text/plain")
					c.String(http.StatusBadRequest, err.Error())
					return nil // stop processing
				}
			}

			if config.ValidateResponse {
				// Intercept and validate the response
				rw := c.Response().Writer
				resRec := &ResponseRecorder{buffer: &bytes.Buffer{}, original: rw}
				c.Response().Writer = resRec

				defer func() {
					responseValidationInput := &openapi3filter.ResponseValidationInput{
						RequestValidationInput: requestValidationInput,
						Status:                 resRec.status,
						Header:                 resRec.Header(),
					}
					responseValidationInput.SetBodyBytes(resRec.buffer.Bytes())

					if err = openapi3filter.ValidateResponse(
						ctx,
						responseValidationInput,
					); err != nil {
						// write error response
						c.Response().Header().Set(echo.HeaderContentType, "text/plain")
						c.String(http.StatusInternalServerError, err.Error())
					} else {
						// Write original response
						rw.WriteHeader(resRec.status)
						resRec.buffer.WriteTo(rw)
					}
				}()

				defer func() {
					// reset the response, using the original ResponseWriter
					c.SetResponse(echo.NewResponse(rw, c.Echo()))
				}()
			}

			return next(c)
		}
	}
}
