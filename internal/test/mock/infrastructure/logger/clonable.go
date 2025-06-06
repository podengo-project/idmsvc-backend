// Code generated by mockery v2.53.3. DO NOT EDIT.

package logger

import mock "github.com/stretchr/testify/mock"

// Clonable is an autogenerated mock type for the Clonable type
type Clonable struct {
	mock.Mock
}

// Clone provides a mock function with no fields
func (_m *Clonable) Clone() interface{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Clone")
	}

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// NewClonable creates a new instance of Clonable. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClonable(t interface {
	mock.TestingT
	Cleanup(func())
}) *Clonable {
	mock := &Clonable{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
