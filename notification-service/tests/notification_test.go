package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Iksolot21/serverless-task-management/notification-service/config"
	"github.com/Iksolot21/serverless-task-management/notification-service/pb"
	"github.com/Iksolot21/serverless-task-management/notification-service/service"
)

type MockEmailSender struct {
	Sent bool
}

func (m *MockEmailSender) SendEmail(toEmail, subject, body string) error {
	m.Sent = true
	return nil
}

func TestSendNotification(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file %v", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config %v", err)
	}
	emailSender := &MockEmailSender{Sent: false}

	req := &pb.SendNotificationRequest{
		ToEmail: "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	resp, err := service.SendNotification(context.Background(), req, emailSender)
	if err != nil {
		t.Errorf("SendNotification failed: %v", err)
	}

	if resp.Success != true {
		t.Error("Notification send is not successfull")
	}
	if emailSender.Sent == false {
		t.Error("Email should be sent")
	}

}

func TestGrpcConnection(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file %v", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config %v", err)
	}

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNotificationServiceClient(conn)

	resp, err := c.SendNotification(context.Background(), &pb.SendNotificationRequest{
		ToEmail: "test@test.com",
		Subject: "grpc test",
		Body:    "grpc test body",
	})
	if err != nil {
		t.Fatalf("Could not register user: %v", err)
	}

	if !resp.Success {
		t.Error("GRPC request failed")
	}
}
