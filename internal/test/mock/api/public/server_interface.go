// Code generated by mockery v2.46.0. DO NOT EDIT.

package public

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"

	public "github.com/podengo-project/idmsvc-backend/internal/api/public"

	uuid "github.com/google/uuid"
)

// ServerInterface is an autogenerated mock type for the ServerInterface type
type ServerInterface struct {
	mock.Mock
}

// CreateDomainToken provides a mock function with given fields: ctx, params
func (_m *ServerInterface) CreateDomainToken(ctx echo.Context, params public.CreateDomainTokenParams) error {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for CreateDomainToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, public.CreateDomainTokenParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDomain provides a mock function with given fields: ctx, _a1, params
func (_m *ServerInterface) DeleteDomain(ctx echo.Context, _a1 uuid.UUID, params public.DeleteDomainParams) error {
	ret := _m.Called(ctx, _a1, params)

	if len(ret) == 0 {
		panic("no return value specified for DeleteDomain")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID, public.DeleteDomainParams) error); ok {
		r0 = rf(ctx, _a1, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSigningKeys provides a mock function with given fields: ctx, params
func (_m *ServerInterface) GetSigningKeys(ctx echo.Context, params public.GetSigningKeysParams) error {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetSigningKeys")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, public.GetSigningKeysParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HostConf provides a mock function with given fields: ctx, inventoryId, fqdn, params
func (_m *ServerInterface) HostConf(ctx echo.Context, inventoryId uuid.UUID, fqdn string, params public.HostConfParams) error {
	ret := _m.Called(ctx, inventoryId, fqdn, params)

	if len(ret) == 0 {
		panic("no return value specified for HostConf")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID, string, public.HostConfParams) error); ok {
		r0 = rf(ctx, inventoryId, fqdn, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListDomains provides a mock function with given fields: ctx, params
func (_m *ServerInterface) ListDomains(ctx echo.Context, params public.ListDomainsParams) error {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for ListDomains")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, public.ListDomainsParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadDomain provides a mock function with given fields: ctx, _a1, params
func (_m *ServerInterface) ReadDomain(ctx echo.Context, _a1 uuid.UUID, params public.ReadDomainParams) error {
	ret := _m.Called(ctx, _a1, params)

	if len(ret) == 0 {
		panic("no return value specified for ReadDomain")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID, public.ReadDomainParams) error); ok {
		r0 = rf(ctx, _a1, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterDomain provides a mock function with given fields: ctx, params
func (_m *ServerInterface) RegisterDomain(ctx echo.Context, params public.RegisterDomainParams) error {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for RegisterDomain")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, public.RegisterDomainParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateDomainAgent provides a mock function with given fields: ctx, _a1, params
func (_m *ServerInterface) UpdateDomainAgent(ctx echo.Context, _a1 uuid.UUID, params public.UpdateDomainAgentParams) error {
	ret := _m.Called(ctx, _a1, params)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDomainAgent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID, public.UpdateDomainAgentParams) error); ok {
		r0 = rf(ctx, _a1, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateDomainUser provides a mock function with given fields: ctx, _a1, params
func (_m *ServerInterface) UpdateDomainUser(ctx echo.Context, _a1 uuid.UUID, params public.UpdateDomainUserParams) error {
	ret := _m.Called(ctx, _a1, params)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDomainUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context, uuid.UUID, public.UpdateDomainUserParams) error); ok {
		r0 = rf(ctx, _a1, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewServerInterface creates a new instance of ServerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServerInterface {
	mock := &ServerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
