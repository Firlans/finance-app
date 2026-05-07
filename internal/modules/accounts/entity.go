package accounts

type ListAccountsRequest struct {
	Page     int `query:"page" validate:"gte=1"`
	PageSize int `query:"page_size" validate:"gte=1,lte=100"`
}

type CreateAccountRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description string  `json:"description,omitempty" validate:"max=255"`
	Balance     float64 `json:"balance" validate:"required,gte=0"`
}

type UpdateAccountRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=255"`
	Balance     *float64 `json:"balance,omitempty" validate:"omitempty,gte=0"`
}

type Account struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	AccountName string  `json:"account_name"`
	Description string  `json:"description,omitempty"`
	Balance     float64 `json:"balance"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
