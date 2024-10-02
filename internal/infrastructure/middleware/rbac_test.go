package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	rbac_data "github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware/rbac-data"
	api_builder "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	client_rbac "github.com/podengo-project/idmsvc-backend/internal/test/mock/interface/client/rbac"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func helperRbacSetupEcho(config *RBACConfig, method, path string, status int) *echo.Echo {
	e := echo.New()
	e.Use(ContextLogConfig(&LogConfig{}))
	e.Use(CreateContext())
	e.Use(ParseXRHIDMiddlewareWithConfig(&ParseXRHIDMiddlewareConfig{}))
	e.Use(EnforceIdentityWithConfig(&IdentityConfig{}))
	e.Use(RBACWithConfig(config))
	e.Add(method, path, func(c echo.Context) error {
		return c.NoContent(status)
	})
	return e
}

func helperRbacCreateMapping() rbac_data.RBACMap {
	service := rbac_data.RBACService("idmsvc")
	resourceToken := rbac_data.RBACResource("token")
	resourceDomains := rbac_data.RBACResource("domains")
	// This map fit the one at internal/infrastructure/router/rbac.yaml
	data := rbac_data.NewRBACMapBuilder().
		Add("/domains/token", http.MethodPost, rbac_data.NewRbacPermission(service, resourceToken, rbac_data.RbacVerbCreate)).
		Add("/domains", http.MethodGet, rbac_data.NewRbacPermission(service, resourceDomains, rbac_data.RbacVerbRead)).
		Add("/domains/:uuid", http.MethodGet, rbac_data.NewRbacPermission(service, resourceDomains, rbac_data.RbacVerbRead)).
		Add("/domains/:uuid", http.MethodPatch, rbac_data.NewRbacPermission(service, resourceDomains, rbac_data.RbacVerbUpdate)).
		Add("/domains/:uuid", http.MethodDelete, rbac_data.NewRbacPermission(service, resourceDomains, rbac_data.RbacVerbDelete)).
		Build()
	return data
}

func helperRbacSetup(t *testing.T, prefix string, skipper echo_middleware.Skipper) (RBACConfig, identity.XRHID, bool) {
	rbacMap := helperRbacCreateMapping()
	rbacClientMock := client_rbac.NewRbac(t)
	rbacConfig := RBACConfig{
		Skipper:       skipper,
		Prefix:        prefix,
		PermissionMap: rbacMap,
		Client:        rbacClientMock,
	}
	xrhid := api_builder.NewUserXRHID().
		WithOrgID("12345").
		WithUserID("test").
		Build()
	return rbacConfig, xrhid, false
}

func TestRBACWithConfigErrors(t *testing.T) {
	prefix := "/api/idmsvc/v1"

	assert.Panics(t, func() {
		RBACWithConfig(nil)
	}, "Fails with nil configuration")

	assert.Panics(t, func() {
		RBACWithConfig(&RBACConfig{
			Prefix: "",
		})
	}, "Fails with empty Prefix")

	assert.Panics(t, func() {
		RBACWithConfig(&RBACConfig{
			Prefix: prefix,
			Client: nil,
		})
	}, "Fails with nil client")

	assert.Panics(t, func() {
		RBACWithConfig(&RBACConfig{
			Prefix: "",
		})
	}, "Fails with nil configuration")
}

func helperRbacSkipper(c echo.Context) bool {
	return c.Path() == "/api/idmsvc/v1/openapi.json"
}

func TestRBACWithConfigSkipper(t *testing.T) {
	prefix := "/api/idmsvc/v1"
	rbacConfig, xrhid, shouldReturn := helperRbacSetup(t, prefix, helperRbacSkipper)
	if shouldReturn {
		return
	}
	method := http.MethodGet
	path := "/api/idmsvc/v1/openapi.json"
	statusExpected := http.StatusOK
	e := helperRbacSetupEcho(&rbacConfig, method, path, statusExpected)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, "http://localhost:8000"+path, http.NoBody)
	req.Header.Add(header.HeaderXRHID, header.EncodeXRHID(&xrhid))
	e.ServeHTTP(rec, req)
	assert.Equal(t, statusExpected, rec.Code)
}

func TestRBACWithConfigFailGetPermission(t *testing.T) {
	prefix := "/api/idmsvc/v1"
	rbacMap := helperRbacCreateMapping()
	rbacClientMock := client_rbac.NewRbac(t)
	rbacConfig := RBACConfig{
		Prefix:        prefix,
		PermissionMap: rbacMap,
		Client:        rbacClientMock,
	}
	xrhid := api_builder.NewUserXRHID().
		WithOrgID("12345").
		WithUserID("test").
		Build()
	method := http.MethodGet
	path := "/api/idmsvc/v1/openapi.json"
	statusExpected := http.StatusUnauthorized
	e := helperRbacSetupEcho(&rbacConfig, method, path, statusExpected)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, "http://localhost:8000"+path, http.NoBody)
	req.Header.Add(header.HeaderXRHID, header.EncodeXRHID(&xrhid))
	e.ServeHTTP(rec, req)
	assert.Equal(t, statusExpected, rec.Code)
}

func TestRBACWithConfigFailRbacClient(t *testing.T) {
	prefix := "/api/idmsvc/v1"
	rbacMap := helperRbacCreateMapping()
	rbacClientMock := client_rbac.NewRbac(t)
	rbacConfig := RBACConfig{
		Prefix:        prefix,
		PermissionMap: rbacMap,
		Client:        rbacClientMock,
	}
	xrhid := api_builder.NewUserXRHID().
		WithOrgID("12345").
		WithUserID("test").
		Build()
	method := http.MethodGet
	path := prefix + "/domains"
	statusExpected := http.StatusNotFound
	e := helperRbacSetupEcho(&rbacConfig, method, path, http.StatusOK)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, "http://localhost:8000"+path, http.NoBody)
	req.Header.Add(header.HeaderXRHID, header.EncodeXRHID(&xrhid))
	permission := rbac_data.NewRbacPermission("idmsvc", "domains", rbac_data.RbacVerbRead)
	rbacClientMock.On("IsAllowed", mock.Anything, string(header.EncodeXRHID(&xrhid)), string(permission)).Return(false, echo.NewHTTPError(http.StatusNotFound))
	e.ServeHTTP(rec, req)
	assert.Equal(t, statusExpected, rec.Code)
}

func TestRBACWithConfigUnauthorized(t *testing.T) {
	prefix := "/api/idmsvc/v1"
	rbacMap := helperRbacCreateMapping()
	rbacClientMock := client_rbac.NewRbac(t)
	rbacConfig := RBACConfig{
		Prefix:        prefix,
		PermissionMap: rbacMap,
		Client:        rbacClientMock,
	}
	xrhid := api_builder.NewUserXRHID().
		WithOrgID("12345").
		WithUserID("test").
		Build()
	method := http.MethodGet
	path := prefix + "/domains"
	statusExpected := http.StatusUnauthorized
	e := helperRbacSetupEcho(&rbacConfig, method, path, http.StatusOK)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, "http://localhost:8000"+path, http.NoBody)
	req.Header.Add(header.HeaderXRHID, header.EncodeXRHID(&xrhid))
	permission := rbac_data.NewRbacPermission("idmsvc", "domains", rbac_data.RbacVerbRead)
	rbacClientMock.On("IsAllowed", mock.Anything, string(header.EncodeXRHID(&xrhid)), string(permission)).Return(false, nil)
	e.ServeHTTP(rec, req)
	assert.Equal(t, statusExpected, rec.Code)
}

func TestRBACWithConfigGrantedPermission(t *testing.T) {
	prefix := "/api/idmsvc/v1"
	rbacMap := helperRbacCreateMapping()
	rbacClientMock := client_rbac.NewRbac(t)
	rbacConfig := RBACConfig{
		Prefix:        prefix,
		PermissionMap: rbacMap,
		Client:        rbacClientMock,
	}
	xrhid := api_builder.NewUserXRHID().
		WithOrgID("12345").
		WithUserID("test").
		Build()
	method := http.MethodGet
	path := prefix + "/domains"
	statusExpected := http.StatusOK
	e := helperRbacSetupEcho(&rbacConfig, method, path, http.StatusOK)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, "http://localhost:8000"+path, http.NoBody)
	req.Header.Add(header.HeaderXRHID, header.EncodeXRHID(&xrhid))
	permission := rbac_data.NewRbacPermission("idmsvc", "domains", rbac_data.RbacVerbRead)
	rbacClientMock.On("IsAllowed", mock.Anything, string(header.EncodeXRHID(&xrhid)), string(permission)).Return(true, nil)
	e.ServeHTTP(rec, req)
	assert.Equal(t, statusExpected, rec.Code)
}

func TestRBACWithConfigFailGettingUserPermissions(t *testing.T) {
	prefix := "/api/idmsvc/v1"
	rbacMap := helperRbacCreateMapping()
	rbacClientMock := client_rbac.NewRbac(t)
	rbacConfig := RBACConfig{
		Prefix:        prefix,
		PermissionMap: rbacMap,
		Client:        rbacClientMock,
	}
	xrhid := api_builder.NewUserXRHID().
		WithOrgID("12345").
		WithUserID("test").
		Build()
	method := http.MethodGet
	path := "/api/idmsvc/v1/openapi.json"
	statusExpected := http.StatusUnauthorized
	e := helperRbacSetupEcho(&rbacConfig, method, path, statusExpected)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, "http://localhost:8000"+path, http.NoBody)
	req.Header.Add(header.HeaderXRHID, header.EncodeXRHID(&xrhid))
	e.ServeHTTP(rec, req)
	assert.Equal(t, statusExpected, rec.Code)
}
