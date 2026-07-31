package payments_test

import (
	"context"
	"testing"
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/payments"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/transactions"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) SavePayment(ctx context.Context, payment *payments.CreatePaymentRequest) error {
	args := m.Called(ctx, payment)
	payment.ID = 1
	return args.Error(0)
}

func (m *MockPaymentRepo) FindPaymentByID(ctx context.Context, id int) (*payments.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*payments.Payment), args.Error(1)
}

func (m *MockPaymentRepo) FindPaymentsByLoanID(ctx context.Context, loanID int) (*[]payments.Payment, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]payments.Payment), args.Error(1)
}

func (m *MockPaymentRepo) FindPaymentsByTransactionID(ctx context.Context, transactionID int) (*[]payments.Payment, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]payments.Payment), args.Error(1)
}

func (m *MockPaymentRepo) FindFirstPaymentByLoanID(ctx context.Context, loanID int) (*payments.Payment, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*payments.Payment), args.Error(1)
}

func (m *MockPaymentRepo) UpdatePayment(ctx context.Context, payment *payments.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepo) DeletePayment(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentRepo) GetLoanTypeByID(ctx context.Context, loanID int) (string, error) {
	args := m.Called(ctx, loanID)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentRepo) GetHutangCategoryID(ctx context.Context) (*int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) Save(ctx context.Context, txn *transactions.Transaction) error {
	args := m.Called(ctx, txn)
	txn.ID = 100
	return args.Error(0)
}

func (m *MockTransactionRepo) GetTransactions(ctx context.Context, userID string, from string, to string, page int) ([]transactions.Transaction, error) {
	return nil, nil
}

func (m *MockTransactionRepo) GetTransactionByID(ctx context.Context, id int) (*transactions.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactions.Transaction), args.Error(1)
}

func (m *MockTransactionRepo) UpdateTransaction(ctx context.Context, txn *transactions.Transaction) error {
	args := m.Called(ctx, txn)
	return args.Error(0)
}

func (m *MockTransactionRepo) DeleteTransaction(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTransactionRepo) IsTransactionLinkedToPayment(ctx context.Context, transactionID int) (bool, error) {
	return false, nil
}

func (m *MockTransactionRepo) CreateTransfer(ctx context.Context, userID string, req *transactions.CreateTransferRequest) error {
	return nil
}

func (m *MockTransactionRepo) GetTransferByID(ctx context.Context, userID string, id int) (*transactions.CreateTransferRequest, error) {
	return nil, nil
}

func (m *MockTransactionRepo) UpdateTransfer(ctx context.Context, userID string, id int, req *transactions.CreateTransferRequest) error {
	return nil
}

func (m *MockTransactionRepo) GetCategoryIDByName(ctx context.Context, name string) (*int, error) {
	return nil, nil
}

func TestDetermineTransactionType(t *testing.T) {
	ctx := context.Background()
	val := validator.New()

	t.Run("Receivable - Decrease (Pembayaran) -> Credit", func(t *testing.T) {
		mockPaymentRepo := new(MockPaymentRepo)
		mockTxnRepo := new(MockTransactionRepo)
		uc := payments.NewUseCase(mockPaymentRepo, mockTxnRepo, val)

		mockPaymentRepo.On("GetLoanTypeByID", ctx, 1).Return("receivable", nil)
		mockPaymentRepo.On("GetHutangCategoryID", ctx).Return(nil, nil)
		mockTxnRepo.On("Save", ctx, mock.MatchedBy(func(txn *transactions.Transaction) bool {
			return txn.TransactionType == "credit"
		})).Return(nil)
		mockPaymentRepo.On("SavePayment", ctx, mock.Anything).Return(nil)

		req := &payments.CreatePaymentRequest{
			LoanID:      1,
			Amount:      100000,
			Type:        "decrease",
			PaymentDate: time.Now(),
			Transaction: &payments.TransactionInput{
				AccountID: 1,
			},
		}

		_, err := uc.Save(ctx, req)
		assert.NoError(t, err)
		mockTxnRepo.AssertExpectations(t)
	})

	t.Run("Receivable - Increase (Penambahan) -> Debit", func(t *testing.T) {
		mockPaymentRepo := new(MockPaymentRepo)
		mockTxnRepo := new(MockTransactionRepo)
		uc := payments.NewUseCase(mockPaymentRepo, mockTxnRepo, val)

		mockPaymentRepo.On("GetLoanTypeByID", ctx, 1).Return("receivable", nil)
		mockPaymentRepo.On("GetHutangCategoryID", ctx).Return(nil, nil)
		mockTxnRepo.On("Save", ctx, mock.MatchedBy(func(txn *transactions.Transaction) bool {
			return txn.TransactionType == "debit"
		})).Return(nil)
		mockPaymentRepo.On("SavePayment", ctx, mock.Anything).Return(nil)

		req := &payments.CreatePaymentRequest{
			LoanID:      1,
			Amount:      100000,
			Type:        "increase",
			PaymentDate: time.Now(),
			Transaction: &payments.TransactionInput{
				AccountID: 1,
			},
		}

		_, err := uc.Save(ctx, req)
		assert.NoError(t, err)
		mockTxnRepo.AssertExpectations(t)
	})

	t.Run("Debt - Decrease (Pembayaran) -> Debit", func(t *testing.T) {
		mockPaymentRepo := new(MockPaymentRepo)
		mockTxnRepo := new(MockTransactionRepo)
		uc := payments.NewUseCase(mockPaymentRepo, mockTxnRepo, val)

		mockPaymentRepo.On("GetLoanTypeByID", ctx, 2).Return("debt", nil)
		mockPaymentRepo.On("GetHutangCategoryID", ctx).Return(nil, nil)
		mockTxnRepo.On("Save", ctx, mock.MatchedBy(func(txn *transactions.Transaction) bool {
			return txn.TransactionType == "debit"
		})).Return(nil)
		mockPaymentRepo.On("SavePayment", ctx, mock.Anything).Return(nil)

		req := &payments.CreatePaymentRequest{
			LoanID:      2,
			Amount:      100000,
			Type:        "decrease",
			PaymentDate: time.Now(),
			Transaction: &payments.TransactionInput{
				AccountID: 1,
			},
		}

		_, err := uc.Save(ctx, req)
		assert.NoError(t, err)
		mockTxnRepo.AssertExpectations(t)
	})

	t.Run("Debt - Increase (Penambahan) -> Credit", func(t *testing.T) {
		mockPaymentRepo := new(MockPaymentRepo)
		mockTxnRepo := new(MockTransactionRepo)
		uc := payments.NewUseCase(mockPaymentRepo, mockTxnRepo, val)

		mockPaymentRepo.On("GetLoanTypeByID", ctx, 2).Return("debt", nil)
		mockPaymentRepo.On("GetHutangCategoryID", ctx).Return(nil, nil)
		mockTxnRepo.On("Save", ctx, mock.MatchedBy(func(txn *transactions.Transaction) bool {
			return txn.TransactionType == "credit"
		})).Return(nil)
		mockPaymentRepo.On("SavePayment", ctx, mock.Anything).Return(nil)

		req := &payments.CreatePaymentRequest{
			LoanID:      2,
			Amount:      100000,
			Type:        "increase",
			PaymentDate: time.Now(),
			Transaction: &payments.TransactionInput{
				AccountID: 1,
			},
		}

		_, err := uc.Save(ctx, req)
		assert.NoError(t, err)
		mockTxnRepo.AssertExpectations(t)
	})
}
