// Code generated by mockery v2.0.0. DO NOT EDIT.

package gawetest

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Job is an autogenerated mock type for the Job type
type Job struct {
	mock.Mock
}

// Exec provides a mock function with given fields: ctx
func (_m *Job) Exec(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// JobID provides a mock function with given fields:
func (_m *Job) JobID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// JobType provides a mock function with given fields:
func (_m *Job) JobType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Tags provides a mock function with given fields:
func (_m *Job) Tags() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
