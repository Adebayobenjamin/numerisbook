package repositories

import (
	"context"
	"fmt"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repositories_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/repositories/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type reminderRepository struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

// UpsertReminders implements repositories_interfaces.ReminderRepository.
func (r *reminderRepository) UpsertReminders(ctx context.Context, reminders []models.InvoiceReminder) error {
	query := `
		INSERT INTO invoice_reminders (
			invoice_id, 
			customer_id, 
			schedule, 
			reminder_date,
			deleted_at
		) 
		VALUES (
			:invoice_id, 
			:customer_id, 
			:schedule, 
			:reminder_date,
			:deleted_at
		)
		ON DUPLICATE KEY UPDATE 
			customer_id = VALUES(customer_id),
			reminder_date = VALUES(reminder_date),
			deleted_at = VALUES(deleted_at),
			updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.NamedExecContext(ctx, query, reminders)
	if err != nil {
		return fmt.Errorf("failed to upsert reminders: %w", err)
	}

	return nil
}

func NewReminderRepository(
	db *sqlx.DB, logger *zerolog.Logger,
) repositories_interfaces.ReminderRepository {
	return &reminderRepository{
		db:     db,
		logger: logger,
	}
}
