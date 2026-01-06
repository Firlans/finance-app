package infra

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogger(config *viper.Viper) *logrus.Logger {
	log := logrus.New()

	// 1. Output ke Stdout (Terminal) - Standard Cloud Native
	log.SetOutput(os.Stdout)

	// 2. Set Level dari Config (String lebih aman daripada Integer)
	// Contoh config: "info", "debug", "warn", "error"
	lvl, err := logrus.ParseLevel(config.GetString("log.level"))
	if err != nil {
		lvl = logrus.InfoLevel // Default jika config salah
	}
	log.SetLevel(lvl)

	// 3. Formatter Cerdas
	// Jika di Production -> JSON (untuk mesin)
	// Jika di Development -> Text (biar enak dibaca mata manusia warna-warni)
	if config.GetString("app.env") == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	return log
}
