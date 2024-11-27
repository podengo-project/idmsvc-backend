// Code generated by mockery v2.46.0. DO NOT EDIT.

package pendo

import (
	context "context"

	pendo "github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
	mock "github.com/stretchr/testify/mock"
)

// Pendo is an autogenerated mock type for the Pendo type
type Pendo struct {
	mock.Mock
}

// SendTrackEvent provides a mock function with given fields: ctx, track
func (_m *Pendo) SendTrackEvent(ctx context.Context, track *pendo.TrackRequest) error {
	ret := _m.Called(ctx, track)

	if len(ret) == 0 {
		panic("no return value specified for SendTrackEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *pendo.TrackRequest) error); ok {
		r0 = rf(ctx, track)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetMetadata provides a mock function with given fields: ctx, kind, group, metrics
func (_m *Pendo) SetMetadata(ctx context.Context, kind pendo.Kind, group pendo.Group, metrics pendo.SetMetadataRequest) (*pendo.SetMetadataResponse, error) {
	ret := _m.Called(ctx, kind, group, metrics)

	if len(ret) == 0 {
		panic("no return value specified for SetMetadata")
	}

	var r0 *pendo.SetMetadataResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pendo.Kind, pendo.Group, pendo.SetMetadataRequest) (*pendo.SetMetadataResponse, error)); ok {
		return rf(ctx, kind, group, metrics)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pendo.Kind, pendo.Group, pendo.SetMetadataRequest) *pendo.SetMetadataResponse); ok {
		r0 = rf(ctx, kind, group, metrics)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pendo.SetMetadataResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, pendo.Kind, pendo.Group, pendo.SetMetadataRequest) error); ok {
		r1 = rf(ctx, kind, group, metrics)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPendo creates a new instance of Pendo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPendo(t interface {
	mock.TestingT
	Cleanup(func())
}) *Pendo {
	mock := &Pendo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
