package models

import "time"

type InvoiceStatus string

const (
	InvoiceStatusDraft          InvoiceStatus = "draft"
	InvoiceStatusSent           InvoiceStatus = "sent"
	InvoiceStatusPaid           InvoiceStatus = "paid"
	InvoiceStatusPendingPayment InvoiceStatus = "pending payment"
)

// Invoice represents an invoice entity
type Invoice struct {
	ID              uint              `db:"id" json:"id,omitempty"`
	InvoiceNumber   string            `db:"invoice_number" json:"invoice_number,omitempty"`
	Sender          *Sender           `db:"sender" json:"sender,omitempty"`
	CustomerID      uint              `db:"customer_id" json:"customer_id,omitempty"`
	Customer        *Customer         `db:"customer" json:"customer,omitempty"`
	IssueDate       time.Time         `db:"issue_date" json:"issue_date,omitempty"`
	DueDate         time.Time         `db:"due_date" json:"due_date,omitempty"`
	TotalAmountDue  float64           `db:"total_amount_due" json:"total_amount_due,omitempty"`
	Subtotal        float64           `db:"subtotal" json:"subtotal,omitempty"`
	IsFullyPaid     bool              `db:"is_fully_paid" json:"is_fully_paid,omitempty"`
	BillingCurrency string            `db:"billing_currency" json:"billing_currency,omitempty"`
	Items           []InvoiceItem     `db:"items" json:"items,omitempty"`
	Discount        float64           `db:"discount" json:"discount,omitempty"`
	Payments        []Payment         `db:"payments" json:"payments,omitempty"`
	Reminders       []InvoiceReminder `db:"reminders" json:"reminders,omitempty"`
	Status          InvoiceStatus     `db:"status" json:"status,omitempty"`
	PaymentInfo     *PaymentInfo      `db:"payment_info" json:"payment_info,omitempty"`
	ShareableLink   *string           `db:"shareable_link" json:"shareable_link,omitempty"`
	Notes           string            `db:"notes" json:"notes,omitempty"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt       *time.Time        `db:"deleted_at" json:"deleted_at,omitempty"`
}
