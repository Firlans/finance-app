package budget

import (
	"context"
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/pkg/validators"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var (
	// Max budget: 1 billion
	MaxBudgetAmount = decimal.NewFromInt(1_000_000_000)
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
	// 1. Struct validation
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Validate UUID
	if err := validators.ValidateUUID(userID); err != nil {
		u.log.WithError(err).Warn("Invalid user ID format")
		return nil, ErrInternalServer
	}

	// 3. Validate date with range check
	parsedDate, err := validators.ValidateDate(req.Date)
	if err != nil {
		return nil, err
	}

	// 4. Validate decimal with max amount check
	decBudget, err := validators.ValidateDecimal(req.Budget, MaxBudgetAmount)
	if err != nil {
		return nil, err
	}

	// 5. Create entity
	budget := &MonthlyBudget{
		ID:        uuid.New().String(),
		UserID:    userID,
		Budget:    decBudget,
		Date:      parsedDate,
		CreatedAt: time.Now(),
	}

	// 6. Save to DB
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

	// 2. Validate UUIDs
	if err := validators.ValidateUUID(userID); err != nil {
		u.log.WithError(err).Warn("Invalid user ID format")
		return nil, ErrInternalServer
	}

	if err := validators.ValidateUUID(budgetID); err != nil {
		return nil, validators.ErrInvalidUUID
	}

	// 3. Cek apakah budget exists dan milik user ini
	existingBudget, err := u.repo.FindByID(ctx, budgetID)
	if err != nil {
		u.log.WithError(err).Error("Failed to find budget")
		return nil, ErrInternalServer
	}

	if existingBudget == nil {
		return nil, ErrBudgetNotFound
	}

	if existingBudget.UserID != userID {
		u.log.WithFields(logrus.Fields{
			"user_id":   userID,
			"budget_id": budgetID,
			"owner_id":  existingBudget.UserID,
		}).Warn("Forbidden budget access attempt")
		return nil, ErrForbidden
	}

	// 4. Parse dan validasi data baru
	parsedDate, err := validators.ValidateDate(req.Date)
	if err != nil {
		return nil, err
	}

	decBudget, err := validators.ValidateDecimal(req.Budget, MaxBudgetAmount)
	if err != nil {
		return nil, err
	}

	// 5. Update data
	existingBudget.Budget = decBudget
	existingBudget.Date = parsedDate

	if err := u.repo.Update(ctx, existingBudget); err != nil {
		u.log.WithError(err).Error("Failed to update budget")
		return nil, ErrInternalServer
	}

	return existingBudget, nil
}

func (u useCase) GetBudgets(ctx context.Context, userID string, req *ListBudgetRequest) ([]*MonthlyBudget, error) {
	// Validate user ID
	if err := validators.ValidateUUID(userID); err != nil {
		u.log.WithError(err).Warn("Invalid user ID format")
		return nil, ErrInternalServer
	}

	// Default to current month
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

	// Parse start date if provided
	if req.StartDate != "" {
		if t, err := validators.ValidateDate(req.StartDate); err == nil {
			startDate = t
		} else {
			return nil, err
		}
	}

	// Parse end date if provided
	if req.EndDate != "" {
		if t, err := validators.ValidateDate(req.EndDate); err == nil {
			endDate = t.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		} else {
			return nil, err
		}
	}

	budgets, err := u.repo.FindMany(ctx, userID, startDate, endDate)
	if err != nil {
		u.log.WithError(err).Error("Failed to fetch budgets")
		return nil, ErrInternalServer
	}

	return budgets, nil
}
