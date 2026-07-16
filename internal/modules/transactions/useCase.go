package transactions

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	Save(ctx context.Context, transaction *Transaction) error
	GetTransactions(ctx context.Context, userID string, from string, to string, page int) ([]Transaction, error)
	GetTransactionByID(ctx context.Context, id int) (*Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
	DeleteTransaction(ctx context.Context, id int) error
}

type useCase struct {
	repo     Repository
	validate *validator.Validate
}

func NewUseCase(repo Repository, validate *validator.Validate) UseCase {
	return &useCase{repo: repo, validate: validate}
}

func (uc *useCase) Save(ctx context.Context, transaction *Transaction) error {
	if transaction == nil {
		return nil
	}

	// Audit fields
	transaction.CreatedAt = time.Now().UTC()
	transaction.UpdatedAt = transaction.CreatedAt

	if err := uc.validate.Struct(transaction); err != nil {
		return err
	}

	return uc.repo.Save(ctx, transaction)
}

func (uc *useCase) GetTransactions(ctx context.Context, userID string, from string, to string, page int) ([]Transaction, error) {
	if userID == "" {
		return nil, nil
	}
	return uc.repo.GetTransactions(ctx, userID, from, to, page)
}

func (uc *useCase) GetTransactionByID(ctx context.Context, id int) (*Transaction, error) {
	if id == 0 {
		return nil, nil
	}

	return uc.repo.GetTransactionByID(ctx, id)
}

func (uc *useCase) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	if transaction == nil {
		return nil
	}

	req := *transaction

	if err := uc.validate.Struct(&req); err != nil {
		return err
	}

	existing, err := uc.repo.GetTransactionByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if existing == nil {
		return nil
	}

	existing.Amount = req.Amount
	existing.TransactionType = req.TransactionType
	existing.Description = req.Description
	existing.AccountID = req.AccountID
	existing.CategoryID = req.CategoryID

	// Only update transaction_date if client provides it.
	if !req.TransactionDate.IsZero() {
		existing.TransactionDate = req.TransactionDate
	}

	existing.UpdatedAt = time.Now().UTC()

	return uc.repo.UpdateTransaction(ctx, existing)

}

func (uc *useCase) DeleteTransaction(ctx context.Context, id int) error {
	if id == 0 {
		return nil
	}
	account, err := uc.repo.GetTransactionByID(ctx, id)
	if err != nil {
		return err
	}

	if account == nil {
		return nil
	}

	return uc.repo.DeleteTransaction(ctx, id)
}
