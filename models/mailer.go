package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "test@gocourse.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	PlainText string
	HTML      string
}

type SMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type EmailService struct {
	DefaultSender string
	dialer        *mail.Dialer
}

func NewMailService(config SMTPConfig) *EmailService {
	return &EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.User, config.Pass),
	}
}

func (es *EmailService) Send(email Email) error {
	mail := mail.NewMessage()
	mail.SetHeader("From", email.From)
	mail.SetHeader("To", email.To)
	mail.SetHeader("Subject", email.Subject)
	mail.SetBody("text/plain", email.PlainText)
	mail.AddAlternative("text/html", email.HTML)

	err := es.dialer.DialAndSend(mail)

	if err != nil {
		return fmt.Errorf("Error in sending mail")
	}

	return nil
}
