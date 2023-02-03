// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
)

// User is an autogenerated mock type for the User type
type User struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, login
func (_m *User) Delete(ctx context.Context, login string) error {
	ret := _m.Called(ctx, login)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, login)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByLogin provides a mock function with given fields: ctx, login
func (_m *User) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	ret := _m.Called(ctx, login)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, login, password
func (_m *User) Login(ctx context.Context, login string, password string) (string, string, error) {
	ret := _m.Called(ctx, login, password)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, login, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string, string) string); ok {
		r1 = rf(ctx, login, password)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, login, password)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Refresh provides a mock function with given fields: ctx, login, userRefreshToken
func (_m *User) Refresh(ctx context.Context, login string, userRefreshToken string) (string, string, error) {
	ret := _m.Called(ctx, login, userRefreshToken)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, login, userRefreshToken)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string, string) string); ok {
		r1 = rf(ctx, login, userRefreshToken)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, login, userRefreshToken)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Signup provides a mock function with given fields: ctx, user
func (_m *User) Signup(ctx context.Context, user *model.User) (string, string, *model.User, error) {
	ret := _m.Called(ctx, user)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) string); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, *model.User) string); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 *model.User
	if rf, ok := ret.Get(2).(func(context.Context, *model.User) *model.User); ok {
		r2 = rf(ctx, user)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*model.User)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, *model.User) error); ok {
		r3 = rf(ctx, user)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Update provides a mock function with given fields: ctx, login, user
func (_m *User) Update(ctx context.Context, login string, user *model.User) error {
	ret := _m.Called(ctx, login, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.User) error); ok {
		r0 = rf(ctx, login, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUser interface {
	mock.TestingT
	Cleanup(func())
}

// NewUser creates a new instance of User. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUser(t mockConstructorTestingTNewUser) *User {
	mock := &User{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
