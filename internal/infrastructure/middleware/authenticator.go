package middleware

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type XRhIValidator interface {
	ValidateXRhIdentity(xrhi *identity.XRHID) error
}

func NewAuthenticator(v XRhIValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

func checkGuardsAuthenticate(v XRhIValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if v == nil {
		return fmt.Errorf("'v' is nil")
	}
	if ctx == nil {
		return fmt.Errorf("'ctx' is nil")
	}
	if input == nil {
		return fmt.Errorf("'input' is nil")
	}
	return nil
}

func Authenticate(v XRhIValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	var (
		err  error
		data *identity.XRHID
	)
	if err = checkGuardsAuthenticate(v, ctx, input); err != nil {
		return err
	}
	if input.SecuritySchemeName != "x-rh-identity" {
		return fmt.Errorf("security scheme '%s' != 'x-rh-identity'", input.SecuritySchemeName)
	}

	// domainCtx, ok := ctx.(DomainContextInterface)
	domainCtx, ok := ctx.Value("oapi-codegen/echo-context").(DomainContextInterface)
	if !ok {
		return fmt.Errorf("'ctx' does not match a 'DomainContextInterface'")
	}
	data = domainCtx.XRHID()
	if err = v.ValidateXRhIdentity(data); err != nil {
		return fmt.Errorf("No valid " + headerXRhIdentity)
	}

	return nil
}
