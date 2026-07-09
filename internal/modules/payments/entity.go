package payments

import (
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/transactions"
)

type CreatePaymentRequest struct {
	ID            int               `json:"id,omitempty"`
	TransactionID *int              `json:"transaction_id,omitempty"`
	Transaction   *TransactionInput `json:"transaction,omitempty" validate:"omitempty"`
	LoanID        int               `json:"loan_id" validate:"required,gt=0"`
}

type UpdatePaymentRequest struct {
	LoanID      *int              `json:"loan_id,omitempty" validate:"omitempty,gt=0"`
	Transaction *TransactionInput `json:"transaction,omitempty" validate:"omitempty"`
}

type TransactionInput struct {
	Amount          float64 `json:"amount" validate:"required,gte=0"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=debit credit"`
	Description     string  `json:"description,omitempty" validate:"max=255"`
	AccountID       int     `json:"account_id" validate:"required,gt=0"`
	CategoryID      *int    `json:"category_id,omitempty"`
}

type CreatePaymentResponse struct {
	ID int `json:"id"`
}

type Payment struct {
	ID            int                       `json:"id"`
	TransactionID *int                      `json:"transaction_id"`
	Transaction   *transactions.Transaction `json:"transaction,omitempty"`
	LoanID        int                       `json:"loan_id"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

type ListPaymentsRequest struct {
	Page     int `query:"page" validate:"gte=1"`
	PageSize int `query:"page_size" validate:"gte=1,lte=100"`
}
