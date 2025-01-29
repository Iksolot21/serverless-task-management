package service

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/Iksolot21/serverless-task-management/notification-service/pb"
)

type EmailSender interface {
	SendEmail(toEmail, subject, body string) error
}

type SmtpEmailSender struct {
	config struct {
		host      string
		port      string
		user      string
		password  string
		fromEmail string
	}
}

func NewSmtpEmailSender(host, port, user, password, fromEmail string) *SmtpEmailSender {
	return &SmtpEmailSender{
		config: struct {
			host      string
			port      string
			user      string
			password  string
			fromEmail string
		}{host: host, port: port, user: user, password: password, fromEmail: fromEmail},
	}
}
func (s *SmtpEmailSender) SendEmail(toEmail, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.user, s.config.password, s.config.host)
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body + "\r\n")

	addr := fmt.Sprintf("%s:%s", s.config.host, s.config.port)
	err := smtp.SendMail(addr, auth, s.config.fromEmail, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil

}

func SendNotification(ctx context.Context, req *pb.SendNotificationRequest, emailSender EmailSender) (*pb.SendNotificationResponse, error) {
	err := emailSender.SendEmail(req.ToEmail, req.Subject, req.Body)
	if err != nil {
		return &pb.SendNotificationResponse{Success: false, Message: "Failed to send notification"}, fmt.Errorf("failed to send notification %w", err)
	}
	return &pb.SendNotificationResponse{Success: true, Message: "Notification sent successfully"}, nil
}
