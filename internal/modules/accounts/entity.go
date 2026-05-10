package accounts

import "time"

type ListAccountsRequest struct {
	Page     int `query:"page" validate:"gte=1"`
	PageSize int `query:"page_size" validate:"gte=1,lte=100"`
}

type CreateAccountRequest struct {
	ID          int     `json:"id,omitempty"`
	UserID      string  `json:"user_id,omitempty" validate:"omitempty,uuid4"`
	AccountName string  `json:"account_name" validate:"required,min=3,max=100"`
	Description string  `json:"description,omitempty" validate:"max=255"`
	Balance     float64 `json:"balance" validate:"gte=0"`
}

type CreateAccountResponse struct {
	ID int `json:"id"`
}

type UpdateAccountRequest struct {
	AccountName *string  `json:"account_name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=255"`
	Balance     *float64 `json:"balance,omitempty" validate:"omitempty,gte=0"`
}

type Account struct {
	ID          int       `json:"id"`
	UserID      string    `json:"user_id"`
	AccountName string    `json:"account_name"`
	Description string    `json:"description,omitempty"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
}
