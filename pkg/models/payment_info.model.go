package models

import "time"

type PaymentInfo struct {
	ID            uint       `db:"id" json:"id"`
	InvoiceID     uint       `db:"invoice_id" json:"invoice_id"`
	BankName      string     `db:"bank_name" json:"bank_name"`
	AccountNumber string     `db:"account_number" json:"account_number"`
	AccountName   string     `db:"account_name" json:"account_name"`
	AchRoutingNo  string     `db:"ach_routing_no" json:"ach_routing_no"`
	BankAddress   string     `db:"bank_address" json:"bank_address"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
}
