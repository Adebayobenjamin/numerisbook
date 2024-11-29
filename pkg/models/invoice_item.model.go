package models

import "time"

type InvoiceItem struct {
	ID          uint       `db:"id" json:"id"`
	InvoiceID   uint       `db:"invoice_id" json:"invoice_id"`
	Description string     `db:"description" json:"description"`
	Quantity    int        `db:"quantity" json:"quantity"`
	UnitPrice   float64    `db:"unit_price" json:"unit_price"`
	TotalPrice  float64    `db:"total_price" json:"total_price"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}
