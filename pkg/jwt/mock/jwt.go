// Code generated by MockGen. DO NOT EDIT.
// Source: ./jwt.go

// Package jwtmock is a generated GoMock package.
package jwtmock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Decode mocks base method.
func (m *MockManager) Decode(token string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decode indicates an expected call of Decode.
func (mr *MockManagerMockRecorder) Decode(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockManager)(nil).Decode), token)
}

// Issue mocks base method.
func (m *MockManager) Issue(id string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Issue", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Issue indicates an expected call of Issue.
func (mr *MockManagerMockRecorder) Issue(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Issue", reflect.TypeOf((*MockManager)(nil).Issue), id)
}
