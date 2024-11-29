package request_dto

import (
	"time"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type CreateInvoiceRequest struct {
	Sender            Sender                                  `json:"sender" binding:"required"`
	IssueDate         time.Time                               `json:"issue_date" binding:"required"`
	DueDate           time.Time                               `json:"due_date" binding:"required"`
	BillingCurrency   string                                  `json:"billing_currency" binding:"required"`
	Items             []InvoiceItem                           `json:"items" binding:"required"`
	Discount          float64                                 `json:"discount"`
	Notes             string                                  `json:"notes"`
	ReminderSchedules map[models.InvoiceReminderSchedule]bool `json:"reminder_schedules"`
	PaymentInfo       PaymentInfo                             `json:"payment_info"`
}

type Sender struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
	Email   string `json:"email" binding:"required"`
}

type InvoiceItem struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
}

type PaymentInfo struct {
	BankName      string `json:"bank_name" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required"`
	AccountName   string `json:"account_name" binding:"required"`
	AchRoutingNo  string `json:"ach_routing_no"`
	BankAddress   string `json:"bank_address"`
}
