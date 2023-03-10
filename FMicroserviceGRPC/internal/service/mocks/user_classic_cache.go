// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	model "GolangInternship/FMicroserviceGRPC/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UserClassicCache is an autogenerated mock type for the UserClassicCache type
type UserClassicCache struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *UserClassicCache) CreateUser(ctx context.Context, user *model.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByLogin provides a mock function with given fields: ctx, login
func (_m *UserClassicCache) GetByLogin(ctx context.Context, login string) (*model.User, bool, error) {
	ret := _m.Called(ctx, login)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, login)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewUserClassicCache interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserClassicCache creates a new instance of UserClassicCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserClassicCache(t mockConstructorTestingTNewUserClassicCache) *UserClassicCache {
	mock := &UserClassicCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
