// Code generated by mockery v2.40.3. DO NOT EDIT.

package service

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// MockRbac is an autogenerated mock type for the MockRbac type
type MockRbac struct {
	mock.Mock
}

// GetBaseURL provides a mock function with given fields:
func (_m *MockRbac) GetBaseURL() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetBaseURL")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetPermissions provides a mock function with given fields: data
func (_m *MockRbac) SetPermissions(data []string) {
	_m.Called(data)
}

// WaitAddress provides a mock function with given fields: timeout
func (_m *MockRbac) WaitAddress(timeout time.Duration) error {
	ret := _m.Called(timeout)

	if len(ret) == 0 {
		panic("no return value specified for WaitAddress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Duration) error); ok {
		r0 = rf(timeout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockRbac creates a new instance of MockRbac. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRbac(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRbac {
	mock := &MockRbac{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
