package infra

import (
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/infra/middleware"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/budget"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/history"
	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/user" // Import module User
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

	// Auth Middleware
	authMiddleware := middleware.AuthMiddleware(config.Config)

	// Setup Logger Middleware
	logConfig := logger.Config{
		TimeFormat: "2006-01-02T15:04:05-0700",
		TimeZone:   "Asia/Jakarta",
	}

	// Cek Environment dari Config
	env := config.Config.GetString("app.env")

	if env == "production" {
		// Production: JSON Format (Mesin)
		logConfig.Format = "{\"time\":\"${time}\",\"status\":${status},\"method\":\"${method}\",\"path\":\"${path}\",\"latency\":\"${latency}\"}\n"
	} else {
		// Development: Colorized Text (Manusia)
		// Format default Fiber sudah cukup bagus untuk dev
		logConfig.Format = "[${time}] ${status} - ${latency} ${method} ${path}\n"
	}

	config.App.Use(logger.New(logConfig))

	// User Module
	userRepo := user.NewRepository(config.DB)
	userUseCase := user.NewUseCase(userRepo, config.Log, config.Validate, config.Config)
	userHandler := user.NewHandler(userUseCase)

	userHandler.RegisterRoutes(config.App, authMiddleware)

	budgetRepo := budget.NewRepository(config.DB)
	budgetUseCase := budget.NewUseCase(budgetRepo, config.Log, config.Validate)
	budgetHandler := budget.NewHandler(budgetUseCase)

	budgetHandler.RegisterRoutes(config.App, authMiddleware)

	historyRepo := history.NewRepository(config.DB)
	historyUseCase := history.NewUseCase(historyRepo, config.Log, config.Validate)
	historyHandler := history.NewHandler(historyUseCase)

	historyHandler.RegisterRoutes(config.App, authMiddleware)
}
