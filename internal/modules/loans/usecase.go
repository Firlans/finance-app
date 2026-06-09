package loans

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/payments"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/transactions"
	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	Save(ctx context.Context, loan *CreateLoanRequest) (int, error)
	GetLoans(ctx context.Context, userID string) (*[]Loan, error)
	GetLoanByID(ctx context.Context, id int) (*Loan, error)
	Update(ctx context.Context, id int, loan *UpdateLoanRequest) error
	Delete(ctx context.Context, id int) error
}

type useCase struct {
	loanRepo        Repository
	paymentRepo     payments.Repository
	transactionRepo transactions.Repository
	validate        *validator.Validate
}

func NewUseCase(loanRepo Repository, paymentRepo payments.Repository, transactionRepo transactions.Repository, validate *validator.Validate) UseCase {
	return &useCase{
		loanRepo:        loanRepo,
		paymentRepo:     paymentRepo,
		transactionRepo: transactionRepo,
		validate:        validate,
	}
}

func (uc *useCase) Save(ctx context.Context, loan *CreateLoanRequest) (int, error) {
	if loan == nil {
		return 0, errors.New("loan is nil")
	}

	if err := uc.validate.Struct(loan); err != nil {
		return 0, err
	}

	if err := uc.loanRepo.SaveLoan(ctx, loan); err != nil {
		return 0, err
	}

	if loan.AccountID > 0 {
		transactionType := "credit"
		if loan.LoanType == "debt" {
			transactionType = "debit"
		}

		txn := &transactions.Transaction{
			Amount:          loan.Balance,
			TransactionType: transactionType,
			Description:     fmt.Sprintf("Loan payment for %s", loan.Name),
			AccountID:       loan.AccountID,
			CategoryID:      nil,
			UserID:          loan.UserID,
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
		}

		if err := uc.transactionRepo.Save(ctx, txn); err != nil {
			return 0, err
		}

		payment := &payments.CreatePaymentRequest{
			TransactionID: txn.ID,
			LoanID:        loan.ID,
		}

		if err := uc.paymentRepo.SavePayment(ctx, payment); err != nil {
			return 0, err
		}
	} else {
		placeholderPayment := &payments.CreatePaymentRequest{
			TransactionID: 0,
			LoanID:        loan.ID,
		}

		if err := uc.paymentRepo.SavePayment(ctx, placeholderPayment); err != nil {
			return 0, err
		}
	}

	return loan.ID, nil
}

func (uc *useCase) GetLoans(ctx context.Context, userID string) (*[]Loan, error) {
	if userID == "" {
		return nil, nil
	}

	return uc.loanRepo.FindLoansByUserID(ctx, userID)
}

func (uc *useCase) GetLoanByID(ctx context.Context, id int) (*Loan, error) {
	if id == 0 {
		return nil, nil
	}

	return uc.loanRepo.FindLoanByID(ctx, id)
}

func (uc *useCase) Update(ctx context.Context, id int, loan *UpdateLoanRequest) error {
	if loan == nil {
		return errors.New("loan is nil")
	}

	if err := uc.validate.Struct(loan); err != nil {
		return err
	}

	log.Println("masuk")
	if err := uc.loanRepo.UpdateLoan(ctx, id, loan); err != nil {
		return err
	}
	if loan.AccountID == nil {
		return nil
	}

	existingLoan, err := uc.loanRepo.FindLoanByID(ctx, id)
	if err != nil {
		return err
	}
	if existingLoan == nil {
		return errors.New("loan not found")
	}

	payment, err := uc.paymentRepo.FindFirstPaymentByLoanID(ctx, id)
	if err != nil {
		return err
	}

	if payment == nil {
		amount := existingLoan.Balance
		if loan.Balance != nil {
			amount = *loan.Balance
		}

		transactionType := "credit"
		loanType := existingLoan.LoanType
		if loan.LoanType != nil {
			loanType = *loan.LoanType
		}
		if loanType == "debt" {
			transactionType = "debit"
		}

		txn := &transactions.Transaction{
			Amount:          amount,
			TransactionType: transactionType,
			Description:     fmt.Sprintf("Loan payment for %s", existingLoan.Name),
			AccountID:       *loan.AccountID,
			CategoryID:      nil,
			UserID:          existingLoan.UserID,
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
		}

		if err := uc.transactionRepo.Save(ctx, txn); err != nil {
			return err
		}

		createPayment := &payments.CreatePaymentRequest{
			TransactionID: txn.ID,
			LoanID:        id,
		}

		return uc.paymentRepo.SavePayment(ctx, createPayment)
	}

	if payment.TransactionID == 0 {
		amount := existingLoan.Balance
		if loan.Balance != nil {
			amount = *loan.Balance
		}

		transactionType := "credit"
		loanType := existingLoan.LoanType
		if loan.LoanType != nil {
			loanType = *loan.LoanType
		}
		if loanType == "debt" {
			transactionType = "debit"
		}

		txn := &transactions.Transaction{
			Amount:          amount,
			TransactionType: transactionType,
			Description:     fmt.Sprintf("Loan payment for %s", existingLoan.Name),
			AccountID:       *loan.AccountID,
			CategoryID:      nil,
			UserID:          existingLoan.UserID,
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
		}

		if err := uc.transactionRepo.Save(ctx, txn); err != nil {
			return err
		}

		payment.TransactionID = txn.ID
		return uc.paymentRepo.UpdatePayment(ctx, payment)
	}

	txn, err := uc.transactionRepo.GetTransactionByID(ctx, payment.TransactionID)
	if err != nil {
		return err
	}
	if txn == nil {
		return errors.New("associated transaction not found")
	}

	loanType := existingLoan.LoanType
	if loan.LoanType != nil {
		loanType = *loan.LoanType
	}
	transactionType := "credit"
	if loanType == "debt" {
		transactionType = "debit"
	}

	if loan.Balance != nil {
		txn.Amount = *loan.Balance
	}

	if loan.Name != nil {
		txn.Description = fmt.Sprintf("Loan payment for %s", *loan.Name)
	}

	txn.TransactionType = transactionType
	txn.AccountID = *loan.AccountID
	txn.CategoryID = nil
	txn.UpdatedAt = time.Now().UTC()

	if err := uc.transactionRepo.UpdateTransaction(ctx, txn); err != nil {
		return err
	}

	return uc.paymentRepo.UpdatePayment(ctx, payment)
}

func (uc *useCase) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return nil
	}

	return uc.loanRepo.DeleteLoan(ctx, id)
}
