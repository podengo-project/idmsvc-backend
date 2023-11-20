package middleware

import (
	"context"
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock_middleware "github.com/podengo-project/idmsvc-backend/internal/test/mock/infrastructure/middleware"
)

func TestCheckGuardsAuthenticate(t *testing.T) {
	// v is nil
	err := checkGuardsAuthenticate(nil, nil, nil)
	assert.EqualError(t, err, "code=500, message='v' cannot be nil")

	// ctx is nil
	v := mock_middleware.NewXRhIValidator(t)
	err = checkGuardsAuthenticate(v, nil, nil)
	assert.EqualError(t, err, "code=500, message='ctx' cannot be nil")

	// input is nil
	ctx := context.Background()
	err = checkGuardsAuthenticate(v, ctx, nil)
	assert.EqualError(t, err, "code=500, message='input' cannot be nil")

	// Success case
	input := openapi3filter.AuthenticationInput{}
	err = checkGuardsAuthenticate(v, ctx, &input)
	assert.NoError(t, err)
}

func TestNewAuthenticator(t *testing.T) {
	var a openapi3filter.AuthenticationFunc

	// No panics
	assert.NotPanics(t, func() {
		a = NewAuthenticator(nil)
	})
	require.NotNil(t, a)

	// No Panics
	v := mock_middleware.NewXRhIValidator(t)
	assert.NotPanics(t, func() {
		a = NewAuthenticator(v)
	})

	// Wrong security schema name
	err := a(nil, nil)
	assert.EqualError(t, err, "code=500, message='ctx' cannot be nil")

	// Wron
	input := openapi3filter.AuthenticationInput{
		SecuritySchemeName: "no-x-rh-identity",
	}
	ctx := context.WithValue(context.Background(), "oapi-codegen/echo-context", nil)
	err = a(ctx, &input)
	assert.EqualError(t, err, fmt.Sprintf("security scheme '%s' != 'x-rh-identity'", "no-x-rh-identity"))

	// Wrong context type assertion
	ctx = context.WithValue(context.Background(), "oapi-codegen/echo-context", nil)
	input.SecuritySchemeName = "x-rh-identity"
	err = a(ctx, &input)
	assert.EqualError(t, err, "'ctx' does not match a 'DomainContextInterface'")

	// Error on ValidateXRhIdentity
	xrhid := identity.XRHID{}
	mockContext := mock_middleware.NewDomainContextInterface(t)
	mockContext.On("XRHID").Return(&xrhid)
	v.On("ValidateXRhIdentity", &xrhid).Return(fmt.Errorf("any error"))
	ctx = context.WithValue(context.Background(), "oapi-codegen/echo-context", mockContext)
	err = a(ctx, &input)
	assert.EqualError(t, err, "No valid "+headerXRhIdentity)
	mockContext.AssertExpectations(t)
	v.AssertExpectations(t)

	// Error on ValidateXRhIdentity
	mockContext = mock_middleware.NewDomainContextInterface(t)
	mockContext.On("XRHID").Return(&xrhid)
	v = mock_middleware.NewXRhIValidator(t)
	v.On("ValidateXRhIdentity", &xrhid).Return(nil)
	ctx = context.WithValue(context.Background(), "oapi-codegen/echo-context", mockContext)
	a = NewAuthenticator(v)
	err = a(ctx, &input)
	assert.NoError(t, err)
	mockContext.AssertExpectations(t)
	v.AssertExpectations(t)
}
