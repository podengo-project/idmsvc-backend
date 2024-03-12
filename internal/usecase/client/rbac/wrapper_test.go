package rbac

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var (
		client ClientInterface
		err    error
	)

	assert.PanicsWithValue(t, "application is an empty string", func() {
		New("", nil)
	}, "Panic on empty application string")

	assert.PanicsWithValue(t, "rbacClient is nil", func() {
		New("idmsvc", nil)
	}, "Panic on nil rbac client")

	assert.NotPanics(t, func() {
		client, err = NewClientWithResponses("http://localhost:8000/api/rbac/v1")
		New("idmsvc", client)
	}, "Instantiate the rbac wrapper")
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestMatchPermission(t *testing.T) {
	c := &rbacWrapper{
		application: "idmsvc",
		client:      nil,
	}
	assert.True(t, c.matchPermission("idmsvc", "domains", "create", "idmsvc:domains:create"))
	assert.False(t, c.matchPermission("idmsvc", "domains", "read", "idmsvc:domains:create"))

	assert.True(t, c.matchPermission("idmsvc", "domains", wildcard, "idmsvc:domains:create"))
	assert.True(t, c.matchPermission("idmsvc", wildcard, wildcard, "idmsvc:domains:create"))

	assert.False(t, c.matchPermission("idmsvc", "domains", wildcard, "rbac:permission:read"))
	assert.False(t, c.matchPermission("idmsvc", wildcard, wildcard, "rbac:permission:read"))

	assert.True(t, c.matchPermission(wildcard, wildcard, wildcard, "idmsvc:domains:create"))
}

func TestMatchPermissionLabel(t *testing.T) {
	c := &rbacWrapper{
		application: "idmsvc",
		client:      nil,
	}
	assert.True(t, c.matchPermissionLabel("read", "read"))
	assert.False(t, c.matchPermissionLabel("write", "read"))
	assert.True(t, c.matchPermissionLabel(wildcard, "read"))
}

func TestDecomposePermission(t *testing.T) {
	c := &rbacWrapper{
		application: "idmsvc",
		client:      nil,
	}

	assert.PanicsWithValue(t, "wrong permission tuple", func() {
		c.decomposePermission("idmsvc")
	}, "Panic on tuple different to 3 items")

	assert.NotPanics(t, func() {
		c.decomposePermission("idmsvc:domains:read")
	}, "No panic on tuple of 3 items")
}

func TestAddXRHID(t *testing.T) {
	var err error
	c := &rbacWrapper{
		application: "idmsvc",
		client:      nil,
	}

	xrhid := builder_api.NewUserXRHID().Build()
	xrhidRaw := header.EncodeXRHID(&xrhid)
	assert.NotPanics(t, func() {
		ctx := ContextWithXRHID(context.Background(), xrhidRaw)
		req := httptest.NewRequest(http.MethodGet, "/domains", http.NoBody)
		err = c.addXRHID(ctx, req)
	}, "No panic on normal operation")
	require.NoError(t, err)
}
