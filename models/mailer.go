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

	switch {
	case email.PlainText != "" && email.HTML != "":
		mail.SetBody("text/plain", email.PlainText)
		mail.AddAlternative("text/html", email.HTML)
	case email.PlainText != "":
		mail.SetBody("text/plain", email.PlainText)
	case email.HTML != "":
		mail.SetBody("text/html", email.HTML)
	}

	err := es.dialer.DialAndSend(mail)

	if err != nil {
		return fmt.Errorf("Error in sending mail")
	}

	return nil
}

func (es *EmailService) SendForgetPasswordEmail(to, requestUrl string) error {
	email := Email{
		From:      "support@golangpractice.com",
		To:        to,
		Subject:   "Forget Password Email",
		PlainText: "To reset your password please follow the link : " + requestUrl,
		HTML:      `<p> To reset your password please follow the link : <a href="` + requestUrl + `"> ` + requestUrl + ` </a> </p>`,
	}

	err := es.Send(email)

	if err != nil {
		return fmt.Errorf("Error in sending forget password mail to" + to)
	}

	return nil
}
