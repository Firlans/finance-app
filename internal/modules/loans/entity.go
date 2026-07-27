package loans

import "time"

type CreateLoanRequest struct {
	ID        int     `json:"id,omitempty"`
	UserID    string  `json:"user_id,omitempty" validate:"omitempty,uuid4"`
	Name      string  `json:"name" validate:"required,min=3,max=50"`
	Amount    float64 `json:"amount" validate:"gte=0"`
	LoanType  string  `json:"loan_type" validate:"required,oneof=debt receivable"`
	AccountID *int    `json:"account_id,omitempty" validate:"omitempty,gt=0"`
}

type CreateLoanResponse struct {
	ID int `json:"id"`
}

type UpdateLoanRequest struct {
	Name      *string  `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	LoanType  *string  `json:"loan_type,omitempty" validate:"omitempty,oneof=debt receivable"`
}

type Loan struct {
	ID                int       `json:"id"`
	UserID            string    `json:"user_id"`
	Name              string    `json:"name"`
	LoanType          string    `json:"loan_type"`
	OutstandingAmount float64   `json:"outstanding_amount"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ListLoansRequest struct {
	Page     int `query:"page" validate:"gte=1"`
	PageSize int `query:"page_size" validate:"gte=1,lte=100"`
}
