// Code generated by mockery v2.40.3. DO NOT EDIT.

package presenter

import (
	model "github.com/podengo-project/idmsvc-backend/internal/domain/model"
	mock "github.com/stretchr/testify/mock"

	public "github.com/podengo-project/idmsvc-backend/internal/api/public"

	repository "github.com/podengo-project/idmsvc-backend/internal/interface/repository"
)

// DomainPresenter is an autogenerated mock type for the DomainPresenter type
type DomainPresenter struct {
	mock.Mock
}

// CreateDomainToken provides a mock function with given fields: token
func (_m *DomainPresenter) CreateDomainToken(token *repository.DomainRegToken) (*public.DomainRegToken, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for CreateDomainToken")
	}

	var r0 *public.DomainRegToken
	var r1 error
	if rf, ok := ret.Get(0).(func(*repository.DomainRegToken) (*public.DomainRegToken, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(*repository.DomainRegToken) *public.DomainRegToken); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.DomainRegToken)
		}
	}

	if rf, ok := ret.Get(1).(func(*repository.DomainRegToken) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: domain
func (_m *DomainPresenter) Get(domain *model.Domain) (*public.Domain, error) {
	ret := _m.Called(domain)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *public.Domain
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Domain) (*public.Domain, error)); ok {
		return rf(domain)
	}
	if rf, ok := ret.Get(0).(func(*model.Domain) *public.Domain); ok {
		r0 = rf(domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.Domain)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Domain) error); ok {
		r1 = rf(domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: count, offset, limit, data
func (_m *DomainPresenter) List(count int64, offset int, limit int, data []model.Domain) (*public.ListDomainsResponseSchema, error) {
	ret := _m.Called(count, offset, limit, data)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 *public.ListDomainsResponseSchema
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int, int, []model.Domain) (*public.ListDomainsResponseSchema, error)); ok {
		return rf(count, offset, limit, data)
	}
	if rf, ok := ret.Get(0).(func(int64, int, int, []model.Domain) *public.ListDomainsResponseSchema); ok {
		r0 = rf(count, offset, limit, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.ListDomainsResponseSchema)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int, int, []model.Domain) error); ok {
		r1 = rf(count, offset, limit, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: domain
func (_m *DomainPresenter) Register(domain *model.Domain) (*public.Domain, error) {
	ret := _m.Called(domain)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 *public.Domain
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Domain) (*public.Domain, error)); ok {
		return rf(domain)
	}
	if rf, ok := ret.Get(0).(func(*model.Domain) *public.Domain); ok {
		r0 = rf(domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.Domain)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Domain) error); ok {
		r1 = rf(domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAgent provides a mock function with given fields: domain
func (_m *DomainPresenter) UpdateAgent(domain *model.Domain) (*public.Domain, error) {
	ret := _m.Called(domain)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAgent")
	}

	var r0 *public.Domain
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Domain) (*public.Domain, error)); ok {
		return rf(domain)
	}
	if rf, ok := ret.Get(0).(func(*model.Domain) *public.Domain); ok {
		r0 = rf(domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.Domain)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Domain) error); ok {
		r1 = rf(domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: domain
func (_m *DomainPresenter) UpdateUser(domain *model.Domain) (*public.Domain, error) {
	ret := _m.Called(domain)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 *public.Domain
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Domain) (*public.Domain, error)); ok {
		return rf(domain)
	}
	if rf, ok := ret.Get(0).(func(*model.Domain) *public.Domain); ok {
		r0 = rf(domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.Domain)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Domain) error); ok {
		r1 = rf(domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDomainPresenter creates a new instance of DomainPresenter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDomainPresenter(t interface {
	mock.TestingT
	Cleanup(func())
}) *DomainPresenter {
	mock := &DomainPresenter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
