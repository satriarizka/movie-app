package mailer

import (
	"fmt"
	"movie-app/internal/config"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
}

func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{
		dialer: gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass),
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "no-reply@bioskop.com")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
