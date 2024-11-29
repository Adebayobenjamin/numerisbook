package repositories

import (
	"context"
	"database/sql"
	"fmt"

	response_dto "github.com/Adebayobenjamin/numerisbook/pkg/dtos/response"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repositories_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/repositories/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type invoiceRepository struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

// UpdateInvoiceStatus implements repositories_interfaces.InvoiceRepository.
func (i *invoiceRepository) UpdateInvoiceStatus(ctx context.Context, invoiceID uint, status models.InvoiceStatus) error {
	query := `
		UPDATE invoices 
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`

	_, err := i.db.ExecContext(ctx, query, status, invoiceID)
	if err != nil {
		return fmt.Errorf("failed to update invoice status: %w", err)
	}

	return nil
}

// CreateInvoiceWithItems implements repositories_interfaces.InvoiceRepository.
func (i *invoiceRepository) CreateInvoiceWithItems(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	tx, err := i.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert invoice first
	invoiceQuery := `
		INSERT INTO invoices (
			invoice_number, customer_id, issue_date, due_date,
			total_amount_due, subtotal, is_fully_paid, billing_currency,
			discount, status, notes, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	invoiceResult, err := tx.ExecContext(ctx, invoiceQuery,
		invoice.InvoiceNumber,
		invoice.CustomerID,
		invoice.IssueDate,
		invoice.DueDate,
		invoice.TotalAmountDue,
		invoice.Subtotal,
		invoice.IsFullyPaid,
		invoice.BillingCurrency,
		invoice.Discount,
		models.InvoiceStatusPendingPayment,
		invoice.Notes)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	invoiceID, _ := invoiceResult.LastInsertId()

	// Insert sender
	senderQuery := `
		INSERT INTO senders (name, phone, address, email, invoice_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err = tx.ExecContext(ctx, senderQuery,
		invoice.Sender.Name,
		invoice.Sender.Phone,
		invoice.Sender.Address,
		invoice.Sender.Email,
		invoiceID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sender: %w", err)
	}

	// Insert invoice items
	itemQuery := `
		INSERT INTO invoice_items (
			invoice_id, description, quantity, unit_price, total_price,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	for _, item := range invoice.Items {
		_, err = tx.ExecContext(ctx, itemQuery,
			invoiceID,
			item.Description,
			item.Quantity,
			item.UnitPrice,
			item.TotalPrice)
		if err != nil {
			return nil, fmt.Errorf("failed to create invoice item: %w", err)
		}
	}

	// Insert payment info
	paymentInfoQuery := `
		INSERT INTO payment_info (invoice_id, bank_name, account_number, account_name, ach_routing_no, bank_address, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err = tx.ExecContext(ctx, paymentInfoQuery,
		invoiceID,
		invoice.PaymentInfo.BankName,
		invoice.PaymentInfo.AccountNumber,
		invoice.PaymentInfo.AccountName,
		invoice.PaymentInfo.AchRoutingNo,
		invoice.PaymentInfo.BankAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment info: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Fetch the created invoice
	return i.GetByIDAndCutomerID(ctx, uint(invoiceID), invoice.CustomerID)
}

// DuplicateInvoice implements repositories_interfaces.InvoiceRepository.
func (i *invoiceRepository) DuplicateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	tx, err := i.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert new invoice
	invoiceQuery := `
		INSERT INTO invoices (
			invoice_number, customer_id, issue_date, due_date,
			total_amount_due, subtotal, is_fully_paid, billing_currency,
			discount, status, notes, created_at, updated_at
		)
		SELECT 
			CONCAT(invoice_number, '-copy'), customer_id, issue_date, due_date,
			total_amount_due, subtotal, FALSE, billing_currency,
			discount, 'draft', notes, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		FROM invoices 
		WHERE id = ? AND deleted_at IS NULL`

	invoiceResult, err := tx.ExecContext(ctx, invoiceQuery, invoice.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to duplicate invoice: %w", err)
	}

	newInvoiceID, _ := invoiceResult.LastInsertId()

	// Duplicate sender
	senderQuery := `
		INSERT INTO senders (name, phone, address, email, invoice_id, created_at, updated_at)
		SELECT name, phone, address, email, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		FROM senders
		WHERE invoice_id = ?`

	_, err = tx.ExecContext(ctx, senderQuery, newInvoiceID, invoice.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to duplicate sender: %w", err)
	}

	// Duplicate invoice items
	itemsQuery := `
		INSERT INTO invoice_items (
			invoice_id, description, quantity, unit_price, total_price,
			created_at, updated_at
		)
		SELECT 
			?, description, quantity, unit_price, total_price,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		FROM invoice_items
		WHERE invoice_id = ? AND deleted_at IS NULL`

	_, err = tx.ExecContext(ctx, itemsQuery, newInvoiceID, invoice.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to duplicate invoice items: %w", err)
	}

	// Duplicate payment info
	paymentInfoQuery := `
		INSERT INTO payment_info (
			invoice_id, bank_name, account_number, account_name,
			ach_routing_no, bank_address, created_at, updated_at
		)
		SELECT 
			?, bank_name, account_number, account_name,
			ach_routing_no, bank_address, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		FROM payment_info
		WHERE invoice_id = ? AND deleted_at IS NULL`

	_, err = tx.ExecContext(ctx, paymentInfoQuery, newInvoiceID, invoice.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to duplicate payment info: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Get the newly created invoice with all its relations
	return i.GetByIDAndCutomerID(ctx, uint(newInvoiceID), invoice.CustomerID)
}

// GetAllCustomerInvoices implements repositories_interfaces.InvoiceRepository.
func (i *invoiceRepository) GetAllCustomerInvoices(ctx context.Context, customerID uint, limit int, offset int) ([]models.Invoice, error) {
	query := `
		SELECT 
			invoice_number,
			issue_date,
			due_date,
			total_amount_due,
			subtotal,
			status
		FROM invoices
		WHERE customer_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	var invoices []models.Invoice
	err := i.db.SelectContext(ctx, &invoices, query, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer invoices: %w", err)
	}

	return invoices, nil
}

// GetByIDAndCutomerID implements repositories_interfaces.InvoiceRepository.
func (i *invoiceRepository) GetByIDAndCutomerID(ctx context.Context, id uint, customerID uint) (*models.Invoice, error) {
	query := `
		SELECT 
			i.*,
			s.id as "sender.id",
			s.name as "sender.name",
			s.phone as "sender.phone",
			s.address as "sender.address",
			s.email as "sender.email",
			c.id as "customer.id",
			c.name as "customer.name",
			c.phone as "customer.phone",
			c.address as "customer.address",
			c.email as "customer.email"
		FROM invoices i
		LEFT JOIN senders s ON i.id = s.invoice_id
		LEFT JOIN customers c ON i.customer_id = c.id
		WHERE i.id = ? AND i.customer_id = ? AND i.deleted_at IS NULL`

	var invoice models.Invoice
	err := i.db.GetContext(ctx, &invoice, query, id, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// Get invoice items
	itemsQuery := `
		SELECT * FROM invoice_items 
		WHERE invoice_id = ? AND deleted_at IS NULL`

	if err := i.db.SelectContext(ctx, &invoice.Items, itemsQuery, id); err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}

	return &invoice, nil
}

// GetDetails implements repositories_interfaces.InvoiceRepository.
func (i *invoiceRepository) GetDetails(ctx context.Context, invoiceID uint) (*response_dto.GetInvoiceDetailsResponse, error) {
	// Get main invoice details with sender, customer and payment info
	query := `
		SELECT 
			i.*,
			s.id as "sender.id",
			s.name as "sender.name",
			s.phone as "sender.phone",
			s.address as "sender.address",
			s.email as "sender.email",
			c.id as "customer.id",
			c.name as "customer.name",
			c.phone as "customer.phone",
			c.address as "customer.address",
			c.email as "customer.email",
			p.id as "payment_information.id",
			p.bank_name as "payment_information.bank_name",
			p.account_number as "payment_information.account_number",
			p.account_name as "payment_information.account_name",
			p.ach_routing_no as "payment_information.ach_routing_no",
			p.bank_address as "payment_information.bank_address"
		FROM invoices i
		LEFT JOIN senders s ON i.id = s.invoice_id
		LEFT JOIN customers c ON i.customer_id = c.id
		LEFT JOIN payment_info p ON i.id = p.invoice_id
		WHERE i.id = ? AND i.deleted_at IS NULL`

	details := &response_dto.GetInvoiceDetailsResponse{}
	err := i.db.GetContext(ctx, details, query, invoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice details: %w", err)
	}

	// Get invoice items
	itemsQuery := `
		SELECT * FROM invoice_items 
		WHERE invoice_id = ? AND deleted_at IS NULL`
	if err := i.db.SelectContext(ctx, &details.Items, itemsQuery, invoiceID); err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}

	// Get payments
	paymentsQuery := `
		SELECT * FROM payments 
		WHERE invoice_id = ? AND deleted_at IS NULL 
		ORDER BY date DESC`
	if err := i.db.SelectContext(ctx, &details.Payments, paymentsQuery, invoiceID); err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	// Get reminders
	remindersQuery := `
		SELECT * FROM invoice_reminders 
		WHERE invoice_id = ? AND deleted_at IS NULL 
		ORDER BY reminder_date ASC`
	if err := i.db.SelectContext(ctx, &details.Reminders, remindersQuery, invoiceID); err != nil {
		return nil, fmt.Errorf("failed to get reminders: %w", err)
	}

	return details, nil
}

// GetStatistics implements repositories_interfaces.InvoiceRepository.
func (i *invoiceRepository) GetStatistics(ctx context.Context, customerID uint) (*response_dto.GetInvoiceStatisticsResponse, error) {
	query := `
		SELECT
			SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END) as total_paid,
			SUM(CASE WHEN status = 'paid' THEN total_amount_due ELSE 0 END) as total_paid_amount,
			SUM(CASE WHEN status != 'paid' AND due_date < CURRENT_TIMESTAMP THEN 1 ELSE 0 END) as total_over_due,
			SUM(CASE WHEN status != 'paid' AND due_date < CURRENT_TIMESTAMP THEN total_amount_due ELSE 0 END) as total_over_due_amount,
			SUM(CASE WHEN status = 'draft' THEN 1 ELSE 0 END) as total_draft,
			SUM(CASE WHEN status = 'draft' THEN total_amount_due ELSE 0 END) as total_draft_amount,
			SUM(CASE WHEN status = 'pending payment' THEN 1 ELSE 0 END) as total_unpaid,
			SUM(CASE WHEN status = 'pending payment' THEN total_amount_due ELSE 0 END) as total_unpaid_amount
		FROM invoices
		WHERE customer_id = ? AND deleted_at IS NULL`

	stats := &response_dto.GetInvoiceStatisticsResponse{}
	err := i.db.GetContext(ctx, stats, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice statistics: %w", err)
	}

	return stats, nil
}

// UpdateShareableLink implements repositories_interfaces.InvoiceRepository.
func (i *invoiceRepository) UpdateShareableLink(ctx context.Context, invoiceID uint, link string) error {
	query := `
		UPDATE invoices 
		SET shareable_link = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL`

	result, err := i.db.ExecContext(ctx, query, link, invoiceID)
	if err != nil {
		return fmt.Errorf("failed to update shareable link: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("invoice not found")
	}

	return nil
}

func NewInvoiceRepository(
	db *sqlx.DB,
	logger *zerolog.Logger,
) repositories_interfaces.InvoiceRepository {
	return &invoiceRepository{
		db:     db,
		logger: logger,
	}
}
