// Code generated by MockGen. DO NOT EDIT.
// Source: log.go

// Package mock_log is a generated GoMock package.
package mock_log

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockinnerLog is a mock of innerLog interface.
type MockinnerLog struct {
	ctrl     *gomock.Controller
	recorder *MockinnerLogMockRecorder
}

// MockinnerLogMockRecorder is the mock recorder for MockinnerLog.
type MockinnerLogMockRecorder struct {
	mock *MockinnerLog
}

// NewMockinnerLog creates a new mock instance.
func NewMockinnerLog(ctrl *gomock.Controller) *MockinnerLog {
	mock := &MockinnerLog{ctrl: ctrl}
	mock.recorder = &MockinnerLogMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockinnerLog) EXPECT() *MockinnerLogMockRecorder {
	return m.recorder
}

// Errorln mocks base method.
func (m *MockinnerLog) Errorln(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorln", varargs...)
}

// Errorln indicates an expected call of Errorln.
func (mr *MockinnerLogMockRecorder) Errorln(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorln", reflect.TypeOf((*MockinnerLog)(nil).Errorln), args...)
}

// Infoln mocks base method.
func (m *MockinnerLog) Infoln(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infoln", varargs...)
}

// Infoln indicates an expected call of Infoln.
func (mr *MockinnerLogMockRecorder) Infoln(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infoln", reflect.TypeOf((*MockinnerLog)(nil).Infoln), args...)
}

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockLogger) Error(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), args...)
}

// Info mocks base method.
func (m *MockLogger) Info(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), args...)
}
