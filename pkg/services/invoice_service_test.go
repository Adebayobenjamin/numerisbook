package services

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	request_dto "github.com/Adebayobenjamin/numerisbook/pkg/dtos/request"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repository_mocks "github.com/Adebayobenjamin/numerisbook/pkg/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupInvoiceTest(t *testing.T) (*repository_mocks.MockInvoiceRepository, *repository_mocks.MockPaymentRepository, *invoiceService) {
	ctrl := gomock.NewController(t)
	mockInvoiceRepo := repository_mocks.NewMockInvoiceRepository(ctrl)
	mockPaymentRepo := repository_mocks.NewMockPaymentRepository(ctrl)
	service := NewInvoiceService(mockInvoiceRepo, mockPaymentRepo).(*invoiceService)
	return mockInvoiceRepo, mockPaymentRepo, service
}

func TestSetInvoiceStatusIfFullyPaid(t *testing.T) {
	mockInvoiceRepo, mockPaymentRepo, service := setupInvoiceTest(t)
	ctx := context.Background()

	tests := []struct {
		name           string
		invoice        *models.Invoice
		totalPayments  float64
		mockSetup      func()
		expectedStatus models.InvoiceStatus
		wantErr        bool
		errMsg         string
	}{
		{
			name: "fully paid invoice",
			invoice: &models.Invoice{
				ID:             1,
				TotalAmountDue: 100.0,
			},
			totalPayments: 100.0,
			mockSetup: func() {
				mockPaymentRepo.EXPECT().
					GetTotalInvoicePayments(ctx, uint(1)).
					Return(100.0, nil)
				mockInvoiceRepo.EXPECT().
					UpdateInvoiceStatus(ctx, uint(1), models.InvoiceStatusPaid).
					Return(nil)
			},
			expectedStatus: models.InvoiceStatusPaid,
			wantErr:        false,
		},
		{
			name: "partially paid invoice",
			invoice: &models.Invoice{
				ID:             1,
				TotalAmountDue: 100.0,
			},
			totalPayments: 50.0,
			mockSetup: func() {
				mockPaymentRepo.EXPECT().
					GetTotalInvoicePayments(ctx, uint(1)).
					Return(50.0, nil)
			},
			wantErr: false,
		},
		{
			name: "error getting payments",
			invoice: &models.Invoice{
				ID:             1,
				TotalAmountDue: 100.0,
			},
			mockSetup: func() {
				mockPaymentRepo.EXPECT().
					GetTotalInvoicePayments(ctx, uint(1)).
					Return(0.0, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "failed to get total invoice payments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.SetInvoiceStatusIfFullyPaid(ctx, tt.invoice)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				if tt.totalPayments >= tt.invoice.TotalAmountDue {
					assert.Equal(t, tt.expectedStatus, tt.invoice.Status)
				}
			}
		})
	}
}

func TestCreateInvoice(t *testing.T) {
	mockInvoiceRepo, _, service := setupInvoiceTest(t)
	ctx := context.Background()

	validRequest := &request_dto.CreateInvoiceRequest{
		DueDate: time.Now().Add(24 * time.Hour),
		Items: []request_dto.InvoiceItem{
			{UnitPrice: 100.0, Quantity: 2},
			{UnitPrice: 50.0, Quantity: 1},
		},
		Discount: 20.0,
	}

	tests := []struct {
		name       string
		customerID uint
		request    *request_dto.CreateInvoiceRequest
		mockSetup  func()
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "successful creation",
			customerID: 1,
			request:    validRequest,
			mockSetup: func() {
				mockInvoiceRepo.EXPECT().
					CreateInvoiceWithItems(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, invoice *models.Invoice) (*models.Invoice, error) {
						assert.Equal(t, float64(250), invoice.Subtotal)       // (100*2 + 50*1)
						assert.Equal(t, float64(230), invoice.TotalAmountDue) // 250 - 20 (discount)
						assert.NotEmpty(t, invoice.InvoiceNumber)
						return invoice, nil
					})
			},
			wantErr: false,
		},
		{
			name:       "past due date",
			customerID: 1,
			request: &request_dto.CreateInvoiceRequest{
				DueDate: time.Now().Add(-24 * time.Hour),
			},
			mockSetup: func() {},
			wantErr:   true,
			errMsg:    "due date cannot be in the past",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			invoice, err := service.CreateInvoice(ctx, tt.customerID, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, invoice)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, invoice)
			}
		})
	}
}

func TestGetCustomerInvoices(t *testing.T) {
	mockInvoiceRepo, _, service := setupInvoiceTest(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		limit      int
		page       int
		customerID uint
		mockSetup  func()
		wantErr    bool
	}{
		{
			name:       "successful retrieval",
			limit:      10,
			page:       1,
			customerID: 1,
			mockSetup: func() {
				mockInvoiceRepo.EXPECT().
					GetAllCustomerInvoices(ctx, uint(1), 10, 0).
					Return([]models.Invoice{{ID: 1}, {ID: 2}}, nil)
			},
			wantErr: false,
		},
		{
			name:       "repository error",
			limit:      10,
			page:       1,
			customerID: 1,
			mockSetup: func() {
				mockInvoiceRepo.EXPECT().
					GetAllCustomerInvoices(ctx, uint(1), 10, 0).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			invoices, err := service.GetCustomerInvoices(ctx, tt.limit, tt.page, tt.customerID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, invoices)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, invoices)
				assert.Len(t, invoices, 2)
			}
		})
	}
}

func TestValidatePaymentAmount(t *testing.T) {
	_, mockPaymentRepo, service := setupInvoiceTest(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		amount        float64
		invoice       *models.Invoice
		isPartial     bool
		totalPayments float64
		mockSetup     func()
		wantErr       bool
		errMsg        string
	}{
		{
			name:   "valid full payment",
			amount: 100.0,
			invoice: &models.Invoice{
				ID:             1,
				TotalAmountDue: 100.0,
			},
			isPartial: false,
			mockSetup: func() {
				mockPaymentRepo.EXPECT().
					GetTotalInvoicePayments(ctx, uint(1)).
					Return(0.0, nil)
			},
			wantErr: false,
		},
		{
			name:   "valid partial payment",
			amount: 50.0,
			invoice: &models.Invoice{
				ID:             1,
				TotalAmountDue: 100.0,
			},
			isPartial: true,
			mockSetup: func() {
				mockPaymentRepo.EXPECT().
					GetTotalInvoicePayments(ctx, uint(1)).
					Return(0.0, nil)
			},
			wantErr: false,
		},
		{
			name:   "payment exceeds total",
			amount: 150.0,
			invoice: &models.Invoice{
				ID:             1,
				TotalAmountDue: 100.0,
			},
			isPartial: false,
			mockSetup: func() {
				mockPaymentRepo.EXPECT().
					GetTotalInvoicePayments(ctx, uint(1)).
					Return(0.0, nil)
			},
			wantErr: true,
			errMsg:  "payment amount exceeds invoice total amount",
		},
		{
			name:   "insufficient non-partial payment",
			amount: 50.0,
			invoice: &models.Invoice{
				ID:             1,
				TotalAmountDue: 100.0,
			},
			isPartial: false,
			mockSetup: func() {
				mockPaymentRepo.EXPECT().
					GetTotalInvoicePayments(ctx, uint(1)).
					Return(0.0, nil)
			},
			wantErr: true,
			errMsg:  "payment amount is less than invoice total amount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.ValidatePaymentAmount(ctx, tt.amount, tt.invoice, tt.isPartial)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetShareableLink(t *testing.T) {
	mockInvoiceRepo, _, service := setupInvoiceTest(t)
	ctx := context.Background()

	originalURL := os.Getenv("FRONTEND_URL")
	os.Setenv("FRONTEND_URL", "http://example.com")
	defer os.Setenv("FRONTEND_URL", originalURL)

	tests := []struct {
		name      string
		invoice   *models.Invoice
		mockSetup func()
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "successful link generation",
			invoice: &models.Invoice{ID: 1},
			mockSetup: func() {
				mockInvoiceRepo.EXPECT().
					UpdateShareableLink(ctx, uint(1), "http://example.com/invoice/1").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "update error",
			invoice: &models.Invoice{ID: 1},
			mockSetup: func() {
				mockInvoiceRepo.EXPECT().
					UpdateShareableLink(ctx, uint(1), "http://example.com/invoice/1").
					Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "failed to update shareable link",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			link, err := service.GetShareableLink(ctx, tt.invoice)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Empty(t, link)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "http://example.com/invoice/1", link)
			}
		})
	}
}
