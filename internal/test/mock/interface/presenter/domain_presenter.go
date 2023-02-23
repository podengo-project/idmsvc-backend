// Code generated by mockery v2.16.0. DO NOT EDIT.

package presenter

import (
	model "github.com/hmsidm/internal/domain/model"
	mock "github.com/stretchr/testify/mock"

	public "github.com/hmsidm/internal/api/public"
)

// DomainPresenter is an autogenerated mock type for the DomainPresenter type
type DomainPresenter struct {
	mock.Mock
}

// Create provides a mock function with given fields: domain
func (_m *DomainPresenter) Create(domain *model.Domain) (*public.CreateDomainResponseSchema, error) {
	ret := _m.Called(domain)

	var r0 *public.CreateDomainResponseSchema
	if rf, ok := ret.Get(0).(func(*model.Domain) *public.CreateDomainResponseSchema); ok {
		r0 = rf(domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.CreateDomainResponseSchema)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Domain) error); ok {
		r1 = rf(domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: domain
func (_m *DomainPresenter) Get(domain *model.Domain) (*public.ReadDomainResponseSchema, error) {
	ret := _m.Called(domain)

	var r0 *public.ReadDomainResponseSchema
	if rf, ok := ret.Get(0).(func(*model.Domain) *public.ReadDomainResponseSchema); ok {
		r0 = rf(domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.ReadDomainResponseSchema)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Domain) error); ok {
		r1 = rf(domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: prefix, offset, count, data
func (_m *DomainPresenter) List(prefix string, offset int64, count int32, data []model.Domain) (*public.ListDomainsResponseSchema, error) {
	ret := _m.Called(prefix, offset, count, data)

	var r0 *public.ListDomainsResponseSchema
	if rf, ok := ret.Get(0).(func(string, int64, int32, []model.Domain) *public.ListDomainsResponseSchema); ok {
		r0 = rf(prefix, offset, count, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.ListDomainsResponseSchema)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int64, int32, []model.Domain) error); ok {
		r1 = rf(prefix, offset, count, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewDomainPresenter interface {
	mock.TestingT
	Cleanup(func())
}

// NewDomainPresenter creates a new instance of DomainPresenter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDomainPresenter(t mockConstructorTestingTNewDomainPresenter) *DomainPresenter {
	mock := &DomainPresenter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}