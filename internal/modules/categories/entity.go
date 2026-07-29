package categories

import "time"

type Category struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	TypeCategory string     `json:"type_category"`
	UserID       string     `json:"user_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type CreateCategoryRequest struct {
	Name         string  `json:"name" validate:"required"`
	Description  *string `json:"description,omitempty"`
	TypeCategory string  `json:"type_category,omitempty" validate:"omitempty,oneof=expense income both loan transfer fee"`
	UserID       string  `json:"user_id,omitempty" validate:"omitempty,uuid4"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	TypeCategory *string `json:"type_category,omitempty" validate:"omitempty,oneof=expense income both loan transfer fee"`
	UserID       *string `json:"user_id,omitempty" validate:"omitempty,uuid4"`
}

type CategoryResponse struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	TypeCategory string    `json:"type_category"`
	CreatedAt    time.Time `json:"created_at"`
}
