package rbac_data

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRbacServiceValidator(t *testing.T) {
	const service = RBACService("myservice")
	out := NewRbacServiceValidator(service)

	assert.PanicsWithValue(t, "service is undefined", func() {
		out(RBACServiceUndefined)
	}, "Panic on undefined service")

	assert.PanicsWithValue(t, "Service must be 'myservice'", func() {
		out(RBACService("anotherservice"))
	}, "Panic on wrong service")

	var result RBACService
	assert.NotPanics(t, func() {
		result = out(service)
	})
	assert.Equal(t, service, result, "validate service")
}

func TestDefaultRbacServiceValidator(t *testing.T) {
	assert.PanicsWithValue(t, "service is undefined", func() {
		DefaultRbacServiceValidator("")
	}, "Panic on undefined service")

	assert.NotPanics(t, func() {
		DefaultRbacServiceValidator("someservice")
	}, "Success scenario")
}

func TestDefaultRbacResourceValidator(t *testing.T) {
	assert.PanicsWithValue(t, "resource is undefined", func() {
		DefaultRbacResourceValidator(RBACResourceUndefined)
	}, "Panics on undefined resource")

	assert.NotPanics(t, func() {
		DefaultRbacResourceValidator("myresource")
	}, "Success scenario")
}

func TestSetRbacServiceValidator(t *testing.T) {
	assert.PanicsWithValue(t, "service validator is nil", func() {
		SetRbacServiceValidator(nil)
	}, "Panics on nil validator")

	assert.NotPanics(t, func() {
		SetRbacServiceValidator(DefaultRbacServiceValidator)
	}, "Success scenario")
}

func TestSetRbacResourceValidator(t *testing.T) {
	assert.PanicsWithValue(t, "resource validator is nil", func() {
		SetRbacResourceValidator(nil)
	}, "Panics on nil validator")

	assert.NotPanics(t, func() {
		SetRbacResourceValidator(DefaultRbacResourceValidator)
	}, "Succesful scenario")
}

func TestDefaultRbacVerbValidator(t *testing.T) {
	assert.PanicsWithValue(t, "Verb is undefined", func() {
		DefaultRbacVerbValidator(RbacVerbUndefined)
	}, "Panics on undefined verb")

	assert.PanicsWithValue(t, "Verb 'unsupported' is not supported", func() {
		DefaultRbacVerbValidator("unsupported")
	}, "Panics on undefined verb")

	validverbs := []RBACVerb{
		RbacVerbAny,
		RbacVerbRead,
		RbacVerbWrite,
		RbacVerbCreate,
		RbacVerbUpdate,
		RbacVerbDelete,
		RbacVerbLink,
		RbacVerbUnlink,
		RbacVerbOrder,
		RbacVerbExecute,
	}
	for _, verb := range validverbs {
		var result RBACVerb
		assert.NotPanics(t, func() {
			result = DefaultRbacVerbValidator(verb)
		}, "Panics on undefined verb")
		assert.Equal(t, verb, result)
	}
}

func TestNewRbacResourceValidator(t *testing.T) {
	var (
		result        RBACResource
		resource      = RBACResource("test")
		otherresource = RBACResource("otherresource")
	)
	validator := NewRbacResourceValidator()
	require.NotNil(t, validator)

	validator = NewRbacResourceValidator(resource)
	require.NotNil(t, validator)
	assert.PanicsWithValue(t, fmt.Sprintf("resource '%s' is not valid", otherresource), func() {
		_ = validator(otherresource)
	}, "Panic on otherresource")

	validator = NewRbacResourceValidator(resource)
	require.NotNil(t, validator)
	assert.PanicsWithValue(t, fmt.Sprintf("resource '%s' is not valid", RBACResourceUndefined), func() {
		result = validator(RBACResource(RBACResourceUndefined))
	}, "Panic on Undefined resource")

	validator = NewRbacResourceValidator(resource)
	require.NotNil(t, validator)
	assert.PanicsWithValue(t, fmt.Sprintf("resource '%s' is not valid", RBACResourceAny), func() {
		result = validator(RBACResource(RBACResourceAny))
	}, "Panic on * any resource is indicated")

	validator = NewRbacResourceValidator(resource)
	require.NotNil(t, validator)
	assert.NotPanics(t, func() {
		result = validator(RBACResource(resource))
	}, "Success scenario")
	assert.Equal(t, RBACResource(resource), result)
}

func helperSetValidators(service RBACService, resource RBACResource) {
	SetRbacServiceValidator(NewRbacServiceValidator(service))
	SetRbacResourceValidator(NewRbacResourceValidator(resource))
}

func TestNewRbacPermission(t *testing.T) {
	service := RBACService("myservice")
	resource := RBACResource("myresource")
	helperSetValidators(service, resource)
	expected := RBACPermission(
		fmt.Sprintf(
			"%s:%s:%s",
			string(service),
			string(resource),
			string(RbacVerbRead),
		),
	)
	result := NewRbacPermission(service, resource, RbacVerbRead)
	assert.Equal(t, expected, result)
}

func TestGetPermissionGuards(t *testing.T) {
	var err error
	service := RBACService("idmsvc")
	resource := RBACResource("domains")
	helperSetValidators(service, resource)
	rbacMap := NewRBACMapBuilder().
		Add("/domains", http.MethodGet, NewRbacPermission(service, resource, RbacVerbRead)).
		Build()
	assert.NotEqual(t, RBACMap{}, rbacMap)

	err = rbacMap.getPermissionGuards("", "", "")
	require.EqualError(t, err, "prefix is an empty string when it was expected for instance '/api/someservice/v1'")

	err = rbacMap.getPermissionGuards("/api/idmsvc/v1", "", "")
	require.EqualError(t, err, "route is an empty string")

	err = rbacMap.getPermissionGuards("/api/idmsvc/v1", "/api/idmsvc/v1/domains", "")
	require.EqualError(t, err, "method is an empty string")

	err = rbacMap.getPermissionGuards("/api/idmsvc/v1", "/api/idmsvc/v1/domains", http.MethodGet)
	require.NoError(t, err, "Success scenario")
}

func TestGetPermission(t *testing.T) {
	service := RBACService("idmsvc")
	resource := RBACResource("domains")
	var err error

	helperSetValidators(service, resource)
	rbacMap := NewRBACMapBuilder().
		Add("/domains", http.MethodGet, NewRbacPermission(service, resource, RbacVerbRead)).
		Build()

	var permission RBACPermission
	permission, err = rbacMap.GetPermission("", "", "")
	require.EqualError(t, err, "code=401, message=prefix is an empty string when it was expected for instance '/api/someservice/v1'")
	assert.Equal(t, RBACPermissionUndefined, permission)

	permission, err = rbacMap.GetPermission("/api/idmsvc/v1", "", http.MethodGet)
	require.EqualError(t, err, "code=401, message=route is an empty string")
	assert.Equal(t, RBACPermissionUndefined, permission)

	permission, err = rbacMap.GetPermission("/api/idmsvc/v1", "/api/idmsvc/v1/domains/token", http.MethodGet)
	require.EqualError(t, err, "code=401, message=Route not authorized")
	assert.Equal(t, RBACPermissionUndefined, permission)

	permission, err = rbacMap.GetPermission("/api/idmsvc/v1", "/api/idmsvc/v1/domains", http.MethodPost)
	require.EqualError(t, err, "code=401, message=Method not authorized")
	assert.Equal(t, RBACPermissionUndefined, permission)

	permission, err = rbacMap.GetPermission("/api/idmsvc/v1", "/api/idmsvc/v1/domains", http.MethodGet)
	require.NoError(t, err)
	assert.Equal(t, RBACPermission("idmsvc:domains:read"), permission)
}
