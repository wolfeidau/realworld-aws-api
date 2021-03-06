// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/wolfeidau/realworld-aws-api/internal/stores (interfaces: Customers)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	stores "github.com/wolfeidau/realworld-aws-api/internal/stores"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	reflect "reflect"
)

// MockCustomers is a mock of Customers interface
type MockCustomers struct {
	ctrl     *gomock.Controller
	recorder *MockCustomersMockRecorder
}

// MockCustomersMockRecorder is the mock recorder for MockCustomers
type MockCustomersMockRecorder struct {
	mock *MockCustomers
}

// NewMockCustomers creates a new mock instance
func NewMockCustomers(ctrl *gomock.Controller) *MockCustomers {
	mock := &MockCustomers{ctrl: ctrl}
	mock.recorder = &MockCustomersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCustomers) EXPECT() *MockCustomersMockRecorder {
	return m.recorder
}

// CreateCustomer mocks base method
func (m *MockCustomers) CreateCustomer(arg0 context.Context, arg1, arg2 string, arg3 protoreflect.ProtoMessage) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCustomer", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCustomer indicates an expected call of CreateCustomer
func (mr *MockCustomersMockRecorder) CreateCustomer(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomer", reflect.TypeOf((*MockCustomers)(nil).CreateCustomer), arg0, arg1, arg2, arg3)
}

// GetCustomer mocks base method
func (m *MockCustomers) GetCustomer(arg0 context.Context, arg1 string, arg2 protoreflect.ProtoMessage) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomer", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCustomer indicates an expected call of GetCustomer
func (mr *MockCustomersMockRecorder) GetCustomer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomer", reflect.TypeOf((*MockCustomers)(nil).GetCustomer), arg0, arg1, arg2)
}

// ListCustomers mocks base method
func (m *MockCustomers) ListCustomers(arg0 context.Context, arg1 string, arg2 int) (string, []stores.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCustomers", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]stores.Record)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListCustomers indicates an expected call of ListCustomers
func (mr *MockCustomersMockRecorder) ListCustomers(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCustomers", reflect.TypeOf((*MockCustomers)(nil).ListCustomers), arg0, arg1, arg2)
}
