// Code generated by MockGen. DO NOT EDIT.
// Source: status.go

// Package mock_status is a generated GoMock package.
package mock_status

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	auth "github.com/gotd/td/telegram/auth"
)

// tgAuth is a mock of tgAuth interface.
type tgAuth struct {
	ctrl     *gomock.Controller
	recorder *tgAuthMockRecorder
}

// tgAuthMockRecorder is the mock recorder for tgAuth.
type tgAuthMockRecorder struct {
	mock *tgAuth
}

// NewtgAuth creates a new mock instance.
func NewtgAuth(ctrl *gomock.Controller) *tgAuth {
	mock := &tgAuth{ctrl: ctrl}
	mock.recorder = &tgAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *tgAuth) EXPECT() *tgAuthMockRecorder {
	return m.recorder
}

// Status mocks base method.
func (m *tgAuth) Status(ctx context.Context) (*auth.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", ctx)
	ret0, _ := ret[0].(*auth.Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status.
func (mr *tgAuthMockRecorder) Status(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*tgAuth)(nil).Status), ctx)
}

// MockCheckerAuthStatus is a mock of CheckerAuthStatus interface.
type MockCheckerAuthStatus struct {
	ctrl     *gomock.Controller
	recorder *MockCheckerAuthStatusMockRecorder
}

// MockCheckerAuthStatusMockRecorder is the mock recorder for MockCheckerAuthStatus.
type MockCheckerAuthStatusMockRecorder struct {
	mock *MockCheckerAuthStatus
}

// NewMockCheckerAuthStatus creates a new mock instance.
func NewMockCheckerAuthStatus(ctrl *gomock.Controller) *MockCheckerAuthStatus {
	mock := &MockCheckerAuthStatus{ctrl: ctrl}
	mock.recorder = &MockCheckerAuthStatusMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCheckerAuthStatus) EXPECT() *MockCheckerAuthStatusMockRecorder {
	return m.recorder
}

// CheckAuth mocks base method.
func (m *MockCheckerAuthStatus) CheckAuth(ctx context.Context, auth tgAuth) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuth", ctx, auth)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAuth indicates an expected call of CheckAuth.
func (mr *MockCheckerAuthStatusMockRecorder) CheckAuth(ctx, auth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuth", reflect.TypeOf((*MockCheckerAuthStatus)(nil).CheckAuth), ctx, auth)
}
