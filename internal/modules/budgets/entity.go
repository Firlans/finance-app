package budgets

import (
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Name         string    `json:"name" db:"name"`
	Amount       float64   `json:"amount" db:"amount"`
	IntervalType string    `json:"interval_type" db:"interval_type"`
	Day          *int      `json:"day" db:"day"`
	Date         *int      `json:"date" db:"date"`
	Month        *int      `json:"month" db:"month"`
	Repeat       bool      `json:"repeat" db:"repeat"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type BudgetHistory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	BudgetID  uuid.UUID `json:"budget_id" db:"budget_id"`
	Amount    float64   `json:"amount" db:"amount"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type BudgetCategory struct {
	BudgetID   uuid.UUID `json:"budget_id" db:"budget_id"`
	CategoryID int       `json:"category_id" db:"category_id"`
}

type CreateBudgetRequest struct {
	ID           *uuid.UUID `json:"id"`
	Name         string     `json:"name" validate:"required"`
	Amount       float64    `json:"amount" validate:"required,gt=0"`
	IntervalType string     `json:"interval_type" validate:"required"`
	Day          *int       `json:"day"`
	Date         *int       `json:"date"`
	Month        *int       `json:"month"`
	Repeat       bool       `json:"repeat"`
	CategoryIDs  []int      `json:"category_ids"`
}

type BudgetSummaryResponse struct {
	ID                 uuid.UUID   `json:"id"`
	Name               string      `json:"name"`
	Amount             float64     `json:"amount"`
	IntervalType       string      `json:"interval_type"`
	Day                *int        `json:"day"`
	Date               *int        `json:"date"`
	Month              *int        `json:"month"`
	Repeat             bool        `json:"repeat"`
	CurrentPeriodStart time.Time   `json:"current_period_start"`
	CurrentPeriodEnd   time.Time   `json:"current_period_end"`
	TotalSpent         float64     `json:"total_spent"`
	CategoryIDs        []int       `json:"category_ids"`
}
