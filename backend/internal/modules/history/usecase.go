package history

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type UseCase interface {
	CreateHistory(ctx context.Context, userID string, req *CreateHistoryRequest) (*History, error)
	UpdateHistory(ctx context.Context, userID string, historyID string, req *UpdateHistoryRequest) (*History, error)
	GetHistories(ctx context.Context, userID string, req *ListHistoryRequest) ([]*History, error)
}

type useCase struct {
	repo     Repository
	log      *logrus.Logger
	validate *validator.Validate
}

func NewUseCase(repo Repository, log *logrus.Logger, validate *validator.Validate) UseCase {
	return &useCase{repo: repo, log: log, validate: validate}
}

func (u useCase) CreateHistory(ctx context.Context, userID string, req *CreateHistoryRequest) (*History, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	isOwned, err := u.repo.IsBudgetOwnedByUser(ctx, req.BudgetID, userID)
	if err != nil {
		u.log.WithError(err).Error("Failed to check budget ownership")
		return nil, ErrInternalServer
	}

	if !isOwned {
		return nil, ErrForbidden
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	decAmount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, ErrInvalidAmount
	}

	if decAmount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrAmountMustPositive
	}

	history := &History{
		ID:        uuid.New().String(),
		BudgetID:  req.BudgetID,
		Amount:    decAmount,
		Date:      parsedDate,
		CreatedAt: time.Now(),
	}

	if err := u.repo.Create(ctx, history); err != nil {
		u.log.WithError(err).Error("Failed to create history")
		return nil, ErrInternalServer
	}

	return history, nil
}

func (u useCase) UpdateHistory(ctx context.Context, userID string, historyID string, req *UpdateHistoryRequest) (*History, error) {
	// 1. Validasi input
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Cek ownership melalui join table
	isOwned, err := u.repo.IsHistoryOwnedByUser(ctx, historyID, userID)
	if err != nil {
		u.log.WithError(err).Error("Failed to check history ownership")
		return nil, ErrInternalServer
	}

	if !isOwned {
		return nil, ErrForbidden
	}

	// 3. Get existing history
	existingHistory, err := u.repo.FindByID(ctx, historyID)
	if err != nil {
		u.log.WithError(err).Error("Failed to find history")
		return nil, ErrInternalServer
	}

	if existingHistory == nil {
		return nil, ErrHistoryNotFound
	}

	// 4. Parse dan validasi data baru
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	decAmount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, ErrInvalidAmount
	}

	if decAmount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrAmountMustPositive
	}

	// 5. Update data
	existingHistory.Amount = decAmount
	existingHistory.Date = parsedDate

	if err := u.repo.Update(ctx, existingHistory); err != nil {
		u.log.WithError(err).Error("Failed to update history")
		return nil, ErrInternalServer
	}

	return existingHistory, nil
}

func (u useCase) GetHistories(ctx context.Context, userID string, req *ListHistoryRequest) ([]*History, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	isOwned, err := u.repo.IsBudgetOwnedByUser(ctx, req.BudgetID, userID)
	if err != nil {
		u.log.WithError(err).Error("Failed to check budget ownership")
		return nil, ErrInternalServer
	}

	if !isOwned {
		return nil, ErrForbidden
	}

	histories, err := u.repo.FindManyByBudgetID(ctx, req.BudgetID)
	if err != nil {
		u.log.WithError(err).Error("Failed to fetch histories")
		return nil, ErrInternalServer
	}

	return histories, nil
}
