package budget

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal" // Import wajib
	"github.com/sirupsen/logrus"
)

/* Error Handling */
var (
	ErrInternalServer = errors.New("internal server error")
)

type UseCase interface {
	CreateBudget(ctx context.Context, userID string, req *CreateBudgetRequest) (*MonthlyBudget, error)
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
	// 1. Validasi Input
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Parse Date
	parsedDat, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// 3. Parse Budget (String -> Decimal) [PERBAIKAN DISINI]
	decBudget, err := decimal.NewFromString(req.Budget)
	if err != nil {
		// Return 400 Bad Request jika user mengirim "lima juta" bukannya "5000000"
		return nil, errors.New("invalid budget amount, must be numeric")
	}

	// Validasi Bisnis: Budget harus > 0
	if decBudget.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("budget must be greater than 0")
	}

	// 4. Buat Entity
	budget := &MonthlyBudget{
		ID:        uuid.New().String(),
		UserID:    userID,
		Budget:    decBudget,
		Date:      parsedDat,
		CreatedAt: time.Now(),
	}

	// 5. Simpan Ke DB
	if err := u.repo.Create(ctx, budget); err != nil {
		u.log.WithError(err).Error("Failed to create budget")
		return nil, ErrInternalServer
	}
	return budget, nil
}

func (u useCase) GetBudgets(ctx context.Context, userID string, req *ListBudgetRequest) ([]*MonthlyBudget, error) {
	// 1. Tentukan Default Range (Bulan Ini)
	now := time.Now()
	// Start: Tanggal 1 bulan ini, jam 00:00:00 UTC
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// End: Tanggal terakhir bulan ini, kita set ke jam 23:59:59.999999999
	// Agar mencakup seluruh detik di hari terakhir.
	endDate := startDate.AddDate(0, 1, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

	// 2. Override jika ada filter StartDate dari User
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			startDate = t
		}
	}

	// 3. Override jika ada filter EndDate dari User [PERBAIKAN UTAMA DISINI]
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			// Masalah sebelumnya: t adalah jam 00:00:00.
			// Solusi: Tambahkan waktu agar menjadi 23:59:59 di hari tersebut.
			endDate = t.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		}
	}

	// 4. Panggil Repository
	budgets, err := u.repo.FindMany(ctx, userID, startDate, endDate)
	if err != nil {
		u.log.WithError(err).Error("Failed to fetch budgets")
		return nil, ErrInternalServer
	}

	return budgets, nil
}
