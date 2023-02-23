// Code generated by mockery v2.16.0. DO NOT EDIT.

package handler

import (
	echo "github.com/labstack/echo/v4"

	mock "github.com/stretchr/testify/mock"

	public "github.com/hmsidm/internal/api/public"
)

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

// CheckHost provides a mock function with given fields: ctx, subscriptionManagerId, fqdn, params
func (_m *Application) CheckHost(ctx echo.Context, subscriptionManagerId string, fqdn string, params public.CheckHostParams) error {
	ret := _m.Called(ctx, subscriptionManagerId, fqdn, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, string, string, public.CheckHostParams) error); ok {
		r0 = rf(ctx, subscriptionManagerId, fqdn, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateDomain provides a mock function with given fields: ctx, params
func (_m *Application) CreateDomain(ctx echo.Context, params public.CreateDomainParams) error {
	ret := _m.Called(ctx, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, public.CreateDomainParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDomain provides a mock function with given fields: ctx, uuid, params
func (_m *Application) DeleteDomain(ctx echo.Context, uuid string, params public.DeleteDomainParams) error {
	ret := _m.Called(ctx, uuid, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, string, public.DeleteDomainParams) error); ok {
		r0 = rf(ctx, uuid, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetLivez provides a mock function with given fields: ctx
func (_m *Application) GetLivez(ctx echo.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMetrics provides a mock function with given fields: ctx
func (_m *Application) GetMetrics(ctx echo.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetReadyz provides a mock function with given fields: ctx
func (_m *Application) GetReadyz(ctx echo.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HostConf provides a mock function with given fields: ctx, fqdn, params
func (_m *Application) HostConf(ctx echo.Context, fqdn string, params public.HostConfParams) error {
	ret := _m.Called(ctx, fqdn, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, string, public.HostConfParams) error); ok {
		r0 = rf(ctx, fqdn, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListDomains provides a mock function with given fields: ctx, params
func (_m *Application) ListDomains(ctx echo.Context, params public.ListDomainsParams) error {
	ret := _m.Called(ctx, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, public.ListDomainsParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadDomain provides a mock function with given fields: ctx, uuid, params
func (_m *Application) ReadDomain(ctx echo.Context, uuid string, params public.ReadDomainParams) error {
	ret := _m.Called(ctx, uuid, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, string, public.ReadDomainParams) error); ok {
		r0 = rf(ctx, uuid, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewApplication interface {
	mock.TestingT
	Cleanup(func())
}

// NewApplication creates a new instance of Application. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewApplication(t mockConstructorTestingTNewApplication) *Application {
	mock := &Application{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
