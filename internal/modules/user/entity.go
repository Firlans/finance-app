package user

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	DeletedAt *time.Time
}

// UserResponse: Format standar data user untuk output JSON
type UserResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string    `json:"username" example:"tubagus_aldi"`
	Email     string    `json:"email" example:"tubagus@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2025-01-01T15:00:00Z"`
}

// RegisterRequest: Validasi input saat daftar
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30" example:"tubagus_aldi"`
	Email    string `json:"email" validate:"required,email" example:"tubagus@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"RahasiaNegara123!"`
}

type RegisterResponse struct {
	UserResponse
}

// LoginRequest: Validasi input saat login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"tubagus@example.com"`
	Password string `json:"password" validate:"required" example:"RahasiaNegara123!"`
}

// LoginResponse: WAJIB mengandung Token
type LoginResponse struct {
	AccessToken string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string       `json:"token_type" example:"Bearer"`
	User        UserResponse `json:"user"`
}

type ResetPasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"RahasiaNegara123!"`
	NewPassword     string `json:"new_password" validate:"required,min=6" example:"RahasiaNegara456!"`
}

type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset successful"`
}

type ResetPasswordWithTokenRequest struct {
	Token           string `json:"token" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	NewPassword     string `json:"new_password" validate:"required,min=6" example:"RahasiaNegara456!"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword" example:"RahasiaNegara456!"`
}

type PasswordResetToken struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"tubagus@example.com"`
}

type ForgetPasswordResponse struct {
	Message string `json:"message" example:"Password reset link sent to your email"`
}
