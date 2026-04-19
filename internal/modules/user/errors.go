package user

import "errors"

// Domain Errors - Business Logic Level
var (
	ErrInternalServer            = errors.New("internal server error")
	ErrInvalidCredentials        = errors.New("invalid email or password")
	ErrUserNotFound              = errors.New("user not found")
	ErrInvalidCurrentPassword    = errors.New("invalid current password")
	ErrInvalidPasswordResetToken = errors.New("invalid or expired reset token")
)

// Repository Errors - Data Layer
var (
	ErrEmailTaken    = errors.New("email already taken")
	ErrUsernameTaken = errors.New("username already taken")
)
