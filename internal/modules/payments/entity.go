package payments

import (
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/transactions"
)

type CreatePaymentRequest struct {
	ID            int               `json:"id,omitempty"`
	LoanID        int               `json:"loan_id" validate:"required,gt=0"`
	Amount        float64           `json:"amount" validate:"required,gt=0"`
	Type          string            `json:"type" validate:"required,oneof=increase decrease"`
	PaymentDate   time.Time         `json:"payment_date" validate:"required"`
	TransactionID *int              `json:"-"`
	Transaction   *TransactionInput `json:"transaction,omitempty" validate:"omitempty"`
}

type UpdatePaymentRequest struct {
	Amount      *float64          `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Type        *string           `json:"type,omitempty" validate:"omitempty,oneof=increase decrease"`
	PaymentDate *time.Time        `json:"payment_date,omitempty" validate:"omitempty"`
	Transaction *TransactionInput `json:"transaction,omitempty" validate:"omitempty"`
}

type TransactionInput struct {
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
	LoanID        int                       `json:"loan_id"`
	Amount        float64                   `json:"amount"`
	Type          string                    `json:"type"`
	PaymentDate   time.Time                 `json:"payment_date"`
	TransactionID *int                      `json:"transaction_id"`
	Transaction   *transactions.Transaction `json:"transaction,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

type ListPaymentsRequest struct {
	Page     int `query:"page" validate:"gte=1"`
	PageSize int `query:"page_size" validate:"gte=1,lte=100"`
}
