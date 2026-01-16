package budget

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type UseCase interface {
	CreateBudget(ctx context.Context, userID string, req *CreateBudgetRequest) (*MonthlyBudget, error)
	UpdateBudget(ctx context.Context, userID string, budgetID string, req *UpdateBudgetRequest) (*MonthlyBudget, error)
	GetBudgets(ctx context.Context, userID string, req *ListBudgetRequest) ([]*MonthlyBudget, error)
}

type useCase struct {
	repo     Repository
	log      *logrus.Logger
	validate *validator.Validate
}

func NewUseCase(repo Repository, log *logrus.Logger, validate *validator.Validate) UseCase {
	return &useCase{repo: repo, log: log, validate: validate}
}

func (u useCase) CreateBudget(ctx context.Context, userID string, req *CreateBudgetRequest) (*MonthlyBudget, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	decBudget, err := decimal.NewFromString(req.Budget)
	if err != nil {
		return nil, ErrInvalidBudgetAmount
	}

	if decBudget.LessThanOrEqual(decimal.Zero) {
		return nil, ErrBudgetMustPositive
	}

	budget := &MonthlyBudget{
		ID:        uuid.New().String(),
		UserID:    userID,
		Budget:    decBudget,
		Date:      parsedDate,
		CreatedAt: time.Now(),
	}

	if err := u.repo.Create(ctx, budget); err != nil {
		u.log.WithError(err).Error("Failed to create budget")
		return nil, ErrInternalServer
	}
	return budget, nil
}

func (u useCase) UpdateBudget(ctx context.Context, userID string, budgetID string, req *UpdateBudgetRequest) (*MonthlyBudget, error) {
	// 1. Validasi input
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Cek apakah budget exists dan milik user ini
	existingBudget, err := u.repo.FindByID(ctx, budgetID)
	if err != nil {
		u.log.WithError(err).Error("Failed to find budget")
		return nil, ErrInternalServer
	}

	if existingBudget == nil {
		return nil, ErrBudgetNotFound
	}

	if existingBudget.UserID != userID {
		return nil, ErrForbidden
	}

	// 3. Parse dan validasi data baru
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	decBudget, err := decimal.NewFromString(req.Budget)
	if err != nil {
		return nil, ErrInvalidBudgetAmount
	}

	if decBudget.LessThanOrEqual(decimal.Zero) {
		return nil, ErrBudgetMustPositive
	}

	// 4. Update data
	existingBudget.Budget = decBudget
	existingBudget.Date = parsedDate

	if err := u.repo.Update(ctx, existingBudget); err != nil {
		u.log.WithError(err).Error("Failed to update budget")
		return nil, ErrInternalServer
	}

	return existingBudget, nil
}

func (u useCase) GetBudgets(ctx context.Context, userID string, req *ListBudgetRequest) ([]*MonthlyBudget, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			startDate = t
		}
	}

	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			endDate = t.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		}
	}

	budgets, err := u.repo.FindMany(ctx, userID, startDate, endDate)
	if err != nil {
		u.log.WithError(err).Error("Failed to fetch budgets")
		return nil, ErrInternalServer
	}

	return budgets, nil
}
