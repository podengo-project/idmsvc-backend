// Code generated by mockery v2.46.0. DO NOT EDIT.

package event

import (
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// KafkaHeaders is an autogenerated mock type for the KafkaHeaders type
type KafkaHeaders struct {
	mock.Mock
}

// FromEchoContext provides a mock function with given fields: ctx, _a1
func (_m *KafkaHeaders) FromEchoContext(ctx echo.Context, _a1 string) ([]kafka.Header, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for FromEchoContext")
	}

	var r0 []kafka.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(echo.Context, string) ([]kafka.Header, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(echo.Context, string) []kafka.Header); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]kafka.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(echo.Context, string) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewKafkaHeaders creates a new instance of KafkaHeaders. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKafkaHeaders(t interface {
	mock.TestingT
	Cleanup(func())
}) *KafkaHeaders {
	mock := &KafkaHeaders{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
