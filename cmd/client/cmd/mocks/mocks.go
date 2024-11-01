// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/client/cmd/root.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/MihailSergeenkov/GophKeeper/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockServicer is a mock of Servicer interface.
type MockServicer struct {
	ctrl     *gomock.Controller
	recorder *MockServicerMockRecorder
}

// MockServicerMockRecorder is the mock recorder for MockServicer.
type MockServicerMockRecorder struct {
	mock *MockServicer
}

// NewMockServicer creates a new mock instance.
func NewMockServicer(ctrl *gomock.Controller) *MockServicer {
	mock := &MockServicer{ctrl: ctrl}
	mock.recorder = &MockServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServicer) EXPECT() *MockServicerMockRecorder {
	return m.recorder
}

// AddCard mocks base method.
func (m *MockServicer) AddCard(req *models.AddCardRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCard", req)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCard indicates an expected call of AddCard.
func (mr *MockServicerMockRecorder) AddCard(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCard", reflect.TypeOf((*MockServicer)(nil).AddCard), req)
}

// AddFile mocks base method.
func (m *MockServicer) AddFile(filePath, mark, description string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFile", filePath, mark, description)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFile indicates an expected call of AddFile.
func (mr *MockServicerMockRecorder) AddFile(filePath, mark, description interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFile", reflect.TypeOf((*MockServicer)(nil).AddFile), filePath, mark, description)
}

// AddPassword mocks base method.
func (m *MockServicer) AddPassword(req models.AddPasswordRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPassword", req)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPassword indicates an expected call of AddPassword.
func (mr *MockServicerMockRecorder) AddPassword(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPassword", reflect.TypeOf((*MockServicer)(nil).AddPassword), req)
}

// AddText mocks base method.
func (m *MockServicer) AddText(req models.AddTextRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddText", req)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddText indicates an expected call of AddText.
func (mr *MockServicerMockRecorder) AddText(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddText", reflect.TypeOf((*MockServicer)(nil).AddText), req)
}

// GetCard mocks base method.
func (m *MockServicer) GetCard(id string) (models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCard", id)
	ret0, _ := ret[0].(models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCard indicates an expected call of GetCard.
func (mr *MockServicerMockRecorder) GetCard(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCard", reflect.TypeOf((*MockServicer)(nil).GetCard), id)
}

// GetData mocks base method.
func (m *MockServicer) GetData() []models.UserData {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetData")
	ret0, _ := ret[0].([]models.UserData)
	return ret0
}

// GetData indicates an expected call of GetData.
func (mr *MockServicerMockRecorder) GetData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetData", reflect.TypeOf((*MockServicer)(nil).GetData))
}

// GetFile mocks base method.
func (m *MockServicer) GetFile(id, dir string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile", id, dir)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetFile indicates an expected call of GetFile.
func (mr *MockServicerMockRecorder) GetFile(id, dir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockServicer)(nil).GetFile), id, dir)
}

// GetPassword mocks base method.
func (m *MockServicer) GetPassword(id string) (models.Password, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPassword", id)
	ret0, _ := ret[0].(models.Password)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPassword indicates an expected call of GetPassword.
func (mr *MockServicerMockRecorder) GetPassword(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPassword", reflect.TypeOf((*MockServicer)(nil).GetPassword), id)
}

// GetText mocks base method.
func (m *MockServicer) GetText(id string) (models.Text, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetText", id)
	ret0, _ := ret[0].(models.Text)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetText indicates an expected call of GetText.
func (mr *MockServicerMockRecorder) GetText(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetText", reflect.TypeOf((*MockServicer)(nil).GetText), id)
}

// LoginUser mocks base method.
func (m *MockServicer) LoginUser(req models.CreateUserTokenRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", req)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockServicerMockRecorder) LoginUser(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockServicer)(nil).LoginUser), req)
}

// LogoutUser mocks base method.
func (m *MockServicer) LogoutUser() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogoutUser")
	ret0, _ := ret[0].(error)
	return ret0
}

// LogoutUser indicates an expected call of LogoutUser.
func (mr *MockServicerMockRecorder) LogoutUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogoutUser", reflect.TypeOf((*MockServicer)(nil).LogoutUser))
}

// RegisterUser mocks base method.
func (m *MockServicer) RegisterUser(req models.RegisterUserRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", req)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockServicerMockRecorder) RegisterUser(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockServicer)(nil).RegisterUser), req)
}

// SyncData mocks base method.
func (m *MockServicer) SyncData() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncData")
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncData indicates an expected call of SyncData.
func (mr *MockServicerMockRecorder) SyncData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncData", reflect.TypeOf((*MockServicer)(nil).SyncData))
}
