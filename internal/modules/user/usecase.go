package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type Mailer interface {
	SendPasswordResetEmail(ctx context.Context, to, resetURL string) error
}

type UseCase interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetMe(ctx context.Context, userID string) (*UserResponse, error)
	ResetPassword(ctx context.Context, userID string, req *ResetPasswordRequest) error
	ForgetPassword(ctx context.Context, req *ForgetPasswordRequest) error
	ResetPasswordWithToken(ctx context.Context, req *ResetPasswordWithTokenRequest) error
	RefreshToken(ctx context.Context, userID string, oldToken string, oldTokenExp time.Time) (*LoginResponse, error)
}

type useCase struct {
	repo     Repository
	log      *logrus.Logger
	validate *validator.Validate
	cfg      *viper.Viper
	mailer   Mailer
}

func NewUseCase(repo Repository, log *logrus.Logger, validate *validator.Validate, cfg *viper.Viper, mailer Mailer) UseCase {
	return &useCase{
		repo:     repo,
		log:      log,
		validate: validate,
		cfg:      cfg,
		mailer:   mailer,
	}
}

func (u *useCase) generateAccessToken(user *User) (string, error) {
	tokenTTL := u.cfg.GetDuration("jwt.ttl")
	if tokenTTL == 0 {
		tokenTTL = 24 * time.Hour
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"exp":   now.Add(tokenTTL).Unix(),
		"iat":   now.Unix(),
		"jti":   uuid.NewString(),
		"name":  user.Username,
		"email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := u.cfg.GetString("jwt.secret")
	if jwtSecret == "" {
		u.log.Error("JWT Secret is not configured")
		return "", ErrInternalServer
	}

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		u.log.WithError(err).Error("Failed to sign token")
		return "", ErrInternalServer
	}

	return signedToken, nil
}

// Register Usecase
func (u *useCase) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// 1. Validasi Input
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	// 2. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error("Failed to hash password:", err)
		return nil, ErrInternalServer
	}

	// 3. Construct Entity
	newUser := &User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
		DeletedAt: nil,
	}

	// 4. Simpan ke DB
	if err := u.repo.Save(ctx, newUser); err != nil {
		// Cek error spesifik dari Repository
		if errors.Is(err, ErrEmailTaken) || errors.Is(err, ErrUsernameTaken) {
			return nil, err
		}

		u.log.WithError(err).Error("Failed to save user")
		return nil, ErrInternalServer
	}

	// 5. Return Response
	return &RegisterResponse{
		UserResponse: UserResponse{
			ID:        newUser.ID,
			Username:  newUser.Username,
			Email:     newUser.Email,
			CreatedAt: newUser.CreatedAt,
		},
	}, nil
}

// Login Usecase
func (u *useCase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. Cari User by Email
	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		u.log.WithError(err).Error("Login Failed : Error Finding User")
		return nil, ErrInternalServer
	}

	// 2. Verifikasi Password
	// Jika user nil, return InvalidCredentials
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		u.log.Warnf("Login failed: invalid password for email %s", req.Email)
		return nil, ErrInvalidCredentials
	}

	signedToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// 4. Return Response
	return &LoginResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (u *useCase) GetMe(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		u.log.WithError(err).Error("GetMe: failed to find user")
		return nil, ErrInternalServer
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *useCase) ResetPassword(ctx context.Context, userID string, req *ResetPasswordRequest) error {
	// 1. Validasi Input
	if err := u.validate.Struct(req); err != nil {
		return err
	}

	// 2. Cari User by ID
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		u.log.WithError(err).Error("ResetPassword: failed to find user")
		return ErrInternalServer
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 3. Verifikasi Current Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		u.log.Warnf("ResetPassword failed: invalid current password for user ID %s", userID)
		return ErrInvalidCurrentPassword
	}

	// 4. Hash New Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error("Failed to hash new password:", err)
		return ErrInternalServer
	}

	// 5. Update Password di DB
	user.Password = string(hashed)
	if err := u.repo.Update(ctx, user); err != nil {
		u.log.WithError(err).Error("Failed to update user password")
		return ErrInternalServer
	}

	return nil
}

func (u *useCase) ForgetPassword(ctx context.Context, req *ForgetPasswordRequest) error {
	if err := u.validate.Struct(req); err != nil {
		return err
	}

	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		u.log.WithError(err).Error("ForgetPassword: unable to find user by email")
		return ErrInternalServer
	}

	if user == nil {
		u.log.WithError(err).Error("ForgetPassword: unable to find user by email")
		// For security, don't reveal whether email exists
		return nil
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	passwordResetToken := &PasswordResetToken{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := u.repo.SavePasswordResetToken(ctx, passwordResetToken); err != nil {
		u.log.WithError(err).Error("ForgetPassword: unable to save password reset token")
		return ErrInternalServer
	}

	resetURL := u.cfg.GetString("mail.reset_password_url")
	if resetURL == "" {
		appURL := strings.TrimRight(u.cfg.GetString("app.url"), "/")
		resetURL = fmt.Sprintf("%s/api/users/reset-password?token=%s", appURL, token)
	} else {
		resetURL = fmt.Sprintf("%s?token=%s", strings.TrimRight(resetURL, "/"), token)
	}
	if err := u.mailer.SendPasswordResetEmail(ctx, user.Email, resetURL); err != nil {
		u.log.WithError(err).Error("ForgetPassword: unable to send reset email")
		return ErrInternalServer
	}

	return nil
}

func (u *useCase) ResetPasswordWithToken(ctx context.Context, req *ResetPasswordWithTokenRequest) error {
	if err := u.validate.Struct(req); err != nil {
		return err
	}

	passwordResetToken, err := u.repo.FindPasswordResetToken(ctx, req.Token)
	if err != nil {
		u.log.WithError(err).Error("ResetPasswordWithToken: failed to load token")
		return ErrInternalServer
	}

	if passwordResetToken == nil || time.Now().After(passwordResetToken.ExpiresAt) {
		return ErrInvalidPasswordResetToken
	}

	user, err := u.repo.FindByID(ctx, passwordResetToken.UserID)
	if err != nil {
		u.log.WithError(err).Error("ResetPasswordWithToken: failed to find user")
		return ErrInternalServer
	}

	if user == nil {
		return ErrUserNotFound
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error("Failed to hash new password:", err)
		return ErrInternalServer
	}

	user.Password = string(hashed)
	if err := u.repo.Update(ctx, user); err != nil {
		u.log.WithError(err).Error("ResetPasswordWithToken: failed to update user password")
		return ErrInternalServer
	}

	if err := u.repo.DeletePasswordResetToken(ctx, req.Token); err != nil {
		u.log.WithError(err).Warn("ResetPasswordWithToken: failed to delete token")
	}

	return nil
}

func (u *useCase) RefreshToken(ctx context.Context, userID string, oldToken string, oldTokenExp time.Time) (*LoginResponse, error) {
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		u.log.WithError(err).Error("RefreshToken: failed to find user")
		return nil, ErrInternalServer
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	signedToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}
