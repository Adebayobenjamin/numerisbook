package services

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/helper"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repositories_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/repositories/interfaces"
	services_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/services/interfaces"
)

type auditService struct {
	auditRepository repositories_interfaces.AuditTrailRepository
}

// CreateAuditTrail implements services_interfaces.AuditService.
func (a *auditService) CreateAuditTrail(ctx context.Context, eventType models.EventType, logLevel models.LogLevel, message string, invoiceID uint, customerID uint) error {
	// Log event to audit trail
	return a.auditRepository.
		LogEvent(ctx, eventType, logLevel, message, invoiceID, customerID)
}

// GetAuditTrailsByInvoiceID implements services_interfaces.AuditService.
func (a *auditService) GetAuditTrailsByInvoiceID(ctx context.Context, invoiceID uint, customerID uint, limit int, page int) ([]models.AuditTrail, error) {
	offset := helper.GetOffset(page, limit)
	return a.auditRepository.
		GetByInvoiceIDAndCustomerID(ctx, invoiceID, customerID, limit, offset)
}

// GetCustomerAuditTrails implements services_interfaces.AuditService.
func (a *auditService) GetCustomerAuditTrails(ctx context.Context, customerID uint, limit int, page int) ([]models.AuditTrail, error) {
	offset := helper.GetOffset(page, limit)
	return a.auditRepository.
		GetAllCustomerAuditTrails(ctx, customerID, limit, offset)
}

func NewAuditService(
	auditRepository repositories_interfaces.AuditTrailRepository,
) services_interfaces.AuditService {
	return &auditService{
		auditRepository: auditRepository,
	}
}
