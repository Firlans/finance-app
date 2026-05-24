package user_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/user"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// ==========================================
// 1. MOCK OBJECTS
// ==========================================

// MockRepository memalsukan behavior Repository
type MockRepository struct {
	mock.Mock
}

type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendPasswordResetEmail(ctx context.Context, to, resetURL string) error {
	args := m.Called(ctx, to, resetURL)
	return args.Error(0)
}

func (m *MockRepository) Save(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	// Handle jika return nil (User not found)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) SavePasswordResetToken(ctx context.Context, token *user.PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRepository) FindPasswordResetToken(ctx context.Context, token string) (*user.PasswordResetToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.PasswordResetToken), args.Error(1)
}

func (m *MockRepository) DeletePasswordResetToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

// ==========================================
// 2. HELPER SETUP
// ==========================================

// setupTest mengembalikan Interface UseCase, MockRepo, dan Config untuk dimanipulasi
func setupTest() (user.UseCase, *MockRepository, *MockMailer, *viper.Viper) {
	mockRepo := new(MockRepository)
	mockMailer := new(MockMailer)

	// Logger buang ke tong sampah (supaya terminal bersih)
	log := logrus.New()
	log.SetOutput(io.Discard)

	validate := validator.New()

	// Config default yang valid
	cfg := viper.New()
	cfg.Set("jwt.secret", "secret_key_testing_123")
	cfg.Set("jwt.ttl", "1h")

	useCase := user.NewUseCase(mockRepo, log, validate, cfg, mockMailer)

	return useCase, mockRepo, mockMailer, cfg
}

// ==========================================
// 3. GROUP: REGISTER TESTS
// ==========================================

func TestRegister_Success(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "validuser",
		Email:    "valid@example.com",
		Password: "password123",
	}

	// Expectation: Repo.Save dipanggil sekali dengan data yang cocok
	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(userObj *user.User) bool {
		return userObj.Email == req.Email && userObj.Username == req.Username && userObj.ID != ""
	})).Return(nil)

	// Action
	resp, err := u.Register(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.Email)
	assert.NotEmpty(t, resp.ID)

	mockRepo.AssertExpectations(t)
}

func TestRegister_ValidationError(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "",                  // Invalid: Kosong
		Email:    "invalid-email-fmt", // Invalid: Format salah
		Password: "123",               // Invalid: Terlalu pendek
	}

	// Expectation: Repo.Save TIDAK BOLEH dipanggil
	// (Karena harusnya gagal di layer validasi validator/struct)

	// Action
	resp, err := u.Register(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertNotCalled(t, "Save") // Pastikan tidak tembus ke DB
}

func TestRegister_DuplicateEmail(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "newuser",
		Email:    "taken@example.com",
		Password: "password123",
	}

	// Expectation: Repo return error EmailTaken
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(user.ErrEmailTaken)

	resp, err := u.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, user.ErrEmailTaken, err)
	assert.Nil(t, resp)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "takenuser",
		Email:    "new@example.com",
		Password: "password123",
	}

	// Expectation: Repo return error UsernameTaken
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(user.ErrUsernameTaken)

	resp, err := u.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, user.ErrUsernameTaken, err)
	assert.Nil(t, resp)
}

func TestRegister_RepositoryError(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.RegisterRequest{
		Username: "user",
		Email:    "email@example.com",
		Password: "password123",
	}

	// Expectation: Repo gagal koneksi DB
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db connection lost"))

	resp, err := u.Register(context.Background(), req)

	// Assertions: Harus return ErrInternalServer (jangan bocorkan error asli ke user)
	assert.Error(t, err)
	assert.Equal(t, user.ErrInternalServer, err)
	assert.Nil(t, resp)
}

// ==========================================
// 4. GROUP: LOGIN TESTS
// ==========================================

func TestLogin_Success(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	// Siapkan password hash yang valid
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	dummyUser := &user.User{
		ID:        "uuid-123",
		Username:  "loginuser",
		Email:     "login@example.com",
		Password:  string(hashedPwd),
		CreatedAt: time.Now(),
	}

	req := &user.LoginRequest{
		Email:    "login@example.com",
		Password: "password123", // Password raw yang cocok
	}

	// Expectation: FindByEmail sukses
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(dummyUser, nil)

	// Action
	resp, err := u.Login(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken) // Token harus ada
	assert.Equal(t, "Bearer", resp.TokenType)
	assert.Equal(t, dummyUser.Email, resp.User.Email)
}

func TestLogin_UserNotFound(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.LoginRequest{
		Email:    "ghost@example.com",
		Password: "password123",
	}

	// Expectation: FindByEmail return nil user & nil error (atau error not found, tergantung implementasi repo)
	// Sesuai implementasi mock kita: return nil, nil
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)

	resp, err := u.Login(context.Background(), req)

	// Harus InvalidCredentials, JANGAN "User Not Found" (Security)
	assert.Error(t, err)
	assert.Equal(t, user.ErrInvalidCredentials, err)
	assert.Nil(t, resp)
}

func TestLogin_WrongPassword(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("realpassword"), bcrypt.DefaultCost)
	dummyUser := &user.User{
		Email:    "user@example.com",
		Password: string(hashedPwd),
	}

	req := &user.LoginRequest{
		Email:    "user@example.com",
		Password: "WRONG_PASSWORD",
	}

	// Expectation: User ketemu, tapi nanti bcrypt compare gagal
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(dummyUser, nil)

	resp, err := u.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, user.ErrInvalidCredentials, err)
	assert.Nil(t, resp)
}

func TestLogin_RepositoryError(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.LoginRequest{
		Email:    "error@example.com",
		Password: "any",
	}

	// Expectation: DB mati
	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, errors.New("db timeout"))

	resp, err := u.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, user.ErrInternalServer, err)
	assert.Nil(t, resp)
}

func TestRefreshToken_GeneratesUniqueJWTID(t *testing.T) {
	u, mockRepo, _, cfg := setupTest()

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	dummyUser := &user.User{
		ID:        "uuid-123",
		Username:  "loginuser",
		Email:     "login@example.com",
		Password:  string(hashedPwd),
		CreatedAt: time.Now(),
	}

	loginReq := &user.LoginRequest{
		Email:    dummyUser.Email,
		Password: "password123",
	}

	mockRepo.On("FindByEmail", mock.Anything, loginReq.Email).Return(dummyUser, nil).Once()
	mockRepo.On("FindByID", mock.Anything, dummyUser.ID).Return(dummyUser, nil).Once()

	loginResp, err := u.Login(context.Background(), loginReq)
	assert.NoError(t, err)
	assert.NotNil(t, loginResp)

	refreshResp, err := u.RefreshToken(context.Background(), dummyUser.ID, loginResp.AccessToken, time.Now().Add(time.Hour))
	assert.NoError(t, err)
	assert.NotNil(t, refreshResp)
	assert.NotEqual(t, loginResp.AccessToken, refreshResp.AccessToken)

	parseToken := func(tokenString string) jwt.MapClaims {
		t.Helper()

		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.GetString("jwt.secret")), nil
		})
		assert.NoError(t, parseErr)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		return claims
	}

	loginClaims := parseToken(loginResp.AccessToken)
	refreshClaims := parseToken(refreshResp.AccessToken)

	assert.NotEmpty(t, loginClaims["jti"])
	assert.NotEmpty(t, refreshClaims["jti"])
	assert.NotEqual(t, loginClaims["jti"], refreshClaims["jti"])

	mockRepo.AssertExpectations(t)
}

func TestForgetPassword_ExistingEmail(t *testing.T) {
	u, mockRepo, mockMailer, _ := setupTest()

	req := &user.ForgetPasswordRequest{
		Email: "existing@example.com",
	}

	dummyUser := &user.User{ID: "user-123", Email: req.Email}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(dummyUser, nil)
	mockRepo.On("SavePasswordResetToken", mock.Anything, mock.AnythingOfType("*user.PasswordResetToken")).Return(nil)
	mockMailer.On("SendPasswordResetEmail", mock.Anything, req.Email, mock.AnythingOfType("string")).Return(nil)

	err := u.ForgetPassword(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockMailer.AssertExpectations(t)
}

func TestForgetPassword_UnknownEmail(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.ForgetPasswordRequest{Email: "unknown@example.com"}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)

	err := u.ForgetPassword(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestResetPasswordWithToken_Success(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.ResetPasswordWithTokenRequest{
		Token:           "550e8400-e29b-41d4-a716-446655440000",
		NewPassword:     "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	resetToken := &user.PasswordResetToken{
		Token:     req.Token,
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	dummyUser := &user.User{ID: resetToken.UserID, Email: "existing@example.com", Password: "oldhash"}

	mockRepo.On("FindPasswordResetToken", mock.Anything, req.Token).Return(resetToken, nil)
	mockRepo.On("FindByID", mock.Anything, resetToken.UserID).Return(dummyUser, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
	mockRepo.On("DeletePasswordResetToken", mock.Anything, req.Token).Return(nil)

	err := u.ResetPasswordWithToken(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestResetPasswordWithToken_InvalidToken(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	req := &user.ResetPasswordWithTokenRequest{
		Token:           "550e8400-e29b-41d4-a716-446655440000",
		NewPassword:     "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	mockRepo.On("FindPasswordResetToken", mock.Anything, req.Token).Return(nil, nil)

	err := u.ResetPasswordWithToken(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, user.ErrInvalidPasswordResetToken, err)
}

func TestLogin_MissingJWTSecret(t *testing.T) {
	// Setup Manual Khusus case ini (karena butuh config rusak)
	mockRepo := new(MockRepository)
	log := logrus.New()
	log.SetOutput(io.Discard)
	validate := validator.New()

	// CONFIG RUSAK (Secret Kosong)
	cfg := viper.New()
	cfg.Set("jwt.secret", "")
	mockMailer := new(MockMailer)

	u := user.NewUseCase(mockRepo, log, validate, cfg, mockMailer)

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	dummyUser := &user.User{
		Email:    "test@example.com",
		Password: string(hashedPwd),
	}

	req := &user.LoginRequest{Email: "test@example.com", Password: "pass"}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(dummyUser, nil)

	// Action
	resp, err := u.Login(context.Background(), req)

	// Expectation: Harusnya error internal server (daripada panic)
	assert.Error(t, err)
	assert.Equal(t, user.ErrInternalServer, err)
	assert.Nil(t, resp)
}

// ==========================================
// 5. GROUP: GET ME TESTS (PROFILE)
// ==========================================

func TestGetMe_Success(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	userID := "user-uuid-123"
	dummyUser := &user.User{
		ID:        userID,
		Username:  "architect",
		Email:     "architect@example.com",
		CreatedAt: time.Now(),
	}

	// Expectation: Repo FindByID dipanggil dengan ID yang benar
	mockRepo.On("FindByID", mock.Anything, userID).Return(dummyUser, nil)

	// Action
	resp, err := u.GetMe(context.Background(), userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, dummyUser.ID, resp.ID)
	assert.Equal(t, dummyUser.Email, resp.Email)
}

func TestGetMe_UserNotFound(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	userID := "ghost-uuid"

	// Expectation: User tidak ditemukan (return nil, nil)
	mockRepo.On("FindByID", mock.Anything, userID).Return(nil, nil)

	// Action
	resp, err := u.GetMe(context.Background(), userID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "user not found", err.Error())
}

func TestGetMe_RepositoryError(t *testing.T) {
	u, mockRepo, _, _ := setupTest()

	userID := "error-uuid"

	// Expectation: DB Error
	mockRepo.On("FindByID", mock.Anything, userID).Return(nil, errors.New("db connection failed"))

	// Action
	resp, err := u.GetMe(context.Background(), userID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, user.ErrInternalServer, err)
	assert.Nil(t, resp)
}
