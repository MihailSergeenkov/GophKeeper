// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/server/main.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWebServer is a mock of WebServer interface.
type MockWebServer struct {
	ctrl     *gomock.Controller
	recorder *MockWebServerMockRecorder
}

// MockWebServerMockRecorder is the mock recorder for MockWebServer.
type MockWebServerMockRecorder struct {
	mock *MockWebServer
}

// NewMockWebServer creates a new mock instance.
func NewMockWebServer(ctrl *gomock.Controller) *MockWebServer {
	mock := &MockWebServer{ctrl: ctrl}
	mock.recorder = &MockWebServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebServer) EXPECT() *MockWebServerMockRecorder {
	return m.recorder
}

// ListenAndServe mocks base method.
func (m *MockWebServer) ListenAndServe() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenAndServe")
	ret0, _ := ret[0].(error)
	return ret0
}

// ListenAndServe indicates an expected call of ListenAndServe.
func (mr *MockWebServerMockRecorder) ListenAndServe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServe", reflect.TypeOf((*MockWebServer)(nil).ListenAndServe))
}

// ListenAndServeTLS mocks base method.
func (m *MockWebServer) ListenAndServeTLS(certFile, keyFile string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenAndServeTLS", certFile, keyFile)
	ret0, _ := ret[0].(error)
	return ret0
}

// ListenAndServeTLS indicates an expected call of ListenAndServeTLS.
func (mr *MockWebServerMockRecorder) ListenAndServeTLS(certFile, keyFile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServeTLS", reflect.TypeOf((*MockWebServer)(nil).ListenAndServeTLS), certFile, keyFile)
}
