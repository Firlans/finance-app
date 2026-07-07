package payments

import (
	"context"
	"errors"
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/transactions"
	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	Save(ctx context.Context, payment *CreatePaymentRequest) (int, error)
	Update(ctx context.Context, id int, payment *UpdatePaymentRequest) error
	GetPaymentByID(ctx context.Context, id int) (*Payment, error)
	GetPaymentsByLoanID(ctx context.Context, loanID int) (*[]Payment, error)
	Delete(ctx context.Context, id int) error
}

type useCase struct {
	repo            Repository
	transactionRepo transactions.Repository
	validate        *validator.Validate
}

func NewUseCase(repo Repository, transactionRepo transactions.Repository, validate *validator.Validate) UseCase {
	return &useCase{
		repo:            repo,
		transactionRepo: transactionRepo,
		validate:        validate,
	}
}

func (uc *useCase) Save(ctx context.Context, payment *CreatePaymentRequest) (int, error) {
	if payment == nil {
		return 0, errors.New("payment is nil")
	}

	if err := uc.validate.Struct(payment); err != nil {
		return 0, err
	}

	if payment.Transaction != nil {
		txn := &transactions.Transaction{
			Amount:          payment.Transaction.Amount,
			TransactionType: payment.Transaction.TransactionType,
			Description:     payment.Transaction.Description,
			AccountID:       payment.Transaction.AccountID,
			CategoryID:      payment.Transaction.CategoryID,
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
		}

		if err := uc.transactionRepo.Save(ctx, txn); err != nil {
			return 0, err
		}

		payment.TransactionID = &txn.ID

	}

	if err := uc.repo.SavePayment(ctx, payment); err != nil {
		return 0, err
	}

	return payment.ID, nil
}

func (uc *useCase) Update(ctx context.Context, id int, payment *UpdatePaymentRequest) error {
	if payment == nil {
		return errors.New("payment is nil")
	}

	if err := uc.validate.Struct(payment); err != nil {
		return err
	}

	existingPayment, err := uc.repo.FindPaymentByID(ctx, id)
	if err != nil {
		return err
	}
	if existingPayment == nil {
		return errors.New("payment not found")
	}

	if payment.LoanID != nil {
		existingPayment.LoanID = *payment.LoanID
	}

	if payment.Transaction != nil {
		var txn *transactions.Transaction

		if existingPayment.TransactionID == nil {

			txn = &transactions.Transaction{
				Amount:          payment.Transaction.Amount,
				TransactionType: payment.Transaction.TransactionType,
				Description:     payment.Transaction.Description,
				AccountID:       payment.Transaction.AccountID,
				CategoryID:      payment.Transaction.CategoryID,
				CreatedAt:       time.Now().UTC(),
				UpdatedAt:       time.Now().UTC(),
			}

			if err := uc.transactionRepo.Save(ctx, txn); err != nil {
				return err
			}

			existingPayment.TransactionID = &txn.ID

		} else {
			txn, err = uc.transactionRepo.GetTransactionByID(ctx, *existingPayment.TransactionID)

			if err != nil {
				return err
			}
			if txn == nil {
				return errors.New("associated transaction not found")
			}

			txn.Amount = payment.Transaction.Amount
			txn.TransactionType = payment.Transaction.TransactionType
			txn.Description = payment.Transaction.Description
			txn.AccountID = payment.Transaction.AccountID
			txn.CategoryID = payment.Transaction.CategoryID
			txn.UpdatedAt = time.Now().UTC()

			if err := uc.transactionRepo.UpdateTransaction(ctx, txn); err != nil {
				return err
			}
		}
	}

	return uc.repo.UpdatePayment(ctx, existingPayment)
}

func (uc *useCase) GetPaymentByID(ctx context.Context, id int) (*Payment, error) {
	if id == 0 {
		return nil, nil
	}

	return uc.repo.FindPaymentByID(ctx, id)
}

func (uc *useCase) GetPaymentsByLoanID(ctx context.Context, loanID int) (*[]Payment, error) {
	if loanID == 0 {
		return nil, nil
	}

	return uc.repo.FindPaymentsByLoanID(ctx, loanID)
}

func (uc *useCase) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return nil
	}
	existingPayment, err := uc.repo.FindPaymentByID(ctx, id)
	if err != nil {
		return err
	}
	if existingPayment == nil {
		return errors.New("payment not found")
	}

	// 2. Jika payment ini punya relasi ke Transaction, hapus transaksinya dulu
	if existingPayment.TransactionID != nil {
		// Catatan: Pastikan nama method di transactionRepo kamu adalah DeleteTransaction (atau sesuaikan jika namanya Delete)
		err := uc.transactionRepo.DeleteTransaction(ctx, *existingPayment.TransactionID)
		if err != nil {
			return errors.New("failed to delete associated transaction: " + err.Error())
		}
	}
	return uc.repo.DeletePayment(ctx, id)
}
