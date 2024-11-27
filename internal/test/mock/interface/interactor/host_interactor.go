// Code generated by mockery v2.46.0. DO NOT EDIT.

package interactor

import (
	interactor "github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"

	mock "github.com/stretchr/testify/mock"

	public "github.com/podengo-project/idmsvc-backend/internal/api/public"

	uuid "github.com/google/uuid"
)

// HostInteractor is an autogenerated mock type for the HostInteractor type
type HostInteractor struct {
	mock.Mock
}

// HostConf provides a mock function with given fields: xrhid, inventoryId, fqdn, params, body
func (_m *HostInteractor) HostConf(xrhid *identity.XRHID, inventoryId uuid.UUID, fqdn string, params *public.HostConfParams, body *public.HostConf) (*interactor.HostConfOptions, error) {
	ret := _m.Called(xrhid, inventoryId, fqdn, params, body)

	if len(ret) == 0 {
		panic("no return value specified for HostConf")
	}

	var r0 *interactor.HostConfOptions
	var r1 error
	if rf, ok := ret.Get(0).(func(*identity.XRHID, uuid.UUID, string, *public.HostConfParams, *public.HostConf) (*interactor.HostConfOptions, error)); ok {
		return rf(xrhid, inventoryId, fqdn, params, body)
	}
	if rf, ok := ret.Get(0).(func(*identity.XRHID, uuid.UUID, string, *public.HostConfParams, *public.HostConf) *interactor.HostConfOptions); ok {
		r0 = rf(xrhid, inventoryId, fqdn, params, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*interactor.HostConfOptions)
		}
	}

	if rf, ok := ret.Get(1).(func(*identity.XRHID, uuid.UUID, string, *public.HostConfParams, *public.HostConf) error); ok {
		r1 = rf(xrhid, inventoryId, fqdn, params, body)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewHostInteractor creates a new instance of HostInteractor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHostInteractor(t interface {
	mock.TestingT
	Cleanup(func())
}) *HostInteractor {
	mock := &HostInteractor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
