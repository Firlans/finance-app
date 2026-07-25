package goals

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Name          string    `json:"name" db:"name"`
	TargetAmount  float64   `json:"target_amount" db:"target_amount"`
	CurrentAmount float64   `json:"current_amount" db:"current_amount"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateGoalRequest struct {
	Name          string  `json:"name" validate:"required"`
	TargetAmount  float64 `json:"target_amount" validate:"required,gt=0"`
	CurrentAmount float64 `json:"current_amount" validate:"gte=0"`
}

type UpdateGoalRequest struct {
	Name          string  `json:"name" validate:"required"`
	TargetAmount  float64 `json:"target_amount" validate:"required,gt=0"`
	CurrentAmount float64 `json:"current_amount" validate:"gte=0"`
}

type GoalResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	IsCompleted   bool      `json:"is_completed"`
}
