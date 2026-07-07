package summary

import (
	"context"
	"time"
)

type UseCase interface {
	GetSummaryByModule(ctx context.Context, userID string, module Module, from time.Time, to time.Time) ([]TransactionCategoryBalance, error)
}

type useCase struct{ repo Repository }

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) GetSummaryByModule(ctx context.Context, userID string, module Module, from time.Time, to time.Time) ([]TransactionCategoryBalance, error) {
	if userID == "" {
		return []TransactionCategoryBalance{}, nil
	}
	return uc.repo.GetSummaryByModule(ctx, userID, module, from, to)
}
