package models

import "time"

type InvoiceReminderSchedule string

const (
	InvoiceReminderSchedule14DaysBeforeDue InvoiceReminderSchedule = "14_days_before_due"
	InvoiceReminderSchedule7DaysBeforeDue  InvoiceReminderSchedule = "7_days_before_due"
	InvoiceReminderSchedule3DaysBeforeDue  InvoiceReminderSchedule = "3_days_before_due"
	InvoiceReminderSchedule1DayBeforeDue   InvoiceReminderSchedule = "1_day_before_due"
	InvoiceReminderScheduleOnDue           InvoiceReminderSchedule = "on_due"
)

type InvoiceReminder struct {
	ID           uint                    `db:"id" json:"id"`
	InvoiceID    uint                    `db:"invoice_id" json:"invoice_id"`
	CustomerID   uint                    `db:"customer_id" json:"customer_id"`
	Schedule     InvoiceReminderSchedule `db:"schedule" json:"schedule"`
	ReminderDate time.Time               `db:"reminder_date" json:"reminder_date"`
	CreatedAt    time.Time               `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time               `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time              `db:"deleted_at" json:"deleted_at"`
}
