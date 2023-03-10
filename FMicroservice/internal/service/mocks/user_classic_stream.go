// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	model "GolangInternship/FMicroservice/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UserClassicStream is an autogenerated mock type for the UserClassicStream type
type UserClassicStream struct {
	mock.Mock
}

// ProduceUser provides a mock function with given fields: ctx, user
func (_m *UserClassicStream) ProduceUser(ctx context.Context, user *model.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserClassicStream interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserClassicStream creates a new instance of UserClassicStream. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserClassicStream(t mockConstructorTestingTNewUserClassicStream) *UserClassicStream {
	mock := &UserClassicStream{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
