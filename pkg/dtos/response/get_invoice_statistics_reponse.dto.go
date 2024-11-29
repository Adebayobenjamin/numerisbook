package response_dto

type GetInvoiceStatisticsResponse struct {
	TotalPaid          int     `db:"total_paid" json:"total_paid"`
	TotalPaidAmount    float64 `db:"total_paid_amount" json:"total_paid_amount"`
	TotalOverDue       int     `db:"total_over_due" json:"total_over_due"`
	TotalOverDueAmount float64 `db:"total_over_due_amount" json:"total_over_due_amount"`
	TotalDraft         int     `db:"total_draft" json:"total_draft"`
	TotalDraftAmount   float64 `db:"total_draft_amount" json:"total_draft_amount"`
	TotalUnpaid        int     `db:"total_unpaid" json:"total_unpaid"`
	TotalUnpaidAmount  float64 `db:"total_unpaid_amount" json:"total_unpaid_amount"`
}
