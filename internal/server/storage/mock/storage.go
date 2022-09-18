// Code generated by MockGen. DO NOT EDIT.
// Source: ./interfaces.go

// Package storagemock is a generated GoMock package.
package storagemock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/sergalkin/gophkeeper/internal/server/model"
)

// MockUserServerStorage is a mock of UserServerStorage interface.
type MockUserServerStorage struct {
	ctrl     *gomock.Controller
	recorder *MockUserServerStorageMockRecorder
}

// MockUserServerStorageMockRecorder is the mock recorder for MockUserServerStorage.
type MockUserServerStorageMockRecorder struct {
	mock *MockUserServerStorage
}

// NewMockUserServerStorage creates a new mock instance.
func NewMockUserServerStorage(ctrl *gomock.Controller) *MockUserServerStorage {
	mock := &MockUserServerStorage{ctrl: ctrl}
	mock.recorder = &MockUserServerStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServerStorage) EXPECT() *MockUserServerStorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserServerStorage) Create(ctx context.Context, user model.User) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserServerStorageMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserServerStorage)(nil).Create), ctx, user)
}

// DeleteUser mocks base method.
func (m *MockUserServerStorage) DeleteUser(ctx context.Context, user model.User) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, user)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserServerStorageMockRecorder) DeleteUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserServerStorage)(nil).DeleteUser), ctx, user)
}

// GetByLoginAndPassword mocks base method.
func (m *MockUserServerStorage) GetByLoginAndPassword(ctx context.Context, user model.User) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByLoginAndPassword", ctx, user)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByLoginAndPassword indicates an expected call of GetByLoginAndPassword.
func (mr *MockUserServerStorageMockRecorder) GetByLoginAndPassword(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByLoginAndPassword", reflect.TypeOf((*MockUserServerStorage)(nil).GetByLoginAndPassword), ctx, user)
}

// MockSecretTypeServerStorage is a mock of SecretTypeServerStorage interface.
type MockSecretTypeServerStorage struct {
	ctrl     *gomock.Controller
	recorder *MockSecretTypeServerStorageMockRecorder
}

// MockSecretTypeServerStorageMockRecorder is the mock recorder for MockSecretTypeServerStorage.
type MockSecretTypeServerStorageMockRecorder struct {
	mock *MockSecretTypeServerStorage
}

// NewMockSecretTypeServerStorage creates a new mock instance.
func NewMockSecretTypeServerStorage(ctrl *gomock.Controller) *MockSecretTypeServerStorage {
	mock := &MockSecretTypeServerStorage{ctrl: ctrl}
	mock.recorder = &MockSecretTypeServerStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretTypeServerStorage) EXPECT() *MockSecretTypeServerStorageMockRecorder {
	return m.recorder
}

// GetSecretTypes mocks base method.
func (m *MockSecretTypeServerStorage) GetSecretTypes(ctx context.Context) ([]model.SecretType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretTypes", ctx)
	ret0, _ := ret[0].([]model.SecretType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretTypes indicates an expected call of GetSecretTypes.
func (mr *MockSecretTypeServerStorageMockRecorder) GetSecretTypes(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretTypes", reflect.TypeOf((*MockSecretTypeServerStorage)(nil).GetSecretTypes), ctx)
}

// MockSecretServerStorage is a mock of SecretServerStorage interface.
type MockSecretServerStorage struct {
	ctrl     *gomock.Controller
	recorder *MockSecretServerStorageMockRecorder
}

// MockSecretServerStorageMockRecorder is the mock recorder for MockSecretServerStorage.
type MockSecretServerStorageMockRecorder struct {
	mock *MockSecretServerStorage
}

// NewMockSecretServerStorage creates a new mock instance.
func NewMockSecretServerStorage(ctrl *gomock.Controller) *MockSecretServerStorage {
	mock := &MockSecretServerStorage{ctrl: ctrl}
	mock.recorder = &MockSecretServerStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretServerStorage) EXPECT() *MockSecretServerStorageMockRecorder {
	return m.recorder
}

// CreateSecret mocks base method.
func (m *MockSecretServerStorage) CreateSecret(ctx context.Context, secret model.Secret) (model.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSecret", ctx, secret)
	ret0, _ := ret[0].(model.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSecret indicates an expected call of CreateSecret.
func (mr *MockSecretServerStorageMockRecorder) CreateSecret(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecret", reflect.TypeOf((*MockSecretServerStorage)(nil).CreateSecret), ctx, secret)
}

// DeleteSecret mocks base method.
func (m *MockSecretServerStorage) DeleteSecret(ctx context.Context, secret model.Secret) (model.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecret", ctx, secret)
	ret0, _ := ret[0].(model.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockSecretServerStorageMockRecorder) DeleteSecret(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockSecretServerStorage)(nil).DeleteSecret), ctx, secret)
}

// GetSecret mocks base method.
func (m *MockSecretServerStorage) GetSecret(ctx context.Context, secret model.Secret) (model.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", ctx, secret)
	ret0, _ := ret[0].(model.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockSecretServerStorageMockRecorder) GetSecret(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecretServerStorage)(nil).GetSecret), ctx, secret)
}
