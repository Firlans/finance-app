package budgets

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UseCase interface {
	UpsertBudget(ctx context.Context, userID uuid.UUID, req *CreateBudgetRequest) error
	GetBudgetSummaries(ctx context.Context, userID uuid.UUID, dateStr string) ([]BudgetSummaryResponse, error)
}

type useCase struct {
	repo     Repository
	validate *validator.Validate
}

func NewUseCase(repo Repository, validate *validator.Validate) UseCase {
	return &useCase{
		repo:     repo,
		validate: validate,
	}
}

func (u *useCase) UpsertBudget(ctx context.Context, userID uuid.UUID, req *CreateBudgetRequest) error {
	if err := u.validate.Struct(req); err != nil {
		return err
	}

	budget := &Budget{
		UserID:       userID,
		Name:         req.Name,
		Amount:       req.Amount,
		IntervalType: req.IntervalType,
		Day:          req.Day,
		Date:         req.Date,
		Month:        req.Month,
		Repeat:       req.Repeat,
	}
	
	if req.ID != nil {
		budget.ID = *req.ID
	}

	return u.repo.UpsertBudget(ctx, budget, req.CategoryIDs)
}

func (u *useCase) GetBudgetSummaries(ctx context.Context, userID uuid.UUID, dateStr string) ([]BudgetSummaryResponse, error) {
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return u.repo.GetBudgetSummaries(ctx, userID, dateStr)
}
