package mailer

import (
	"net/smtp"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpAddr = "smtp.gmail.com:587"
)

type SMTP struct {
	Auth smtp.Auth
}

func New(username string, password string) *SMTP {
	return &SMTP{
		Auth: smtp.PlainAuth("", username, password, smtpHost),
	}
}

func (s *SMTP) Send(from string, to []string, msg []byte) error {
	return smtp.SendMail(smtpAddr, s.Auth, from, to, msg)
}
