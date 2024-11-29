package response_dto

import (
	"time"

	"github.com/Adebayobenjamin/numerisbook/pkg/models"
)

type GetInvoiceDetailsResponse struct {
	ID                 uint                     `db:"id" json:"id"`
	InvoiceNumber      string                   `db:"invoice_number" json:"invoice_number"`
	Sender             *models.Sender           `db:"sender" json:"sender"`
	CustomerID         uint                     `db:"customer_id" json:"customer_id"`
	Customer           *models.Customer         `db:"customer" json:"customer"`
	IssueDate          time.Time                `db:"issue_date" json:"issue_date"`
	DueDate            time.Time                `db:"due_date" json:"due_date"`
	TotalAmountDue     float64                  `db:"total_amount_due" json:"total_amount_due"`
	Subtotal           float64                  `db:"subtotal" json:"subtotal"`
	IsFullyPaid        bool                     `db:"is_fully_paid" json:"is_fully_paid"`
	BillingCurrency    string                   `db:"billing_currency" json:"billing_currency"`
	Items              []models.InvoiceItem     `db:"items" json:"items"`
	Reminders          []models.InvoiceReminder `db:"reminders" json:"reminders"`
	Discount           float64                  `db:"discount" json:"discount"`
	Payments           []models.Payment         `db:"payments" json:"payments"`
	Status             models.InvoiceStatus     `db:"status" json:"status"`
	PaymentInformation *models.PaymentInfo      `db:"payment_information" json:"payment_information"`
	ShareableLink      *string                  `db:"shareable_link" json:"shareable_link"`
	Notes              string                   `db:"notes" json:"notes"`
	CreatedAt          time.Time                `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time                `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt          *time.Time               `db:"deleted_at" json:"deleted_at,omitempty"`
}
