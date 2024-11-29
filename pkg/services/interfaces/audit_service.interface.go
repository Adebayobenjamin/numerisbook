package services_interfaces

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type AuditService interface {
	CreateAuditTrail(
		ctx context.Context,
		eventType models.EventType,
		logLevel models.LogLevel,
		message string,
		invoiceID uint,
		customerID uint,
	) error

	GetCustomerAuditTrails(
		ctx context.Context,
		customerID uint,
		limit int,
		page int,
	) ([]models.AuditTrail, error)

	GetAuditTrailsByInvoiceID(
		ctx context.Context,
		invoiceID uint,
		customerID uint,
		limit int,
		page int,
	) ([]models.AuditTrail, error)
}
