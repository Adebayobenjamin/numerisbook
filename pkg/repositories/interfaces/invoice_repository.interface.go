package repositories_interfaces

import (
	"context"

	response_dto "github.com/Adebayobenjamin/numerisbook/pkg/dtos/response"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type InvoiceRepository interface {
	CreateInvoiceWithItems(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	GetByIDAndCutomerID(ctx context.Context, id, customerID uint) (*models.Invoice, error)
	UpdateShareableLink(ctx context.Context, invoiceID uint, link string) error
	GetStatistics(ctx context.Context, customerID uint) (*response_dto.GetInvoiceStatisticsResponse, error)
	GetDetails(ctx context.Context, invoiceID uint) (*response_dto.GetInvoiceDetailsResponse, error)
	DuplicateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error)
	GetAllCustomerInvoices(ctx context.Context, customerID uint, limit int, offset int) ([]models.Invoice, error)
	UpdateInvoiceStatus(ctx context.Context, invoiceID uint, status models.InvoiceStatus) error
}
