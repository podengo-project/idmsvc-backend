package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/stretchr/testify/assert"
)

func helperSetupRBACMiddleware(config *RBACConfig, method, path string) *echo.Echo {
	e := echo.New()
	e.Use(CreateContext())
	e.Use(EnforceIdentityWithConfig(&IdentityConfig{}))
	e.Use(RBACWithConfig(config))
	e.Add(method, path, func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	return e
}

func TestRBACWithConfig(t *testing.T) {
	xrhid := api.NewUserXRHID().Build()
	xrhidString := header.EncodeXRHID(&xrhid)
	rbacMap := RBACMap{}
	rbacConfig := RBACConfig{Skipper: nil, PermissionMap: rbacMap}

	// Fails with nil configuration
	assert.Panics(t, func() {
		RBACWithConfig(nil)
	})

	// Execute middleware
	method := http.MethodGet
	path := "/api/idmsvc/v1/domains"
	e := helperSetupRBACMiddleware(&rbacConfig, method, path)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, http.NoBody)
	req.Header.Add("X-Rh-Identity", xrhidString)
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// TODO HMS-3522
}
