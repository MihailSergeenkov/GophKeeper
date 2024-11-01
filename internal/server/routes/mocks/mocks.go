// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/routes/routes.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	http "net/http"
	reflect "reflect"

	models "github.com/MihailSergeenkov/GophKeeper/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockHandlerer is a mock of Handlerer interface.
type MockHandlerer struct {
	ctrl     *gomock.Controller
	recorder *MockHandlererMockRecorder
}

// MockHandlererMockRecorder is the mock recorder for MockHandlerer.
type MockHandlererMockRecorder struct {
	mock *MockHandlerer
}

// NewMockHandlerer creates a new mock instance.
func NewMockHandlerer(ctrl *gomock.Controller) *MockHandlerer {
	mock := &MockHandlerer{ctrl: ctrl}
	mock.recorder = &MockHandlererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandlerer) EXPECT() *MockHandlererMockRecorder {
	return m.recorder
}

// AddCard mocks base method.
func (m *MockHandlerer) AddCard() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCard")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// AddCard indicates an expected call of AddCard.
func (mr *MockHandlererMockRecorder) AddCard() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCard", reflect.TypeOf((*MockHandlerer)(nil).AddCard))
}

// AddFile mocks base method.
func (m *MockHandlerer) AddFile() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFile")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// AddFile indicates an expected call of AddFile.
func (mr *MockHandlererMockRecorder) AddFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFile", reflect.TypeOf((*MockHandlerer)(nil).AddFile))
}

// AddPassword mocks base method.
func (m *MockHandlerer) AddPassword() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPassword")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// AddPassword indicates an expected call of AddPassword.
func (mr *MockHandlererMockRecorder) AddPassword() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPassword", reflect.TypeOf((*MockHandlerer)(nil).AddPassword))
}

// AddText mocks base method.
func (m *MockHandlerer) AddText() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddText")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// AddText indicates an expected call of AddText.
func (mr *MockHandlererMockRecorder) AddText() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddText", reflect.TypeOf((*MockHandlerer)(nil).AddText))
}

// CreateUserToken mocks base method.
func (m *MockHandlerer) CreateUserToken() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserToken")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// CreateUserToken indicates an expected call of CreateUserToken.
func (mr *MockHandlererMockRecorder) CreateUserToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserToken", reflect.TypeOf((*MockHandlerer)(nil).CreateUserToken))
}

// FetchUserData mocks base method.
func (m *MockHandlerer) FetchUserData() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchUserData")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// FetchUserData indicates an expected call of FetchUserData.
func (mr *MockHandlererMockRecorder) FetchUserData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUserData", reflect.TypeOf((*MockHandlerer)(nil).FetchUserData))
}

// GetCard mocks base method.
func (m *MockHandlerer) GetCard() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCard")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// GetCard indicates an expected call of GetCard.
func (mr *MockHandlererMockRecorder) GetCard() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCard", reflect.TypeOf((*MockHandlerer)(nil).GetCard))
}

// GetFile mocks base method.
func (m *MockHandlerer) GetFile() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// GetFile indicates an expected call of GetFile.
func (mr *MockHandlererMockRecorder) GetFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockHandlerer)(nil).GetFile))
}

// GetPassword mocks base method.
func (m *MockHandlerer) GetPassword() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPassword")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// GetPassword indicates an expected call of GetPassword.
func (mr *MockHandlererMockRecorder) GetPassword() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPassword", reflect.TypeOf((*MockHandlerer)(nil).GetPassword))
}

// GetText mocks base method.
func (m *MockHandlerer) GetText() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetText")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// GetText indicates an expected call of GetText.
func (mr *MockHandlererMockRecorder) GetText() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetText", reflect.TypeOf((*MockHandlerer)(nil).GetText))
}

// Ping mocks base method.
func (m *MockHandlerer) Ping() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockHandlererMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockHandlerer)(nil).Ping))
}

// RegisterUser mocks base method.
func (m *MockHandlerer) RegisterUser() http.HandlerFunc {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser")
	ret0, _ := ret[0].(http.HandlerFunc)
	return ret0
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockHandlererMockRecorder) RegisterUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockHandlerer)(nil).RegisterUser))
}

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// GetUserByID mocks base method.
func (m *MockStorager) GetUserByID(ctx context.Context, userID int) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockStoragerMockRecorder) GetUserByID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockStorager)(nil).GetUserByID), ctx, userID)
}
