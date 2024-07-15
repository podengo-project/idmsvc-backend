// Code generated by mockery v2.40.3. DO NOT EDIT.

package event

import (
	event "github.com/podengo-project/idmsvc-backend/internal/infrastructure/event"

	mock "github.com/stretchr/testify/mock"
)

// EventRouter is an autogenerated mock type for the EventRouter type
type EventRouter struct {
	mock.Mock
}

// Add provides a mock function with given fields: topic, _a1
func (_m *EventRouter) Add(topic string, _a1 event.Eventable) {
	_m.Called(topic, _a1)
}

// NewEventRouter creates a new instance of EventRouter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventRouter(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventRouter {
	mock := &EventRouter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
