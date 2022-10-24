// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/molchalin/workshop3/internal/exchange (interfaces: Client)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// ExchangeRate mocks base method.
func (m *MockClient) ExchangeRate(arg0 context.Context, arg1, arg2 string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeRate", arg0, arg1, arg2)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExchangeRate indicates an expected call of ExchangeRate.
func (mr *MockClientMockRecorder) ExchangeRate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeRate", reflect.TypeOf((*MockClient)(nil).ExchangeRate), arg0, arg1, arg2)
}
