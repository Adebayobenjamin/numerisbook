package services_interfaces

import (
	"context"
	"time"

	request_dto "github.com/Adebayobenjamin/numerisbook/pkg/dtos/request"
	response_dto "github.com/Adebayobenjamin/numerisbook/pkg/dtos/response"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type InvoiceService interface {
	CreateInvoice(ctx context.Context, customerID uint, request *request_dto.CreateInvoiceRequest) (*models.Invoice, error)
	DuplicateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	GetInvoiceByIDandCustomer(ctx context.Context, invoiceID uint, customerID uint) (*models.Invoice, error)
	ConfirmPayment(ctx context.Context, invoiceID uint, amount float64, date time.Time, isPartial bool) error
	ValidatePaymentAmount(ctx context.Context, amount float64, invoice *models.Invoice, isPartial bool) error
	GetInvoiceDetails(ctx context.Context, invoiceID uint) (*response_dto.GetInvoiceDetailsResponse, error)
	GetCustomerInvoices(ctx context.Context, limit int, page int, customerID uint) ([]models.Invoice, error)
	GetShareableLink(ctx context.Context, invoice *models.Invoice) (string, error)
	GetInvoiceStatistics(ctx context.Context, customerID uint) (*response_dto.GetInvoiceStatisticsResponse, error)
	SetInvoiceStatusIfFullyPaid(ctx context.Context, invoice *models.Invoice) error
}
