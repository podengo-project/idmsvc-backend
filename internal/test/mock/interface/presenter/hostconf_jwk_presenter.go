// Code generated by mockery v2.53.3. DO NOT EDIT.

package presenter

import (
	mock "github.com/stretchr/testify/mock"

	public "github.com/podengo-project/idmsvc-backend/internal/api/public"
)

// HostconfJwkPresenter is an autogenerated mock type for the HostconfJwkPresenter type
type HostconfJwkPresenter struct {
	mock.Mock
}

// PublicSigningKeys provides a mock function with given fields: keys, revokedKids
func (_m *HostconfJwkPresenter) PublicSigningKeys(keys []string, revokedKids []string) (*public.SigningKeysResponse, error) {
	ret := _m.Called(keys, revokedKids)

	if len(ret) == 0 {
		panic("no return value specified for PublicSigningKeys")
	}

	var r0 *public.SigningKeysResponse
	var r1 error
	if rf, ok := ret.Get(0).(func([]string, []string) (*public.SigningKeysResponse, error)); ok {
		return rf(keys, revokedKids)
	}
	if rf, ok := ret.Get(0).(func([]string, []string) *public.SigningKeysResponse); ok {
		r0 = rf(keys, revokedKids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*public.SigningKeysResponse)
		}
	}

	if rf, ok := ret.Get(1).(func([]string, []string) error); ok {
		r1 = rf(keys, revokedKids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewHostconfJwkPresenter creates a new instance of HostconfJwkPresenter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHostconfJwkPresenter(t interface {
	mock.TestingT
	Cleanup(func())
}) *HostconfJwkPresenter {
	mock := &HostconfJwkPresenter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
