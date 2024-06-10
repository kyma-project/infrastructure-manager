// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// KubeconfigProvider is an autogenerated mock type for the KubeconfigProvider type
type KubeconfigProvider struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: ctx, shootName
func (_m *KubeconfigProvider) Fetch(ctx context.Context, shootName string) (string, error) {
	ret := _m.Called(ctx, shootName)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, shootName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, shootName)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, shootName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewKubeconfigProvider creates a new instance of KubeconfigProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewKubeconfigProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *KubeconfigProvider {
	mock := &KubeconfigProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
