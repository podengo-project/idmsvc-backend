// Code generated by mockery v2.53.3. DO NOT EDIT.

package repository

import (
	context "context"

	jwk "github.com/lestrrat-go/jwx/v2/jwk"
	mock "github.com/stretchr/testify/mock"

	model "github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/hostconf_jwk/model"
)

// HostconfJwkRepository is an autogenerated mock type for the HostconfJwkRepository type
type HostconfJwkRepository struct {
	mock.Mock
}

// GetPrivateSigningKeys provides a mock function with given fields: ctx
func (_m *HostconfJwkRepository) GetPrivateSigningKeys(ctx context.Context) ([]jwk.Key, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetPrivateSigningKeys")
	}

	var r0 []jwk.Key
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]jwk.Key, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []jwk.Key); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]jwk.Key)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPublicKeyArray provides a mock function with given fields: ctx
func (_m *HostconfJwkRepository) GetPublicKeyArray(ctx context.Context) ([]string, []string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetPublicKeyArray")
	}

	var r0 []string
	var r1 []string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, []string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) []string); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(ctx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// InsertJWK provides a mock function with given fields: ctx, hcjwk
func (_m *HostconfJwkRepository) InsertJWK(ctx context.Context, hcjwk *model.HostconfJwk) error {
	ret := _m.Called(ctx, hcjwk)

	if len(ret) == 0 {
		panic("no return value specified for InsertJWK")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.HostconfJwk) error); ok {
		r0 = rf(ctx, hcjwk)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListJWKs provides a mock function with given fields: ctx
func (_m *HostconfJwkRepository) ListJWKs(ctx context.Context) ([]model.HostconfJwk, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListJWKs")
	}

	var r0 []model.HostconfJwk
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]model.HostconfJwk, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []model.HostconfJwk); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.HostconfJwk)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PurgeExpiredJWKs provides a mock function with given fields: ctx
func (_m *HostconfJwkRepository) PurgeExpiredJWKs(ctx context.Context) ([]model.HostconfJwk, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for PurgeExpiredJWKs")
	}

	var r0 []model.HostconfJwk
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]model.HostconfJwk, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []model.HostconfJwk); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.HostconfJwk)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RevokeJWK provides a mock function with given fields: ctx, kid
func (_m *HostconfJwkRepository) RevokeJWK(ctx context.Context, kid string) (*model.HostconfJwk, error) {
	ret := _m.Called(ctx, kid)

	if len(ret) == 0 {
		panic("no return value specified for RevokeJWK")
	}

	var r0 *model.HostconfJwk
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.HostconfJwk, error)); ok {
		return rf(ctx, kid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.HostconfJwk); ok {
		r0 = rf(ctx, kid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.HostconfJwk)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, kid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewHostconfJwkRepository creates a new instance of HostconfJwkRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHostconfJwkRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *HostconfJwkRepository {
	mock := &HostconfJwkRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
