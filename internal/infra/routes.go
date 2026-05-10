package infra

import (
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra/middleware"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/accounts"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/user"
)

func (app *BootstrapConfig) RegisterRoutes() {

	// ========================================
	// MODULE SETUP
	// ========================================

	// Auth Middleware
	authMiddleware := middleware.AuthMiddleware(app.Config)
	authRateLimiter := middleware.AuthRateLimiter()

	// User Module
	userRepo := user.NewRepository(app.DB)
	mailSender := NewSMTPMailer(app.Config, app.Log)
	userUseCase := user.NewUseCase(userRepo, app.Log, app.Validate, app.Config, mailSender)
	userHandler := user.NewHandler(userUseCase)
	userHandler.RegisterRoutes(app.App, authMiddleware, authRateLimiter)

	// Accounts Module
	accountRepo := accounts.NewRepository(app.DB)
	accountUseCase := accounts.NewUseCase(accountRepo, app.Validate)
	accountHandler := accounts.NewHandler(accountUseCase)
	accountHandler.RegisterRoutes(app.App, authMiddleware, authRateLimiter)
}
