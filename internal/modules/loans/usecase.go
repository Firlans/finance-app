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

	now := time.Now().UTC()

	payment := &payments.CreatePaymentRequest{
		LoanID:      loan.ID,
		Amount:      loan.Amount,
		Type:        "increase",
		PaymentDate: now,
	}

	var transactionType string
	if loan.AccountID != nil && *loan.AccountID > 0 {
		transactionType = "debit"
		if loan.LoanType == "debt" {
			transactionType = "credit"
		}
		
		payment.Transaction = &payments.TransactionInput{
			Description:     fmt.Sprintf("Loan creation: %s", loan.Name),
			AccountID:       *loan.AccountID,
		}
	}

	// use paymentRepo.SavePayment (which now might not create the transaction, so we need to either create it here or let payment usecase do it)
	// Actually, wait, paymentRepo.SavePayment just saves the payment in DB. It doesn't create the transaction!
	// So we DO need to create the transaction here if we are just calling paymentRepo.SavePayment!

	if payment.Transaction != nil {
		txn := &transactions.Transaction{
			Amount:          payment.Amount,
			TransactionType: transactionType,
			Description:     payment.Transaction.Description,
			AccountID:       payment.Transaction.AccountID,
			CategoryID:      payment.Transaction.CategoryID,
			UserID:          loan.UserID,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		if err := uc.transactionRepo.Save(ctx, txn); err != nil {
			return 0, err
		}

		payment.TransactionID = &txn.ID
	}

	if err := uc.paymentRepo.SavePayment(ctx, payment); err != nil {
		return 0, err
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
	return nil
}

func (uc *useCase) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return nil
	}
	paymentsList, err := uc.paymentRepo.FindPaymentsByLoanID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Lakukan looping untuk menghapus setiap transaksi dan pembayaran
	if paymentsList != nil {
		for _, p := range *paymentsList {
			// Jika payment ini punya transaksi, hapus transaksinya terlebih dahulu
			if p.TransactionID != nil {
				// Pastikan method hapus di transactionRepo sesuai dengan interface-mu
				err := uc.transactionRepo.DeleteTransaction(ctx, *p.TransactionID)
				if err != nil {
					return fmt.Errorf("failed to delete transaction for payment %d: %w", p.ID, err)
				}
			}

			// Setelah transaksi terhapus, hapus payment-nya
			err = uc.paymentRepo.DeletePayment(ctx, p.ID)
			if err != nil {
				return fmt.Errorf("failed to delete payment %d: %w", p.ID, err)
			}
		}
	}

	return uc.loanRepo.DeleteLoan(ctx, id)
}
