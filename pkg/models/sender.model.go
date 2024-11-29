package models

import "time"

type Sender struct {
	ID        uint       `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Phone     string     `db:"phone" json:"phone"`
	Address   string     `db:"address" json:"address"`
	Email     string     `db:"email" json:"email"`
	InvoiceID uint       `db:"invoice_id" json:"invoice_id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}
