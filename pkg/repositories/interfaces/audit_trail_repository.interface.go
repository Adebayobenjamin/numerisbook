package repositories_interfaces

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type AuditTrailRepository interface {
	LogEvent(ctx context.Context, eventType models.EventType, logLevel models.LogLevel, message string, invoiceID uint, customerID uint) error
	GetAllCustomerAuditTrails(ctx context.Context, customerID uint, limit int, offset int) ([]models.AuditTrail, error)
	GetByInvoiceIDAndCustomerID(ctx context.Context, invoiceID uint, customerID uint, limit int, offset int) ([]models.AuditTrail, error)
}
