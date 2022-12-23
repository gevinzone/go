// Code generated by MockGen. DO NOT EDIT.
// Source: ./person.go

// Package person_mock is a generated GoMock package.
package person_mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPerson is a mock of Person interface.
type MockPerson struct {
	ctrl     *gomock.Controller
	recorder *MockPersonMockRecorder
}

// MockPersonMockRecorder is the mock recorder for MockPerson.
type MockPersonMockRecorder struct {
	mock *MockPerson
}

// NewMockPerson creates a new mock instance.
func NewMockPerson(ctrl *gomock.Controller) *MockPerson {
	mock := &MockPerson{ctrl: ctrl}
	mock.recorder = &MockPersonMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPerson) EXPECT() *MockPersonMockRecorder {
	return m.recorder
}

// GetName mocks base method.
func (m *MockPerson) GetName(id int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetName indicates an expected call of GetName.
func (mr *MockPersonMockRecorder) GetName(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockPerson)(nil).GetName), id)
}
