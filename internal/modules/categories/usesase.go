package categories

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	Save(ctx context.Context, category *Category) error
	GetCategories(ctx context.Context, userID *string, typeCategory string) ([]Category, error)
	GetCategoryByID(ctx context.Context, id int) (*Category, error)
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type useCase struct {
	repo     Repository
	validate *validator.Validate
}

func NewUseCase(repo Repository, validate *validator.Validate) UseCase {
	return &useCase{repo: repo, validate: validate}
}

func (u *useCase) Save(ctx context.Context, category *Category) error {
	if err := u.validate.Struct(category); err != nil {
		return err
	}
	return u.repo.Save(ctx, category)
}
func (u *useCase) GetCategories(ctx context.Context, userID *string, typeCategory string) ([]Category, error) {
	return u.repo.GetCategories(ctx, userID, typeCategory)
}
func (u *useCase) GetCategoryByID(ctx context.Context, id int) (*Category, error) {
	return u.repo.GetCategoryByID(ctx, id)
}
func (u *useCase) UpdateCategory(ctx context.Context, category *Category) error {
	if err := u.validate.Struct(category); err != nil {
		return err
	}
	existingCategory, err := u.repo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		return err
	}

	if category.Name != "" {
		existingCategory.Name = category.Name
	}
	existingCategory.Description = category.Description
	if category.TypeCategory != "" {
		existingCategory.TypeCategory = category.TypeCategory
	}

	return u.repo.UpdateCategory(ctx, existingCategory)
}
func (u *useCase) DeleteCategory(ctx context.Context, id int) error {
	return u.repo.DeleteCategory(ctx, id)
}
