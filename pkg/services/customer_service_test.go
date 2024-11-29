package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repository_mocks "github.com/Adebayobenjamin/numerisbook/pkg/repositories/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupCustomerTest(t *testing.T) (*repository_mocks.MockCustomerRepository, *customerService) {
	ctrl := gomock.NewController(t)
	mockRepo := repository_mocks.NewMockCustomerRepository(ctrl)
	logger := zerolog.New(nil)
	service := NewCustomerService(&logger, mockRepo).(*customerService)
	return mockRepo, service
}

func TestGetCustomerByID(t *testing.T) {
	mockRepo, service := setupCustomerTest(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		customerID uint
		mockSetup  func() *models.Customer
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "successful customer retrieval",
			customerID: 1,
			mockSetup: func() *models.Customer {
				expected := &models.Customer{
					ID:      1,
					Name:    "John Doe",
					Email:   "john@example.com",
					Phone:   "1234567890",
					Address: "123 Main St",
				}
				mockRepo.EXPECT().
					GetCustomerByID(ctx, uint(1)).
					Return(expected, nil)
				return expected
			},
			wantErr: false,
		},
		{
			name:       "customer not found",
			customerID: 999,
			mockSetup: func() *models.Customer {
				mockRepo.EXPECT().
					GetCustomerByID(ctx, uint(999)).
					Return(nil, errors.New("customer not found"))
				return nil
			},
			wantErr: true,
			errMsg:  "customer not found",
		},
		{
			name:       "database error",
			customerID: 1,
			mockSetup: func() *models.Customer {
				mockRepo.EXPECT().
					GetCustomerByID(ctx, uint(1)).
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

			customer, err := service.GetCustomerByID(ctx, tt.customerID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, customer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, expected, customer)
			}
		})
	}
}

func TestNewCustomerService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := repository_mocks.NewMockCustomerRepository(ctrl)
	logger := zerolog.New(nil)

	service := NewCustomerService(&logger, mockRepo)

	assert.NotNil(t, service)

	// Type assertion to verify the concrete type
	_, ok := service.(*customerService)
	assert.True(t, ok, "service should be of type *customerService")

	// Verify internal fields are set correctly
	customerSvc := service.(*customerService)
	assert.Equal(t, &logger, customerSvc.logger)
	assert.Equal(t, mockRepo, customerSvc.customerRepository)
}
