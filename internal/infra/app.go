package infra

import (
	"strings"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra/middleware"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/budget"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/transaction"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/user"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB       *pgxpool.Pool
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// ========================================
	// GLOBAL MIDDLEWARE (Order matters!)
	// ========================================

	// 1. Request ID (must be first for tracing)
	config.App.Use(middleware.RequestID())

	// 2. Security Headers
	config.App.Use(middleware.SecurityHeaders())

	// 3. CORS
	allowedOrigins := config.Config.GetStringSlice("security.cors.allowed_origins")
	config.App.Use(middleware.CORS(allowedOrigins))

	// 4. Global Rate Limiter
	config.App.Use(middleware.RateLimiter())

	// 5. Logger Middleware
	logConfig := logger.Config{
		TimeFormat: "2006-01-02T15:04:05-0700",
		TimeZone:   "Asia/Jakarta",
	}

	env := config.Config.GetString("app.env")
	if env == "production" {
		logConfig.Format = "{\"time\":\"${time}\",\"status\":${status},\"method\":\"${method}\",\"path\":\"${path}\",\"latency\":\"${latency}\",\"ip\":\"${ip}\",\"request_id\":\"${locals:request_id}\"}\n"
	} else {
		logConfig.Format = "[${time}] ${status} - ${latency} ${method} ${path} | IP: ${ip} | ReqID: ${locals:request_id}\n"
	}
	config.App.Use(logger.New(logConfig))

	// ========================================
	// HEALTH CHECK ENDPOINTS
	// ========================================
	config.App.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": c.Context().Time(),
		})
	})

	config.App.Get("/ready", func(c *fiber.Ctx) error {
		// Check DB connection
		if err := config.DB.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
				"error":  "database not available",
			})
		}
		return c.JSON(fiber.Map{
			"status": "ready",
		})
	})

	// ========================================
	// MODULE SETUP
	// ========================================

	// Auth Middleware
	authMiddleware := middleware.AuthMiddleware(config.Config)
	authRateLimiter := middleware.AuthRateLimiter()

	// User Module
	userRepo := user.NewRepository(config.DB)
	mailSender := NewSMTPMailer(config.Config, config.Log)
	userUseCase := user.NewUseCase(userRepo, config.Log, config.Validate, config.Config, mailSender)
	userHandler := user.NewHandler(userUseCase)
	userHandler.RegisterRoutes(config.App, authMiddleware, authRateLimiter)

	// Budget Module
	budgetRepo := budget.NewRepository(config.DB)
	budgetUseCase := budget.NewUseCase(budgetRepo, config.Log, config.Validate)
	budgetHandler := budget.NewHandler(budgetUseCase)
	budgetHandler.RegisterRoutes(config.App, authMiddleware)

	// Transaction Module (renamed from History)
	transactionRepo := transaction.NewRepository(config.DB)
	transactionUseCase := transaction.NewUseCase(transactionRepo, config.Log, config.Validate)
	transactionHandler := transaction.NewHandler(transactionUseCase)
	transactionHandler.RegisterRoutes(config.App, authMiddleware)

	// ========================================
	// SWAGGER DOCUMENTATION (Dev/Staging Only)
	// ========================================
	if env == "development" || env == "staging" {
		// Import swagger handler hanya jika diperlukan
		// untuk menghindari import di production
		fiberSwagger := getSwaggerHandler()
		if fiberSwagger != nil {
			config.App.Get("/swagger/*", fiberSwagger)
			config.Log.Info("Swagger UI enabled at /swagger/index.html")
		}
	}
}

// getSwaggerHandler returns swagger handler atau nil
func getSwaggerHandler() fiber.Handler {
	// Dynamic import untuk menghindari include di production build
	// Anda bisa gunakan build tags jika perlu
	type swaggerHandler interface {
		WrapHandler(fiber.Handler) fiber.Handler
	}

	// Untuk sekarang, kita return nil dan perlu manual import
	// di main.go jika environment development
	return nil
}

// Helper untuk parse comma-separated string ke slice
func parseCommaSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
