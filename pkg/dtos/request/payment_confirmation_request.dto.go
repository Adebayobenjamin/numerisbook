package request_dto

import "time"

type PaymentConfirmationRequest struct {
	Amount      float64   `json:"amount" binding:"required"`
	PaymentDate time.Time `json:"payment_date" binding:"required"`
	IsPartial   bool      `json:"is_partial"`
}
