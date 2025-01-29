package config

import (
	"os"
)

type Config struct {
	Port      string
	SMTPHost  string
	SMTPPort  string
	SMTPUser  string
	SMTPPass  string
	FromEmail string
}

func LoadConfig() (Config, error) {
	port := os.Getenv("PORT")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("FROM_EMAIL")

	return Config{
		Port:      port,
		SMTPHost:  smtpHost,
		SMTPPort:  smtpPort,
		SMTPUser:  smtpUser,
		SMTPPass:  smtpPass,
		FromEmail: fromEmail,
	}, nil
}
