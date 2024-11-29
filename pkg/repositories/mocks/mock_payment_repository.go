// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/repositories/interfaces/payment_repository.interface.go
//
// Generated by this command:
//
//	mockgen -source=pkg/repositories/interfaces/payment_repository.interface.go -destination=pkg/repositories/mocks/mock_payment_repository.go -package=repository_mocks
//

// Package repository_mocks is a generated GoMock package.
package repository_mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/Adebayobenjamin/numerisbook/pkg/models"
	gomock "go.uber.org/mock/gomock"
)

// MockPaymentRepository is a mock of PaymentRepository interface.
type MockPaymentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentRepositoryMockRecorder
	isgomock struct{}
}

// MockPaymentRepositoryMockRecorder is the mock recorder for MockPaymentRepository.
type MockPaymentRepositoryMockRecorder struct {
	mock *MockPaymentRepository
}

// NewMockPaymentRepository creates a new mock instance.
func NewMockPaymentRepository(ctrl *gomock.Controller) *MockPaymentRepository {
	mock := &MockPaymentRepository{ctrl: ctrl}
	mock.recorder = &MockPaymentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPaymentRepository) EXPECT() *MockPaymentRepositoryMockRecorder {
	return m.recorder
}

// CreatePayment mocks base method.
func (m *MockPaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePayment", ctx, payment)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePayment indicates an expected call of CreatePayment.
func (mr *MockPaymentRepositoryMockRecorder) CreatePayment(ctx, payment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePayment", reflect.TypeOf((*MockPaymentRepository)(nil).CreatePayment), ctx, payment)
}

// GetTotalInvoicePayments mocks base method.
func (m *MockPaymentRepository) GetTotalInvoicePayments(ctx context.Context, invoiceID uint) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalInvoicePayments", ctx, invoiceID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalInvoicePayments indicates an expected call of GetTotalInvoicePayments.
func (mr *MockPaymentRepositoryMockRecorder) GetTotalInvoicePayments(ctx, invoiceID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalInvoicePayments", reflect.TypeOf((*MockPaymentRepository)(nil).GetTotalInvoicePayments), ctx, invoiceID)
}
