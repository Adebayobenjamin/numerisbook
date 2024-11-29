package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repository_mocks "github.com/Adebayobenjamin/numerisbook/pkg/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupAuditTest(t *testing.T) (*repository_mocks.MockAuditTrailRepository, *auditService) {
	ctrl := gomock.NewController(t)
	mockRepo := repository_mocks.NewMockAuditTrailRepository(ctrl)
	service := NewAuditService(mockRepo).(*auditService)
	return mockRepo, service
}

func TestCreateAuditTrail(t *testing.T) {
	mockRepo, service := setupAuditTest(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		eventType  models.EventType
		logLevel   models.LogLevel
		message    string
		invoiceID  uint
		customerID uint
		mockSetup  func()
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "successful audit creation",
			eventType:  models.EventType("INVOICE_CREATED"),
			logLevel:   models.LogLevel("INFO"),
			message:    "Test message",
			invoiceID:  1,
			customerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().
					LogEvent(ctx, models.EventType("INVOICE_CREATED"), models.LogLevel("INFO"), "Test message", uint(1), uint(1)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "repository error",
			eventType:  models.EventType("INVOICE_CREATED"),
			logLevel:   models.LogLevel("INFO"),
			message:    "Test message",
			invoiceID:  1,
			customerID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().
					LogEvent(ctx, models.EventType("INVOICE_CREATED"), models.LogLevel("INFO"), "Test message", uint(1), uint(1)).
					Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.CreateAuditTrail(ctx, tt.eventType, tt.logLevel, tt.message, tt.invoiceID, tt.customerID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAuditTrailsByInvoiceID(t *testing.T) {
	mockRepo, service := setupAuditTest(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		invoiceID  uint
		customerID uint
		limit      int
		page       int
		mockSetup  func() []models.AuditTrail
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "successful retrieval",
			invoiceID:  1,
			customerID: 1,
			limit:      10,
			page:       1,
			mockSetup: func() []models.AuditTrail {
				expected := []models.AuditTrail{{ID: 1, InvoiceID: 1, CustomerID: 1}}
				mockRepo.EXPECT().
					GetByInvoiceIDAndCustomerID(ctx, uint(1), uint(1), 10, 0).
					Return(expected, nil)
				return expected
			},
			wantErr: false,
		},
		{
			name:       "repository error",
			invoiceID:  1,
			customerID: 1,
			limit:      10,
			page:       1,
			mockSetup: func() []models.AuditTrail {
				mockRepo.EXPECT().
					GetByInvoiceIDAndCustomerID(ctx, uint(1), uint(1), 10, 0).
					Return(nil, errors.New("database error"))
				return nil
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tt.mockSetup()

			trails, err := service.GetAuditTrailsByInvoiceID(ctx, tt.invoiceID, tt.customerID, tt.limit, tt.page)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expected, trails)
			}
		})
	}
}

func TestGetCustomerAuditTrails(t *testing.T) {
	mockRepo, service := setupAuditTest(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		customerID uint
		limit      int
		page       int
		mockSetup  func() []models.AuditTrail
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "successful retrieval",
			customerID: 1,
			limit:      10,
			page:       1,
			mockSetup: func() []models.AuditTrail {
				expected := []models.AuditTrail{{ID: 1, CustomerID: 1}}
				mockRepo.EXPECT().
					GetAllCustomerAuditTrails(ctx, uint(1), 10, 0).
					Return(expected, nil)
				return expected
			},
			wantErr: false,
		},
		{
			name:       "repository error",
			customerID: 1,
			limit:      10,
			page:       1,
			mockSetup: func() []models.AuditTrail {
				mockRepo.EXPECT().
					GetAllCustomerAuditTrails(ctx, uint(1), 10, 0).
					Return(nil, errors.New("database error"))
				return nil
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tt.mockSetup()

			trails, err := service.GetCustomerAuditTrails(ctx, tt.customerID, tt.limit, tt.page)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expected, trails)
			}
		})
	}
}
