package infra

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/TubagusAldiMY/finance-tracker-app/backend/internal/modules/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type SMTPMailer struct {
	cfg *viper.Viper
	log *logrus.Logger
}

func NewSMTPMailer(cfg *viper.Viper, log *logrus.Logger) user.Mailer {
	return &SMTPMailer{cfg: cfg, log: log}
}

func (m *SMTPMailer) SendPasswordResetEmail(ctx context.Context, to, resetURL string) error {
	smtpHost := m.cfg.GetString("mail.host")
	smtpPort := m.cfg.GetInt("mail.port")
	log.Printf("SMTP Config - Host: %s, Port: %d", smtpHost, smtpPort)
	if smtpPort == 0 {
		smtpPort = 587
	}
	mailUser := m.cfg.GetString("mail.username")
	mailPass := m.cfg.GetString("mail.password")
	mailFrom := m.cfg.GetString("mail.from")
	if mailFrom == "" {
		mailFrom = "noreply@example.com"
	}

	if smtpHost == "" || mailUser == "" || mailPass == "" {
		return fmt.Errorf("smtp configuration is incomplete")
	}

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", mailUser, mailPass, smtpHost)

	subject := "Password Reset Request"
	body := fmt.Sprintf(`To reset your password, click the link below:

%s

If you did not request a password reset, please ignore this email.`, resetURL)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", mailFrom),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
		"",
		body,
	}, "\r\n")

	return smtp.SendMail(addr, auth, mailFrom, []string{to}, []byte(msg))
}
