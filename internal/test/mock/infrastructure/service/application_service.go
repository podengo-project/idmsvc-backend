// Code generated by mockery v2.46.0. DO NOT EDIT.

package service

import mock "github.com/stretchr/testify/mock"

// ApplicationService is an autogenerated mock type for the ApplicationService type
type ApplicationService struct {
	mock.Mock
}

// Start provides a mock function with given fields:
func (_m *ApplicationService) Start() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *ApplicationService) Stop() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Stop")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewApplicationService creates a new instance of ApplicationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplicationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ApplicationService {
	mock := &ApplicationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
