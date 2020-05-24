// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Authenticator is an autogenerated mock type for the Authenticator type
type Authenticator struct {
	mock.Mock
}

// Authenticate provides a mock function with given fields: authHeader
func (_m *Authenticator) Authenticate(authHeader string) (map[string]interface{}, error) {
	ret := _m.Called(authHeader)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(string) map[string]interface{}); ok {
		r0 = rf(authHeader)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(authHeader)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}