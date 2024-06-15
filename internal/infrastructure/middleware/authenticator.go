package middleware

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
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
		return internal_errors.NilArgError("v")
	}
	if ctx == nil {
		return internal_errors.NilArgError("ctx")
	}
	if input == nil {
		return internal_errors.NilArgError("input")
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
		return fmt.Errorf("No valid " + header.HeaderXRHID)
	}

	return nil
}
