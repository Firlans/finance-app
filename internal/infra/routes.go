package infra

import (
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra/middleware"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/accounts"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/categories"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/transactions"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/user"
	"github.com/gofiber/fiber/v2"
)

type routeRegistrar interface {
	RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler, rateLimiter fiber.Handler)
}

func (app *BootstrapConfig) RegisterRoutes() {
	// ========================================
	// MODULE SETUP
	// ========================================

	// Auth Middleware
	authMiddleware := middleware.AuthMiddleware(app.Config)
	authRateLimiter := middleware.AuthRateLimiter()

	mailSender := NewSMTPMailer(app.Config, app.Log)

	modules := []routeRegistrar{
		user.NewHandler(user.NewUseCase(user.NewRepository(app.DB), app.Log, app.Validate, app.Config, mailSender)),
		accounts.NewHandler(accounts.NewUseCase(accounts.NewRepository(app.DB), app.Validate)),
		transactions.NewHandler(transactions.NewUseCase(transactions.NewRepository(app.DB), app.Validate)),
		categories.NewHandler(categories.NewUseCase(categories.NewRepository(app.DB), app.Validate)),
	}

	for _, module := range modules {
		module.RegisterRoutes(app.App, authMiddleware, authRateLimiter)
	}
}
