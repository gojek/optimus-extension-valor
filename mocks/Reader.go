// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	model "github.com/gojek/optimus-extension-valor/model"
	mock "github.com/stretchr/testify/mock"
)

// Reader is an autogenerated mock type for the Reader type
type Reader struct {
	mock.Mock
}

// ReadAll provides a mock function with given fields:
func (_m *Reader) ReadAll() ([]*model.Data, model.Error) {
	ret := _m.Called()

	var r0 []*model.Data
	if rf, ok := ret.Get(0).(func() []*model.Data); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Data)
		}
	}

	var r1 model.Error
	if rf, ok := ret.Get(1).(func() model.Error); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(model.Error)
		}
	}

	return r0, r1
}

// ReadOne provides a mock function with given fields:
func (_m *Reader) ReadOne() (*model.Data, model.Error) {
	ret := _m.Called()

	var r0 *model.Data
	if rf, ok := ret.Get(0).(func() *model.Data); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Data)
		}
	}

	var r1 model.Error
	if rf, ok := ret.Get(1).(func() model.Error); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(model.Error)
		}
	}

	return r0, r1
}
