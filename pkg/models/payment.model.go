package models

import "time"

type Payment struct {
	ID        uint       `db:"id" json:"id"`
	InvoiceID uint       `db:"invoice_id" json:"invoice_id"`
	Amount    float64    `db:"amount" json:"amount"`
	IsPartial bool       `db:"is_partial" json:"is_partial"`
	Date      time.Time  `db:"date" json:"date"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}
