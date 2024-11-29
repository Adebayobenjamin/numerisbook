package services

import (
	"context"
	"time"

	"github.com/Adebayobenjamin/numerisbook/pkg/helper"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repositories_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/repositories/interfaces"
	services_interfaces "github.com/Adebayobenjamin/numerisbook/pkg/services/interfaces"
)

type reminderService struct {
	reminderRepository repositories_interfaces.ReminderRepository
}

// SetInvoiceReminders implements services_interfaces.RemiderService.
func (r *reminderService) SetInvoiceReminders(ctx context.Context, invoice *models.Invoice, customerID uint, shedules map[models.InvoiceReminderSchedule]bool) error {
	var reminders []models.InvoiceReminder

	for schedule := range shedules {
		remider := models.InvoiceReminder{
			InvoiceID:    invoice.ID,
			CustomerID:   customerID,
			ReminderDate: getReminderDateFromSchedule(schedule, invoice.DueDate),
			Schedule:     schedule,
		}

		if !shedules[schedule] {
			remider.DeletedAt = helper.ReturnPointer(time.Now())
		}

		reminders = append(reminders, remider)
	}

	return r.reminderRepository.UpsertReminders(ctx, reminders)
}

func getReminderDateFromSchedule(schedule models.InvoiceReminderSchedule, invoiceDueDate time.Time) time.Time {
	// NOTE: We are using the invoice due date to calculate the reminder date
	// we do this by adding the negative value of the schedule to the invoice due date
	switch schedule {
	case models.InvoiceReminderSchedule14DaysBeforeDue:
		return invoiceDueDate.AddDate(0, 0, -14)
	case models.InvoiceReminderSchedule7DaysBeforeDue:
		return invoiceDueDate.AddDate(0, 0, -7)
	case models.InvoiceReminderSchedule3DaysBeforeDue:
		return invoiceDueDate.AddDate(0, 0, -3)
	case models.InvoiceReminderSchedule1DayBeforeDue:
		return invoiceDueDate.AddDate(0, 0, -1)
	default:
		return invoiceDueDate
	}
}

func NewReminderService(
	reminderRepository repositories_interfaces.ReminderRepository,
) services_interfaces.RemiderService {
	return &reminderService{
		reminderRepository: reminderRepository,
	}
}
