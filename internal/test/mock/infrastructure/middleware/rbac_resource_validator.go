// Code generated by mockery v2.53.3. DO NOT EDIT.

package middleware

import (
	rbac_data "github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware/rbac-data"
	mock "github.com/stretchr/testify/mock"
)

// RbacResourceValidator is an autogenerated mock type for the RbacResourceValidator type
type RbacResourceValidator struct {
	mock.Mock
}

// Execute provides a mock function with given fields: r
func (_m *RbacResourceValidator) Execute(r rbac_data.RBACResource) rbac_data.RBACResource {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 rbac_data.RBACResource
	if rf, ok := ret.Get(0).(func(rbac_data.RBACResource) rbac_data.RBACResource); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Get(0).(rbac_data.RBACResource)
	}

	return r0
}

// NewRbacResourceValidator creates a new instance of RbacResourceValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRbacResourceValidator(t interface {
	mock.TestingT
	Cleanup(func())
}) *RbacResourceValidator {
	mock := &RbacResourceValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
