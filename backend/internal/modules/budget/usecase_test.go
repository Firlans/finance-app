package budget_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/budget"
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

func (m *MockRepository) Create(ctx context.Context, b *budget.MonthlyBudget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, b *budget.MonthlyBudget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*budget.MonthlyBudget, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*budget.MonthlyBudget), args.Error(1)
}

func (m *MockRepository) FindMany(ctx context.Context, userID string, startDate, endDate time.Time) ([]*budget.MonthlyBudget, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*budget.MonthlyBudget), args.Error(1)
}

// ==========================================
// 2. HELPER SETUP
// ==========================================

func setupTest() (budget.UseCase, *MockRepository) {
	mockRepo := new(MockRepository)

	// Logger discard agar output test bersih
	log := logrus.New()
	log.SetOutput(io.Discard)

	validate := validator.New()

	useCase := budget.NewUseCase(mockRepo, log, validate)
	return useCase, mockRepo
}

// ==========================================
// 3. GROUP: CREATE BUDGET TESTS
// ==========================================

func TestCreateBudget_Success(t *testing.T) {
	u, mockRepo := setupTest()

	userID := "user-123"
	req := &budget.CreateBudgetRequest{
		Budget: "5000000.555",
		Date:   "2025-01-01",
	}

	// Expected: String "5000000.555" akan di-parse jadi decimal dengan 2 desimal
	// decimal.NewFromString("5000000.555") akan jadi 5000000.56 (rounded)

	// Expectation: Repository Create dipanggil dengan data yang benar
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(b *budget.MonthlyBudget) bool {
		// Parse expected value sama seperti di usecase
		expectedBudget, _ := decimal.NewFromString("5000000.555")

		return b.UserID == userID &&
			b.Budget.Equal(expectedBudget) && // Use Equal for decimal comparison
			b.Date.Format("2006-01-02") == "2025-01-01" &&
			b.ID != ""
	})).Return(nil)

	// Action
	res, err := u.CreateBudget(context.Background(), userID, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// Verify the budget amount
	expectedBudget, _ := decimal.NewFromString("5000000.555")
	assert.True(t, res.Budget.Equal(expectedBudget))
	assert.Equal(t, userID, res.UserID)

	mockRepo.AssertExpectations(t)
}

func TestCreateBudget_ValidationError(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"

	// Case: Budget 0 (Min 1)
	req := &budget.CreateBudgetRequest{
		Budget: "0", // String format
		Date:   "2025-01-01",
	}

	// Action
	res, err := u.CreateBudget(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateBudget_InvalidDateFormat(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"

	// Case: Format tanggal salah
	req := &budget.CreateBudgetRequest{
		Budget: "100000",     // String format
		Date:   "01-01-2025", // Seharusnya YYYY-MM-DD
	}

	// Action
	res, err := u.CreateBudget(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateBudget_RepositoryError(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"

	req := &budget.CreateBudgetRequest{
		Budget: "100000", // String format
		Date:   "2025-01-01",
	}

	// Expectation: DB Error
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db down"))

	// Action
	res, err := u.CreateBudget(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, budget.ErrInternalServer, err) // Harus dicovert ke Internal Server Error
	assert.Nil(t, res)
}

// ==========================================
// 4. GROUP: GET BUDGETS TESTS
// ==========================================

func TestGetBudgets_Success_DefaultDate(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"

	// Request kosong -> Default Bulan Ini
	req := &budget.ListBudgetRequest{}

	// Setup Dummy Data
	dummyBudgets := []*budget.MonthlyBudget{
		{ID: "1", Budget: decimal.NewFromInt(100000)},
	}

	// Expectation: FindMany dipanggil dengan Range Bulan Ini
	now := time.Now()
	expectedStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := expectedStart.AddDate(0, 1, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

	mockRepo.On("FindMany", mock.Anything, userID,
		mock.MatchedBy(func(t time.Time) bool { return t.Equal(expectedStart) }),
		mock.MatchedBy(func(t time.Time) bool { return t.Equal(expectedEnd) }),
	).Return(dummyBudgets, nil)

	// Action
	res, err := u.GetBudgets(context.Background(), userID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetBudgets_Success_WithFilter(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"

	// Request dengan Filter
	req := &budget.ListBudgetRequest{
		StartDate: "2024-01-01",
		EndDate:   "2024-12-31",
	}

	dummyBudgets := []*budget.MonthlyBudget{}

	// Expectation: Parse tanggal berhasil
	mockRepo.On("FindMany", mock.Anything, userID, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		start := args.Get(2).(time.Time)
		end := args.Get(3).(time.Time)

		assert.Equal(t, "2024-01-01", start.Format("2006-01-02"))
		// EndDate akan ditambah 23:59:59
		assert.Equal(t, "2024-12-31", end.Format("2006-01-02"))
	}).Return(dummyBudgets, nil)

	// Action
	res, err := u.GetBudgets(context.Background(), userID, req)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, res, 0)
}

func TestGetBudgets_RepositoryError(t *testing.T) {
	u, mockRepo := setupTest()
	userID := "user-123"
	req := &budget.ListBudgetRequest{}

	// Expectation: DB Error
	mockRepo.On("FindMany", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("query failed"))

	// Action
	res, err := u.GetBudgets(context.Background(), userID, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, budget.ErrInternalServer, err)
	assert.Nil(t, res)
}
