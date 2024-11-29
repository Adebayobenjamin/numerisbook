package services

import (
	"context"
	"fmt"
	"os"
	"time"

	request_dto "github.com/Adebayobenjamin/numerisbook/pkg/dtos/request"
	response_dto "github.com/Adebayobenjamin/numerisbook/pkg/dtos/response"
	"github.com/Adebayobenjamin/numerisbook/pkg/helper"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repositories_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/repositories/interfaces"
	services_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/services/interfaces"
)

type invoiceService struct {
	invoiceRepository repositories_interfaces.InvoiceRepository
	paymentRepository repositories_interfaces.PaymentRepository
}

// SetInvoiceStatusIfFullyPaid implements services_interfaces.InvoiceService.
func (i *invoiceService) SetInvoiceStatusIfFullyPaid(ctx context.Context, invoice *models.Invoice) error {
	totalPayments, err := i.paymentRepository.GetTotalInvoicePayments(ctx, invoice.ID)
	if err != nil {
		return fmt.Errorf("failed to get total invoice payments: %w", err)
	}

	// Check if total payments equal the total amount due
	if totalPayments >= invoice.TotalAmountDue {
		invoice.Status = models.InvoiceStatusPaid
		err = i.invoiceRepository.UpdateInvoiceStatus(ctx, invoice.ID, models.InvoiceStatusPaid)
		if err != nil {
			return fmt.Errorf("failed to update invoice status: %w", err)
		}
	}

	return nil
}

// ConfirmPayment implements services_interfaces.InvoiceService.
func (i *invoiceService) ConfirmPayment(ctx context.Context, invoiceID uint, amount float64, date time.Time, isPartial bool) error {
	payment := &models.Payment{
		InvoiceID: invoiceID,
		Amount:    amount,
		IsPartial: isPartial,
		Date:      date,
	}

	err := i.paymentRepository.CreatePayment(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// CreateInvoice implements services_interfaces.InvoiceService.
func (i *invoiceService) CreateInvoice(ctx context.Context, customerID uint, request *request_dto.CreateInvoiceRequest) (*models.Invoice, error) {
	var invoiceToBeCreated models.Invoice

	err := helper.JSONUnmarshalToType(request, &invoiceToBeCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal invoice: %w", err)
	}

	if invoiceToBeCreated.DueDate.Before(time.Now()) {
		return nil, fmt.Errorf("due date cannot be in the past")
	}

	invoiceToBeCreated.InvoiceNumber = helper.GenerateInvoiceNumber()
	invoiceToBeCreated.CustomerID = customerID

	for _, item := range invoiceToBeCreated.Items {
		item.TotalPrice = item.UnitPrice * float64(item.Quantity)
		invoiceToBeCreated.TotalAmountDue += item.TotalPrice
		invoiceToBeCreated.Subtotal += item.TotalPrice
	}

	invoiceToBeCreated.TotalAmountDue -= invoiceToBeCreated.Discount

	invoice, err := i.invoiceRepository.CreateInvoiceWithItems(ctx, &invoiceToBeCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	return invoice, nil
}

// DuplicateInvoice implements services_interfaces.InvoiceService.
func (i *invoiceService) DuplicateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	return i.invoiceRepository.DuplicateInvoice(ctx, invoice)
}

// GetCustomerInvoices implements services_interfaces.InvoiceService.
func (i *invoiceService) GetCustomerInvoices(ctx context.Context, limit int, page int, customerID uint) ([]models.Invoice, error) {
	offset := helper.GetOffset(page, limit)
	return i.invoiceRepository.GetAllCustomerInvoices(ctx, customerID, limit, offset)
}

// GetInvoiceByIDandCustomer implements services_interfaces.InvoiceService.
func (i *invoiceService) GetInvoiceByIDandCustomer(ctx context.Context, invoiceID uint, customerID uint) (*models.Invoice, error) {
	return i.invoiceRepository.GetByIDAndCutomerID(ctx, invoiceID, customerID)
}

// GetInvoiceDetails implements services_interfaces.InvoiceService.
func (i *invoiceService) GetInvoiceDetails(ctx context.Context, invoiceID uint) (*response_dto.GetInvoiceDetailsResponse, error) {
	return i.invoiceRepository.GetDetails(ctx, invoiceID)
}

// GetInvoiceStatistics implements services_interfaces.InvoiceService.
func (i *invoiceService) GetInvoiceStatistics(ctx context.Context, customerID uint) (*response_dto.GetInvoiceStatisticsResponse, error) {
	return i.invoiceRepository.GetStatistics(ctx, customerID)
}

// GetShareableLink implements services_interfaces.InvoiceService.
func (i *invoiceService) GetShareableLink(ctx context.Context, invoice *models.Invoice) (string, error) {
	link := fmt.Sprintf("%s/invoice/%d", os.Getenv("FRONTEND_URL"), invoice.ID)

	err := i.invoiceRepository.UpdateShareableLink(ctx, invoice.ID, link)
	if err != nil {
		return "", fmt.Errorf("failed to update shareable link: %w", err)
	}

	return link, nil
}

// ValidatePaymentAmount implements services_interfaces.InvoiceService.
func (i *invoiceService) ValidatePaymentAmount(ctx context.Context, amount float64, invoice *models.Invoice, isPartial bool) error {
	totalPayments, err := i.paymentRepository.GetTotalInvoicePayments(ctx, invoice.ID)

	if err != nil {
		return fmt.Errorf("failed to get total invoice payments: %w", err)
	}

	if totalPayments+amount > invoice.TotalAmountDue {
		return fmt.Errorf("payment amount exceeds invoice total amount")
	}

	if !isPartial && totalPayments+amount < invoice.TotalAmountDue {
		return fmt.Errorf("payment amount is less than invoice total amount")
	}

	return nil
}

func NewInvoiceService(
	invoiceRepository repositories_interfaces.InvoiceRepository,
	paymentRepository repositories_interfaces.PaymentRepository,
) services_interfaces.InvoiceService {
	return &invoiceService{
		invoiceRepository: invoiceRepository,
		paymentRepository: paymentRepository,
	}
}
