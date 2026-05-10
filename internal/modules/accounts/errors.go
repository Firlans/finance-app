package accounts

import "errors"

var (
	ErrAccountNotFound  = errors.New("account not found")
	ErrInvalidAccountID = errors.New("invalid account ID")
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrInternalServer   = errors.New("internal server error")
)

var (
	ErrAccountNameTaken = errors.New("account name already taken")
)
