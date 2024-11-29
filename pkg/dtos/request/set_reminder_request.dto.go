package request_dto

import "github.com/Adebayobenjamin/numerisbook/pkg/models"

type SetReminderRequest struct {
	Schedules map[models.InvoiceReminderSchedule]bool `json:"schedules" binding:"required"`
}
