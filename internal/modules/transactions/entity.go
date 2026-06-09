package transactions

import "time"

type Transaction struct {
	ID              int       `json:"id"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	Description     string    `json:"description"`
	AccountID       int       `json:"account_id"`
	CategoryID      *int      `json:"category_id,omitempty"`
	UserID          string    `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateTransactionRequest struct {
	Amount          float64 `json:"amount" validate:"required,gte=0"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=debit credit"`
	Description     string  `json:"description,omitempty" validate:"max=255"`
	AccountID       int     `json:"account_id" validate:"required"`
	CategoryID      int     `json:"category_id" validate:"required"`
}

type UpdateTransactionRequest struct {
	Amount          *float64 `json:"amount,omitempty" validate:"omitempty,gte=0"`
	TransactionType *string  `json:"transaction_type,omitempty" validate:"omitempty,oneof=debit credit"`
	Description     *string  `json:"description,omitempty" validate:"omitempty,max=255"`
	AccountID       *int     `json:"account_id" validate:"required"`
	CategoryID      *int     `json:"category_id" validate:"required"`
}
