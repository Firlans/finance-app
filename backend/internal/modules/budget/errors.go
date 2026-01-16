package budget

import "errors"

// Domain Errors - Business Logic Level
var (
	ErrInternalServer = errors.New("internal server error")
	ErrBudgetNotFound = errors.New("budget not found")
	ErrForbidden      = errors.New("you don't have permission to modify this budget")
)

// Validation Errors - Input Level
var (
	ErrInvalidBudgetAmount = errors.New("invalid budget amount, must be numeric")
	ErrInvalidDateFormat   = errors.New("invalid date format, use YYYY-MM-DD")
	ErrBudgetMustPositive  = errors.New("budget must be greater than 0")
)
