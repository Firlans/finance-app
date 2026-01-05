package history

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrInternalServer = errors.New("internal server error")
	ErrForbidden      = errors.New("you don't have permission to access this budget")
	ErrBudgetNotFound = errors.New("budget not found")
)

type UseCase interface {
	CreateHistory(ctx context.Context, userID string, req *CreateHistoryRequest) (*History, error)
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
	//1. Validasi Input
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	//2.Pastikan Budget ID Milik User ini
	isOwned, err := u.repo.IsBudgetOwnedByUser(ctx, req.BudgetID, userID)
	if err != nil {
		u.log.WithError(err).Error("Failed to check budget ownership")
		return nil, ErrInternalServer
	}

	if !isOwned {
		return nil, ErrForbidden
	}

	//3. Parse Date
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	//4. Buat Object
	history := &History{
		ID:        uuid.New().String(),
		BudgetID:  req.BudgetID,
		Amount:    RoundFloat(req.Amount),
		Date:      parsedDate,
		CreatedAt: time.Now(),
	}

	//4. Simpan
	if err := u.repo.Create(ctx, history); err != nil {
		u.log.WithError(err).Error("Failed to create history")
		return nil, ErrInternalServer
	}

	return history, nil
}

func (u useCase) GetHistories(ctx context.Context, userID string, req *ListHistoryRequest) ([]*History, error) {
	// 1. Validasi Input Budget Wajib ada
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. SECURITY CHECK: User Hanya Boleh lihathistory Budget Miliknya
	isOwned, err := u.repo.IsBudgetOwnedByUser(ctx, req.BudgetID, userID)
	if err != nil {
		u.log.WithError(err).Error("Failed to check budget ownership")
		return nil, ErrInternalServer
	}

	if !isOwned {
		return nil, ErrForbidden
	}

	// 3. Ambil Data
	histories, err := u.repo.FindManyByBudgetID(ctx, req.BudgetID)
	if err != nil {
		u.log.WithError(err).Error("Failed to fetch histories")
		return nil, ErrInternalServer
	}

	return histories, nil
}
