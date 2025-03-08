// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// URLDeleter is an autogenerated mock type for the URLDeleter type
type URLDeleter struct {
	mock.Mock
}

// DeleteURL provides a mock function with given fields: alias
func (_m *URLDeleter) DeleteURL(alias string) (int64, error) {
	ret := _m.Called(alias)

	if len(ret) == 0 {
		panic("no return value specified for DeleteURL")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int64, error)); ok {
		return rf(alias)
	}
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewURLDeleter creates a new instance of URLDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLDeleter {
	mock := &URLDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
