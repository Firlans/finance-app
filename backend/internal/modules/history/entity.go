package history

import (
	"time"

	"github.com/shopspring/decimal"
)

type History struct {
	ID        string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	BudgetID  string          `json:"budget_id" example:"bf8a39e8-4226-4d04-a035-6453181878d6"`
	Amount    decimal.Decimal `json:"amount" example:"300000.00"`
	Date      time.Time       `json:"date" example:"2025-01-01T00:00:00Z"`
	CreatedAt time.Time       `json:"-"`
}

type CreateHistoryRequest struct {
	BudgetID string `json:"budget_id" validate:"required,uuid" example:"bf8a39e8-4226-4d04-a035-6453181878d6"`
	Amount   string `json:"amount" validate:"required" example:"50000"`
	Date     string `json:"date" validate:"required,datetime=2006-01-02" example:"2025-01-02"`
}

type ListHistoryRequest struct {
	BudgetID string `query:"budget_id" validate:"required,uuid"`
}
