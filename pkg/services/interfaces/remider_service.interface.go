package services_interfaces

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type RemiderService interface {
	SetInvoiceReminders(ctx context.Context, invoice *models.Invoice, customerID uint, shedules map[models.InvoiceReminderSchedule]bool) error
}
