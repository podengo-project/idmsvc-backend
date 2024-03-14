// Code generated by mockery v2.38.0. DO NOT EDIT.

package middleware

import (
	rbac_data "github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware/rbac-data"
	mock "github.com/stretchr/testify/mock"
)

// RbacServiceValidator is an autogenerated mock type for the RbacServiceValidator type
type RbacServiceValidator struct {
	mock.Mock
}

// Execute provides a mock function with given fields: s
func (_m *RbacServiceValidator) Execute(s rbac_data.RBACService) rbac_data.RBACService {
	ret := _m.Called(s)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 rbac_data.RBACService
	if rf, ok := ret.Get(0).(func(rbac_data.RBACService) rbac_data.RBACService); ok {
		r0 = rf(s)
	} else {
		r0 = ret.Get(0).(rbac_data.RBACService)
	}

	return r0
}

// NewRbacServiceValidator creates a new instance of RbacServiceValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRbacServiceValidator(t interface {
	mock.TestingT
	Cleanup(func())
}) *RbacServiceValidator {
	mock := &RbacServiceValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
