// Code generated by mockery v2.34.2. DO NOT EDIT.

package models

import (
	context "context"

	boil "github.com/volatiletech/sqlboiler/v4/boil"

	mock "github.com/stretchr/testify/mock"
)

// MockSchemaMigrationHook is an autogenerated mock type for the SchemaMigrationHook type
type MockSchemaMigrationHook struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockSchemaMigrationHook) Execute(_a0 context.Context, _a1 boil.ContextExecutor, _a2 *SchemaMigration) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, boil.ContextExecutor, *SchemaMigration) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockSchemaMigrationHook creates a new instance of MockSchemaMigrationHook. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSchemaMigrationHook(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSchemaMigrationHook {
	mock := &MockSchemaMigrationHook{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}