package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/pkg/validators"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var (
	// Max amount: 1 billion (same as budget)
	MaxAmount = decimal.NewFromInt(1_000_000_000)
)

type UseCase interface {
	CreateTransaction(ctx context.Context, userID string, req *CreateTransactionRequest) (*Transaction, error)
	UpdateTransaction(ctx context.Context, userID string, transactionID string, req *UpdateTransactionRequest) (*Transaction, error)
	GetTransactions(ctx context.Context, userID string, req *ListTransactionRequest) ([]*Transaction, error)
}

type useCase struct {
	repo     Repository
	log      *logrus.Logger
	validate *validator.Validate
}

func NewUseCase(repo Repository, log *logrus.Logger, validate *validator.Validate) UseCase {
	return &useCase{repo: repo, log: log, validate: validate}
}

func (u useCase) CreateTransaction(ctx context.Context, userID string, req *CreateTransactionRequest) (*Transaction, error) {
	// 1. Struct validation
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Validate UUIDs
	if err := validators.ValidateUUID(userID); err != nil {
		u.log.WithError(err).Warn("Invalid user ID format")
		return nil, ErrInternalServer
	}

	if err := validators.ValidateUUID(req.AccountID); err != nil {
		return nil, validators.ErrInvalidUUID
	}

	// 3. Check budget ownership
	isOwned, err := u.repo.IsAccountOwnedByUser(ctx, req.AccountID, userID)
	if err != nil {
		u.log.WithError(err).Error("Failed to check budget ownership")
		return nil, ErrInternalServer
	}

	if !isOwned {
		u.log.WithFields(logrus.Fields{
			"user_id":    userID,
			"account_id": req.AccountID,
		}).Warn("Forbidden budget access attempt")
		return nil, ErrForbidden
	}

	// 4. Validate date
	parsedDate, err := validators.ValidateDate(req.Date)
	if err != nil {
		return nil, err
	}

	// 5. Validate amount
	decAmount, err := validators.ValidateDecimal(req.Amount, MaxAmount)
	if err != nil {
		if errors.Is(err, validators.ErrInvalidDecimal) {
			return nil, ErrInvalidAmount
		}
		if err.Error() == "amount must be greater than 0" {
			return nil, ErrAmountMustPositive
		}
		return nil, ErrInternalServer
	}

	// 6. Create entity
	transaction := &Transaction{
		ID:        uuid.New().String(),
		UserID:    userID,
		AccountID: req.AccountID,
		Amount:    decAmount,
		Date:      parsedDate,
		CreatedAt: time.Now(),
	}

	// 7. Save to DB
	if err := u.repo.Create(ctx, transaction); err != nil {
		u.log.WithError(err).Error("Failed to create transaction")
		return nil, ErrInternalServer
	}

	return transaction, nil
}

func (u useCase) UpdateTransaction(ctx context.Context, userID string, transactionID string, req *UpdateTransactionRequest) (*Transaction, error) {
	// 1. Validasi input
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Validate UUIDs
	if err := validators.ValidateUUID(userID); err != nil {
		u.log.WithError(err).Warn("Invalid user ID format")
		return nil, ErrInternalServer
	}

	if err := validators.ValidateUUID(transactionID); err != nil {
		return nil, validators.ErrInvalidUUID
	}

	// 3. Cek ownership melalui join table
	isOwned, err := u.repo.IsTransactionOwnedByUser(ctx, transactionID, userID)
	if err != nil {
		u.log.WithError(err).Error("Failed to check history ownership")
		return nil, ErrInternalServer
	}

	if !isOwned {
		u.log.WithFields(logrus.Fields{
			"user_id":        userID,
			"transaction_id": transactionID,
		}).Warn("Forbidden history access attempt")
		return nil, ErrForbidden
	}

	// 4. Get existing history
	existingTransaction, err := u.repo.FindByID(ctx, transactionID)
	if err != nil {
		u.log.WithError(err).Error("Failed to find history")
		return nil, ErrInternalServer
	}

	if existingTransaction == nil {
		return nil, ErrTransactionNotFound
	}

	// 5. Parse dan validasi data baru
	parsedDate, err := validators.ValidateDate(req.Date)
	if err != nil {
		return nil, err
	}

	decAmount, err := validators.ValidateDecimal(req.Amount, MaxAmount)
	if err != nil {
		if errors.Is(err, validators.ErrInvalidDecimal) {
			return nil, ErrInvalidAmount
		}
		if err.Error() == "amount must be greater than 0" {
			return nil, ErrAmountMustPositive
		}
		return nil, ErrInternalServer
	}

	// 6. Update data
	existingTransaction.Amount = decAmount
	existingTransaction.Date = parsedDate
	existingTransaction.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, existingTransaction); err != nil {
		u.log.WithError(err).Error("Failed to update transaction")
		return nil, ErrInternalServer
	}

	return existingTransaction, nil
}

func (u useCase) GetTransactions(ctx context.Context, userID string, req *ListTransactionRequest) ([]*Transaction, error) {
	// 1. Validate struct
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Validate UUIDs
	if err := validators.ValidateUUID(userID); err != nil {
		u.log.WithError(err).Warn("Invalid user ID format")
		return nil, ErrInternalServer
	}

	if err := validators.ValidateUUID(req.AccountID); err != nil {
		return nil, validators.ErrInvalidUUID
	}

	// 3. Check ownership
	isOwned, err := u.repo.IsAccountOwnedByUser(ctx, req.AccountID, userID)
	if err != nil {
		u.log.WithError(err).Error("Failed to check budget ownership")
		return nil, ErrInternalServer
	}

	if !isOwned {
		u.log.WithFields(logrus.Fields{
			"user_id":    userID,
			"account_id": req.AccountID,
		}).Warn("Forbidden budget access attempt")
		return nil, ErrForbidden
	}

	// 4. Fetch histories
	transactions, err := u.repo.FindManyByAccountID(ctx, req.AccountID)
	if err != nil {
		u.log.WithError(err).Error("Failed to fetch histories")
		return nil, ErrInternalServer
	}

	return transactions, nil
}
