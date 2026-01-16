package history

import "errors"

// Domain Errors - Business Logic Level
var (
	ErrInternalServer  = errors.New("internal server error")
	ErrForbidden       = errors.New("you don't have permission to access this resource")
	ErrBudgetNotFound  = errors.New("budget not found")
	ErrHistoryNotFound = errors.New("history not found")
)

// Validation Errors - Input Level
var (
	ErrInvalidAmount      = errors.New("amount must be a valid number")
	ErrInvalidDateFormat  = errors.New("invalid date format, use YYYY-MM-DD")
	ErrAmountMustPositive = errors.New("amount must be greater than 0")
)
