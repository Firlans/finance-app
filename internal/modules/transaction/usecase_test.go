package transaction_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/transaction"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==========================================
// 1. MOCK OBJECTS
// ==========================================

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, h *transaction.Transaction) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, h *transaction.Transaction) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*transaction.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transaction.Transaction), args.Error(1)
}

func (m *MockRepository) FindManyByAccountID(ctx context.Context, accountID string) ([]*transaction.Transaction, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*transaction.Transaction), args.Error(1)
}

func (m *MockRepository) IsAccountOwnedByUser(ctx context.Context, accountID string, userID string) (bool, error) {
	args := m.Called(ctx, accountID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) IsTransactionOwnedByUser(ctx context.Context, transactionID string, userID string) (bool, error) {
	args := m.Called(ctx, transactionID, userID)
	return args.Bool(0), args.Error(1)
}

// ==========================================
// 2. HELPER SETUP
// ==========================================

func setupTest() (transaction.UseCase, *MockRepository) {
	mockRepo := new(MockRepository)

	// Logger discard agar output test bersih
	log := logrus.New()
	log.SetOutput(io.Discard)

	validate := validator.New()

	useCase := transaction.NewUseCase(mockRepo, log, validate)
	return useCase, mockRepo
}

// ==========================================
// 3. GROUP: CREATE HISTORY TESTS
// ==========================================

func TestCreateHistory_Success(t *testing.T) {
	u, mockRepo := setupTest()

	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &transaction.CreateTransactionRequest{
		AccountID: accountID,
		Amount:    "50000",
		Date:      "2025-01-02",
	}

	// 1. Expect Ownership Check: TRUE
	mockRepo.On("IsAccountOwnedByUser", mock.Anything, accountID, userID).Return(true, nil)

	// 2. Expect Create: Success
	expectedAmount := decimal.NewFromInt(50000)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(h *transaction.Transaction) bool {
		return h.AccountID == accountID &&
			h.Amount.Equal(expectedAmount) &&
			h.Date.Format("2006-01-02") == "2025-01-02" &&
			h.ID != ""
	})).Return(nil)

	// Action
	res, err := u.CreateTransaction(context.Background(), userID, req)

	// Assertions
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.True(t, res.Amount.Equal(expectedAmount))
	}

	mockRepo.AssertExpectations(t)
}

func TestCreateHistory_Forbidden(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	req := &transaction.CreateTransactionRequest{
		AccountID: accountID,
		Amount:    "50000",
		Date:      "2025-01-02",
	}

	// 1. Expect Ownership Check: FALSE
	mockRepo.On("IsAccountOwnedByUser", mock.Anything, accountID, userID).Return(false, nil)

	// Action
	res, err := u.CreateTransaction(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrForbidden, err)
	assert.Nil(t, res)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateHistory_AmountZero(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &transaction.CreateTransactionRequest{
		AccountID: accountID,
		Amount:    "0",
		Date:      "2025-01-02",
	}

	// Mock diperlukan karena validasi amount = 0 terjadi SETELAH ownership check
	mockRepo.On("IsAccountOwnedByUser", mock.Anything, accountID, userID).Return(true, nil)

	// Action
	res, err := u.CreateTransaction(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrAmountMustPositive, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateHistory_InvalidBudgetID(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"

	// Invalid UUID format
	req := &transaction.CreateTransactionRequest{
		AccountID: "invalid-uuid",
		Amount:    "100",
		Date:      "2025-01-02",
	}

	// Action
	res, err := u.CreateTransaction(context.Background(), userID, req)

	// Assertions - Validator akan reject ini di struct validation
	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "IsAccountOwnedByUser")
}

func TestCreateHistory_OwnershipCheckError(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &transaction.CreateTransactionRequest{
		AccountID: accountID,
		Amount:    "100",
		Date:      "2025-01-01",
	}

	// DB Error saat cek ownership
	mockRepo.On("IsAccountOwnedByUser", mock.Anything, accountID, userID).Return(false, errors.New("db error"))

	// Action
	res, err := u.CreateTransaction(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrInternalServer, err)
	assert.Nil(t, res)
}

func TestCreateHistory_InvalidDateFormat(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &transaction.CreateTransactionRequest{
		AccountID: accountID,
		Amount:    "100",
		Date:      "01-01-2025", // Invalid format - will be caught by validator tag
	}

	// Action
	res, err := u.CreateTransaction(context.Background(), userID, req)

	// Assertions
	// Validator will catch this BEFORE any usecase logic
	assert.Error(t, err)
	assert.Nil(t, res)
	// Should not reach any repository calls
	mockRepo.AssertNotCalled(t, "IsAccountOwnedByUser")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateHistory_InvalidAmountFormat(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &transaction.CreateTransactionRequest{
		AccountID: accountID,
		Amount:    "not-a-number",
		Date:      "2025-01-01",
	}

	mockRepo.On("IsAccountOwnedByUser", mock.Anything, accountID, userID).Return(true, nil)

	// Action
	res, err := u.CreateTransaction(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrInvalidAmount, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateHistory_NegativeAmount(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &transaction.CreateTransactionRequest{
		AccountID: accountID,
		Amount:    "-100",
		Date:      "2025-01-01",
	}

	mockRepo.On("IsAccountOwnedByUser", mock.Anything, accountID, userID).Return(true, nil)

	// Action
	res, err := u.CreateTransaction(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrAmountMustPositive, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "Create")
}

// ==========================================
// 4. GROUP: GET HISTORIES TESTS
// ==========================================

func TestGetHistories_Success(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &transaction.ListTransactionRequest{AccountID: accountID}

	dummyHistories := []*transaction.Transaction{
		{ID: "h1", Amount: decimal.NewFromInt(5000)},
	}

	// 1. Cek Milik User
	mockRepo.On("IsAccountOwnedByUser", mock.Anything, accountID, userID).Return(true, nil)

	// 2. Ambil Data
	mockRepo.On("FindManyByAccountID", mock.Anything, accountID).Return(dummyHistories, nil)

	// Action
	res, err := u.GetTransactions(context.Background(), userID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestGetHistories_Forbidden(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	req := &transaction.ListTransactionRequest{AccountID: accountID}

	// 1. Cek Milik User -> False
	mockRepo.On("IsAccountOwnedByUser", mock.Anything, accountID, userID).Return(false, nil)

	// Action
	res, err := u.GetTransactions(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrForbidden, err)
	assert.Nil(t, res)

	mockRepo.AssertNotCalled(t, "FindManyByAccountID")
}

func TestGetHistories_RepositoryError(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "550e8400-e29b-41d4-a716-446655440000"
	accountID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &transaction.ListTransactionRequest{AccountID: accountID}

	mockRepo.On("IsAccountOwnedByUser", mock.Anything, accountID, userID).Return(true, nil)
	mockRepo.On("FindManyByAccountID", mock.Anything, accountID).Return(nil, errors.New("db fail"))

	res, err := u.GetTransactions(context.Background(), userID, req)

	assert.Error(t, err)
	assert.Equal(t, transaction.ErrInternalServer, err)
	assert.Nil(t, res)
}
