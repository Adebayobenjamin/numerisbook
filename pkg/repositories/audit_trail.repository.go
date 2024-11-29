package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repositories_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/repositories/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type auditTrailRepository struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

// GetAllCustomerAuditTrails retrieves all audit trails for a specific customer with pagination
func (a *auditTrailRepository) GetAllCustomerAuditTrails(ctx context.Context, customerID uint, limit int, offset int) ([]models.AuditTrail, error) {
	query := `
        SELECT * FROM audit_trails 
        WHERE customer_id = ? AND deleted_at IS NULL 
        ORDER BY created_at DESC 
        LIMIT ? OFFSET ?`

	var auditTrails []models.AuditTrail
	err := a.db.SelectContext(ctx, &auditTrails, query, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer audit trails: %w", err)
	}

	return auditTrails, nil
}

// GetByInvoiceIDAndCustomerID retrieves a specific audit trail by invoice and customer IDs
func (a *auditTrailRepository) GetByInvoiceIDAndCustomerID(ctx context.Context, invoiceID uint, customerID uint, limit int, offset int) ([]models.AuditTrail, error) {
	query := `
        SELECT * FROM audit_trails 
        WHERE invoice_id = ? AND customer_id = ? AND deleted_at IS NULL 
        ORDER BY created_at DESC 
        LIMIT ? OFFSET ?	`

	var auditTrails []models.AuditTrail
	err := a.db.SelectContext(ctx, &auditTrails, query, invoiceID, customerID, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.AuditTrail{}, fmt.Errorf("audit trail not found")
		}
		return []models.AuditTrail{}, fmt.Errorf("failed to get audit trail: %w", err)
	}

	return auditTrails, nil
}

// LogEvent creates a new audit trail entry
func (a *auditTrailRepository) LogEvent(ctx context.Context, eventType models.EventType, logLevel models.LogLevel, message string, invoiceID uint, customerID uint) error {
	query := `
        INSERT INTO audit_trails (
            event_type,
            log_level,
            message,
            invoice_id,
            customer_id,
            created_at,
            updated_at
        ) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := a.db.ExecContext(ctx, query, eventType, logLevel, message, invoiceID, customerID)
	if err != nil {
		return fmt.Errorf("failed to log audit trail event: %w", err)
	}

	return nil
}

func NewAuditTrailRepository(
	db *sqlx.DB,
	logger *zerolog.Logger,
) repositories_interfaces.AuditTrailRepository {
	return &auditTrailRepository{
		db:     db,
		logger: logger,
	}
}
