package repositories_interfaces

import (
	"context"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type ReminderRepository interface {
	UpsertReminders(ctx context.Context, reminders []models.InvoiceReminder) error
}
