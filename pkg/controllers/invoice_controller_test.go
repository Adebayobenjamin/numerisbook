package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	request_dto "github.com/Adebayobenjamin/numerisbook/pkg/dtos/request"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	services_mocks "github.com/Adebayobenjamin/numerisbook/pkg/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestConfirmPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock services
	mockInvoiceService := services_mocks.NewMockInvoiceService(ctrl)
	mockAuditService := services_mocks.NewMockAuditService(ctrl)
	mockReminderService := services_mocks.NewMockRemiderService(ctrl)
	mockCustomerService := services_mocks.NewMockCustomerService(ctrl)

	logger := zerolog.New(nil)
	controller := NewInvoiceController(&logger, mockInvoiceService, mockAuditService, mockReminderService, mockCustomerService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/invoices/:invoice_id/confirm-payment", controller.ConfirmPayment)

	// Test data
	customer := &models.Customer{ID: 1, Name: "John Doe"}
	invoice := &models.Invoice{ID: 1, InvoiceNumber: "INV-001", Status: "Pending"}
	requestBody := request_dto.PaymentConfirmationRequest{
		Amount:      100.00,
		PaymentDate: time.Now(),
		IsPartial:   false,
	}
	body, _ := json.Marshal(requestBody)

	// Mock context extraction
	mockCustomerService.EXPECT().
		GetCustomerByID(gomock.Any(), gomock.Any()).
		Return(customer, nil).
		AnyTimes()

	mockInvoiceService.EXPECT().
		GetInvoiceByIDandCustomer(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(invoice, nil).
		AnyTimes()

	mockInvoiceService.EXPECT().
		ValidatePaymentAmount(gomock.Any(), requestBody.Amount, invoice, requestBody.IsPartial).
		Return(nil).
		Times(1)

	mockInvoiceService.EXPECT().
		ConfirmPayment(gomock.Any(), invoice.ID, requestBody.Amount, requestBody.PaymentDate, requestBody.IsPartial).
		Return(nil).
		Times(1)

	mockInvoiceService.EXPECT().
		SetInvoiceStatusIfFullyPaid(gomock.Any(), invoice).
		Return(nil).
		Times(1)

	mockAuditService.EXPECT().
		CreateAuditTrail(gomock.Any(), models.EventTypePaymentConfirmed, models.LogLevelInfo, gomock.Any(), invoice.ID, customer.ID).
		Times(1)

	// Make HTTP request
	req := httptest.NewRequest(http.MethodPost, "/invoices/1/confirm-payment", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "payment confirmed successfully")
}
