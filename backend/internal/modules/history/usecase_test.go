package history_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/history"
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

func (m *MockRepository) Create(ctx context.Context, h *history.History) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, h *history.History) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*history.History, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*history.History), args.Error(1)
}

func (m *MockRepository) FindManyByBudgetID(ctx context.Context, budgetID string) ([]*history.History, error) {
	args := m.Called(ctx, budgetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*history.History), args.Error(1)
}

func (m *MockRepository) IsBudgetOwnedByUser(ctx context.Context, budgetID string, userID string) (bool, error) {
	args := m.Called(ctx, budgetID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) IsHistoryOwnedByUser(ctx context.Context, historyID string, userID string) (bool, error) {
	args := m.Called(ctx, historyID, userID)
	return args.Bool(0), args.Error(1)
}

// ==========================================
// 2. HELPER SETUP
// ==========================================

func setupTest() (history.UseCase, *MockRepository) {
	mockRepo := new(MockRepository)

	// Logger discard agar output test bersih
	log := logrus.New()
	log.SetOutput(io.Discard)

	validate := validator.New()

	useCase := history.NewUseCase(mockRepo, log, validate)
	return useCase, mockRepo
}

// ==========================================
// 3. GROUP: CREATE HISTORY TESTS
// ==========================================

func TestCreateHistory_Success(t *testing.T) {
	u, mockRepo := setupTest()

	userID := "user-123"
	budgetID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &history.CreateHistoryRequest{
		BudgetID: budgetID,
		Amount:   "50000",
		Date:     "2025-01-02",
	}

	// 1. Expect Ownership Check: TRUE
	mockRepo.On("IsBudgetOwnedByUser", mock.Anything, budgetID, userID).Return(true, nil)

	// 2. Expect Create: Success
	expectedAmount := decimal.NewFromInt(50000)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(h *history.History) bool {
		return h.BudgetID == budgetID &&
			h.Amount.Equal(expectedAmount) &&
			h.Date.Format("2006-01-02") == "2025-01-02" &&
			h.ID != ""
	})).Return(nil)

	// Action
	res, err := u.CreateHistory(context.Background(), userID, req)

	// Assertions
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.True(t, res.Amount.Equal(expectedAmount))
	}

	mockRepo.AssertExpectations(t)
}

func TestCreateHistory_Forbidden(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	budgetID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	req := &history.CreateHistoryRequest{
		BudgetID: budgetID,
		Amount:   "50000",
		Date:     "2025-01-02",
	}

	// 1. Expect Ownership Check: FALSE
	mockRepo.On("IsBudgetOwnedByUser", mock.Anything, budgetID, userID).Return(false, nil)

	// Action
	res, err := u.CreateHistory(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, history.ErrForbidden, err)
	assert.Nil(t, res)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateHistory_AmountZero(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	budgetID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &history.CreateHistoryRequest{
		BudgetID: budgetID,
		Amount:   "0",
		Date:     "2025-01-02",
	}

	// Mock diperlukan karena validasi amount = 0 terjadi SETELAH ownership check
	mockRepo.On("IsBudgetOwnedByUser", mock.Anything, budgetID, userID).Return(true, nil)

	// Action
	res, err := u.CreateHistory(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, history.ErrAmountMustPositive, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateHistory_InvalidBudgetID(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"

	// Invalid UUID format
	req := &history.CreateHistoryRequest{
		BudgetID: "invalid-uuid",
		Amount:   "100",
		Date:     "2025-01-02",
	}

	// Action
	res, err := u.CreateHistory(context.Background(), userID, req)

	// Assertions - Validator akan reject ini di struct validation
	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "IsBudgetOwnedByUser")
}

func TestCreateHistory_OwnershipCheckError(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	budgetID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &history.CreateHistoryRequest{
		BudgetID: budgetID,
		Amount:   "100",
		Date:     "2025-01-01",
	}

	// DB Error saat cek ownership
	mockRepo.On("IsBudgetOwnedByUser", mock.Anything, budgetID, userID).Return(false, errors.New("db error"))

	// Action
	res, err := u.CreateHistory(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, history.ErrInternalServer, err)
	assert.Nil(t, res)
}

func TestCreateHistory_InvalidDateFormat(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	budgetID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &history.CreateHistoryRequest{
		BudgetID: budgetID,
		Amount:   "100",
		Date:     "01-01-2025", // Invalid format - will be caught by validator tag
	}

	// Action
	res, err := u.CreateHistory(context.Background(), userID, req)

	// Assertions
	// Validator will catch this BEFORE any usecase logic
	assert.Error(t, err)
	assert.Nil(t, res)
	// Should not reach any repository calls
	mockRepo.AssertNotCalled(t, "IsBudgetOwnedByUser")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateHistory_InvalidAmountFormat(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	budgetID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &history.CreateHistoryRequest{
		BudgetID: budgetID,
		Amount:   "not-a-number",
		Date:     "2025-01-01",
	}

	mockRepo.On("IsBudgetOwnedByUser", mock.Anything, budgetID, userID).Return(true, nil)

	// Action
	res, err := u.CreateHistory(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, history.ErrInvalidAmount, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateHistory_NegativeAmount(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	budgetID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &history.CreateHistoryRequest{
		BudgetID: budgetID,
		Amount:   "-100",
		Date:     "2025-01-01",
	}

	mockRepo.On("IsBudgetOwnedByUser", mock.Anything, budgetID, userID).Return(true, nil)

	// Action
	res, err := u.CreateHistory(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, history.ErrAmountMustPositive, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "Create")
}

// ==========================================
// 4. GROUP: GET HISTORIES TESTS
// ==========================================

func TestGetHistories_Success(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	budgetID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &history.ListHistoryRequest{BudgetID: budgetID}

	dummyHistories := []*history.History{
		{ID: "h1", Amount: decimal.NewFromInt(5000)},
	}

	// 1. Cek Milik User
	mockRepo.On("IsBudgetOwnedByUser", mock.Anything, budgetID, userID).Return(true, nil)

	// 2. Ambil Data
	mockRepo.On("FindManyByBudgetID", mock.Anything, budgetID).Return(dummyHistories, nil)

	// Action
	res, err := u.GetHistories(context.Background(), userID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestGetHistories_Forbidden(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	budgetID := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

	req := &history.ListHistoryRequest{BudgetID: budgetID}

	// 1. Cek Milik User -> False
	mockRepo.On("IsBudgetOwnedByUser", mock.Anything, budgetID, userID).Return(false, nil)

	// Action
	res, err := u.GetHistories(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, history.ErrForbidden, err)
	assert.Nil(t, res)

	mockRepo.AssertNotCalled(t, "FindManyByBudgetID")
}

func TestGetHistories_RepositoryError(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	budgetID := "d3b07384-d9a3-432d-a410-6c6734105211"

	req := &history.ListHistoryRequest{BudgetID: budgetID}

	mockRepo.On("IsBudgetOwnedByUser", mock.Anything, budgetID, userID).Return(true, nil)
	mockRepo.On("FindManyByBudgetID", mock.Anything, budgetID).Return(nil, errors.New("db fail"))

	res, err := u.GetHistories(context.Background(), userID, req)

	assert.Error(t, err)
	assert.Equal(t, history.ErrInternalServer, err)
	assert.Nil(t, res)
}
