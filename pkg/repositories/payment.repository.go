package repositories

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repositories_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/repositories/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type paymentRepository struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

// CreatePayment implements repositories_interfaces.PaymentRepository.
func (p *paymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	query := `INSERT INTO payments (invoice_id, amount, is_partial, date, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := p.db.ExecContext(ctx, query, payment.InvoiceID, payment.Amount, payment.IsPartial, payment.Date)
	if err != nil {
		return err
	}
	return nil
}

// GetTotalInvoicePayments implements repositories_interfaces.PaymentRepository.
func (p *paymentRepository) GetTotalInvoicePayments(ctx context.Context, invoiceID uint) (float64, error) {
	query := `SELECT SUM(amount) as amount FROM payments WHERE invoice_id = ?`

	var totalPayments *float64
	err := p.db.GetContext(ctx, &totalPayments, query, invoiceID)
	if err != nil || totalPayments == nil {
		return 0, err
	}

	return *totalPayments, nil
}

func NewPaymentRepository(
	db *sqlx.DB, logger *zerolog.Logger,
) repositories_interfaces.PaymentRepository {
	return &paymentRepository{
		db:     db,
		logger: logger,
	}
}
