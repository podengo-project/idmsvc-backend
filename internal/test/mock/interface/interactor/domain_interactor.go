// Code generated by mockery v2.38.0. DO NOT EDIT.

package interactor

import (
	header "github.com/podengo-project/idmsvc-backend/internal/api/header"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"

	mock "github.com/stretchr/testify/mock"

	model "github.com/podengo-project/idmsvc-backend/internal/domain/model"

	public "github.com/podengo-project/idmsvc-backend/internal/api/public"

	uuid "github.com/google/uuid"
)

// DomainInteractor is an autogenerated mock type for the DomainInteractor type
type DomainInteractor struct {
	mock.Mock
}

// CreateDomainToken provides a mock function with given fields: xrhid, params, body
func (_m *DomainInteractor) CreateDomainToken(xrhid *identity.XRHID, params *public.CreateDomainTokenParams, body *public.DomainRegTokenRequest) (string, public.DomainType, error) {
	ret := _m.Called(xrhid, params, body)

	if len(ret) == 0 {
		panic("no return value specified for CreateDomainToken")
	}

	var r0 string
	var r1 public.DomainType
	var r2 error
	if rf, ok := ret.Get(0).(func(*identity.XRHID, *public.CreateDomainTokenParams, *public.DomainRegTokenRequest) (string, public.DomainType, error)); ok {
		return rf(xrhid, params, body)
	}
	if rf, ok := ret.Get(0).(func(*identity.XRHID, *public.CreateDomainTokenParams, *public.DomainRegTokenRequest) string); ok {
		r0 = rf(xrhid, params, body)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*identity.XRHID, *public.CreateDomainTokenParams, *public.DomainRegTokenRequest) public.DomainType); ok {
		r1 = rf(xrhid, params, body)
	} else {
		r1 = ret.Get(1).(public.DomainType)
	}

	if rf, ok := ret.Get(2).(func(*identity.XRHID, *public.CreateDomainTokenParams, *public.DomainRegTokenRequest) error); ok {
		r2 = rf(xrhid, params, body)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Delete provides a mock function with given fields: xrhid, UUID, params
func (_m *DomainInteractor) Delete(xrhid *identity.XRHID, UUID uuid.UUID, params *public.DeleteDomainParams) (string, uuid.UUID, error) {
	ret := _m.Called(xrhid, UUID, params)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 string
	var r1 uuid.UUID
	var r2 error
	if rf, ok := ret.Get(0).(func(*identity.XRHID, uuid.UUID, *public.DeleteDomainParams) (string, uuid.UUID, error)); ok {
		return rf(xrhid, UUID, params)
	}
	if rf, ok := ret.Get(0).(func(*identity.XRHID, uuid.UUID, *public.DeleteDomainParams) string); ok {
		r0 = rf(xrhid, UUID, params)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*identity.XRHID, uuid.UUID, *public.DeleteDomainParams) uuid.UUID); ok {
		r1 = rf(xrhid, UUID, params)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(2).(func(*identity.XRHID, uuid.UUID, *public.DeleteDomainParams) error); ok {
		r2 = rf(xrhid, UUID, params)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetByID provides a mock function with given fields: xrhid, params
func (_m *DomainInteractor) GetByID(xrhid *identity.XRHID, params *public.ReadDomainParams) (string, error) {
	ret := _m.Called(xrhid, params)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*identity.XRHID, *public.ReadDomainParams) (string, error)); ok {
		return rf(xrhid, params)
	}
	if rf, ok := ret.Get(0).(func(*identity.XRHID, *public.ReadDomainParams) string); ok {
		r0 = rf(xrhid, params)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*identity.XRHID, *public.ReadDomainParams) error); ok {
		r1 = rf(xrhid, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: xrhid, params
func (_m *DomainInteractor) List(xrhid *identity.XRHID, params *public.ListDomainsParams) (string, int, int, error) {
	ret := _m.Called(xrhid, params)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 string
	var r1 int
	var r2 int
	var r3 error
	if rf, ok := ret.Get(0).(func(*identity.XRHID, *public.ListDomainsParams) (string, int, int, error)); ok {
		return rf(xrhid, params)
	}
	if rf, ok := ret.Get(0).(func(*identity.XRHID, *public.ListDomainsParams) string); ok {
		r0 = rf(xrhid, params)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*identity.XRHID, *public.ListDomainsParams) int); ok {
		r1 = rf(xrhid, params)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(*identity.XRHID, *public.ListDomainsParams) int); ok {
		r2 = rf(xrhid, params)
	} else {
		r2 = ret.Get(2).(int)
	}

	if rf, ok := ret.Get(3).(func(*identity.XRHID, *public.ListDomainsParams) error); ok {
		r3 = rf(xrhid, params)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Register provides a mock function with given fields: domainRegKey, xrhid, params, body
func (_m *DomainInteractor) Register(domainRegKey []byte, xrhid *identity.XRHID, params *public.RegisterDomainParams, body *public.Domain) (string, *header.XRHIDMVersion, *model.Domain, error) {
	ret := _m.Called(domainRegKey, xrhid, params, body)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 string
	var r1 *header.XRHIDMVersion
	var r2 *model.Domain
	var r3 error
	if rf, ok := ret.Get(0).(func([]byte, *identity.XRHID, *public.RegisterDomainParams, *public.Domain) (string, *header.XRHIDMVersion, *model.Domain, error)); ok {
		return rf(domainRegKey, xrhid, params, body)
	}
	if rf, ok := ret.Get(0).(func([]byte, *identity.XRHID, *public.RegisterDomainParams, *public.Domain) string); ok {
		r0 = rf(domainRegKey, xrhid, params, body)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func([]byte, *identity.XRHID, *public.RegisterDomainParams, *public.Domain) *header.XRHIDMVersion); ok {
		r1 = rf(domainRegKey, xrhid, params, body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*header.XRHIDMVersion)
		}
	}

	if rf, ok := ret.Get(2).(func([]byte, *identity.XRHID, *public.RegisterDomainParams, *public.Domain) *model.Domain); ok {
		r2 = rf(domainRegKey, xrhid, params, body)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*model.Domain)
		}
	}

	if rf, ok := ret.Get(3).(func([]byte, *identity.XRHID, *public.RegisterDomainParams, *public.Domain) error); ok {
		r3 = rf(domainRegKey, xrhid, params, body)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// UpdateAgent provides a mock function with given fields: xrhid, UUID, params, body
func (_m *DomainInteractor) UpdateAgent(xrhid *identity.XRHID, UUID uuid.UUID, params *public.UpdateDomainAgentParams, body *public.UpdateDomainAgentRequest) (string, *header.XRHIDMVersion, *model.Domain, error) {
	ret := _m.Called(xrhid, UUID, params, body)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAgent")
	}

	var r0 string
	var r1 *header.XRHIDMVersion
	var r2 *model.Domain
	var r3 error
	if rf, ok := ret.Get(0).(func(*identity.XRHID, uuid.UUID, *public.UpdateDomainAgentParams, *public.UpdateDomainAgentRequest) (string, *header.XRHIDMVersion, *model.Domain, error)); ok {
		return rf(xrhid, UUID, params, body)
	}
	if rf, ok := ret.Get(0).(func(*identity.XRHID, uuid.UUID, *public.UpdateDomainAgentParams, *public.UpdateDomainAgentRequest) string); ok {
		r0 = rf(xrhid, UUID, params, body)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*identity.XRHID, uuid.UUID, *public.UpdateDomainAgentParams, *public.UpdateDomainAgentRequest) *header.XRHIDMVersion); ok {
		r1 = rf(xrhid, UUID, params, body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*header.XRHIDMVersion)
		}
	}

	if rf, ok := ret.Get(2).(func(*identity.XRHID, uuid.UUID, *public.UpdateDomainAgentParams, *public.UpdateDomainAgentRequest) *model.Domain); ok {
		r2 = rf(xrhid, UUID, params, body)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*model.Domain)
		}
	}

	if rf, ok := ret.Get(3).(func(*identity.XRHID, uuid.UUID, *public.UpdateDomainAgentParams, *public.UpdateDomainAgentRequest) error); ok {
		r3 = rf(xrhid, UUID, params, body)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// UpdateUser provides a mock function with given fields: xrhid, UUID, params, body
func (_m *DomainInteractor) UpdateUser(xrhid *identity.XRHID, UUID uuid.UUID, params *public.UpdateDomainUserParams, body *public.UpdateDomainUserRequest) (string, *model.Domain, error) {
	ret := _m.Called(xrhid, UUID, params, body)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 string
	var r1 *model.Domain
	var r2 error
	if rf, ok := ret.Get(0).(func(*identity.XRHID, uuid.UUID, *public.UpdateDomainUserParams, *public.UpdateDomainUserRequest) (string, *model.Domain, error)); ok {
		return rf(xrhid, UUID, params, body)
	}
	if rf, ok := ret.Get(0).(func(*identity.XRHID, uuid.UUID, *public.UpdateDomainUserParams, *public.UpdateDomainUserRequest) string); ok {
		r0 = rf(xrhid, UUID, params, body)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*identity.XRHID, uuid.UUID, *public.UpdateDomainUserParams, *public.UpdateDomainUserRequest) *model.Domain); ok {
		r1 = rf(xrhid, UUID, params, body)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.Domain)
		}
	}

	if rf, ok := ret.Get(2).(func(*identity.XRHID, uuid.UUID, *public.UpdateDomainUserParams, *public.UpdateDomainUserRequest) error); ok {
		r2 = rf(xrhid, UUID, params, body)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewDomainInteractor creates a new instance of DomainInteractor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDomainInteractor(t interface {
	mock.TestingT
	Cleanup(func())
}) *DomainInteractor {
	mock := &DomainInteractor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
