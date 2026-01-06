package infra

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	// 1. Setup Timeouts & Limits (Security Best Practice)
	// Default ke nilai aman jika config tidak ada
	readTimeout := time.Second * 30
	writeTimeout := time.Second * 30
	idleTimeout := time.Second * 60
	bodyLimit := 10 * 1024 * 1024 // 10 MB limit untuk mencegah DoS via payload besar

	// Jika Anda ingin mengaturnya via Viper nanti:
	// if config.IsSet("server.read_timeout") { ... }

	app := fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      config.GetBool("web.prefork"),

		// Security: Timeouts
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,

		// Security: Limits
		BodyLimit: bodyLimit,

		// Optimization: Mengurangi alokasi memori untuk header
		DisableStartupMessage: false,
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		// Default status code: 500 Internal Server Error
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		// 1. Deteksi apakah ini error standar Fiber (4xx) atau error validasi kita
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		// 2. Logic Masking untuk Production (Security)
		// Jika error 500 (Server Error), JANGAN kirim detail error asli ke client.
		// Pesan asli "db connection failed" hanya boleh ada di logs (stdout/file), bukan di JSON response.
		if code >= fiber.StatusInternalServerError {
			// Kita override pesan error asli agar tidak bocor
			// Log error asli dilakukan oleh Middleware Logger, bukan disini agar SRP terjaga
			message = "Internal Server Error"
		} else {
			// Jika error < 500 (misal 400 Bad Request, 401 Unauthorized),
			// kita boleh menampilkan pesan aslinya (misal: "password required")
			message = err.Error()
		}

		// 3. Return JSON Response yang Konsisten
		return ctx.Status(code).JSON(fiber.Map{
			"error": message,
			// Opsional: Tambahkan Request ID jika Anda menggunakan middleware RequestID untuk tracing
			"request_id": ctx.GetRespHeader("X-Request-ID"),
		})
	}
}
