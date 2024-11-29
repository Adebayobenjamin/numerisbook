package repositories_interfaces

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetTotalInvoicePayments(ctx context.Context, invoiceID uint) (float64, error)
}
