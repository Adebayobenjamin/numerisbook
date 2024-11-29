package models

import "time"

type EventType string

const (
	EventTypeInvoiceCreated    EventType = "invoice_created"
	EventTypeInvoiceDuplicated EventType = "invoice_duplicated"
	EventTypePaymentConfirmed  EventType = "payment_confirmed"
)

type LogLevel string

const (
	LogLevelInfo    LogLevel = "info"
	LogLevelWarning LogLevel = "warning"
	LogLevelError   LogLevel = "error"
)

// AuditTrail struct represents an audit trail entity used to log actions performed on invoices by customers
type AuditTrail struct {
	ID         uint       `db:"id" json:"id"`
	EventType  EventType  `db:"event_type" json:"event_type"`
	LogLevel   LogLevel   `db:"log_level" json:"log_level"`
	Message    string     `db:"message" json:"message"`
	InvoiceID  uint       `db:"invoice_id" json:"invoice_id"`
	CustomerID uint       `db:"customer_id" json:"customer_id"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
}
