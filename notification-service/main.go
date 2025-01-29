package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"your-repo/notification-service/config"
	"your-repo/notification-service/pb"
	"your-repo/notification-service/service"
     "your-repo/notification-service/logger"
)

type notificationServer struct {
	pb.UnimplementedNotificationServiceServer
	config config.Config
}

func (s *notificationServer) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
    emailSender := service.NewSmtpEmailSender(s.config.SMTPHost, s.config.SMTPPort, s.config.SMTPUser, s.config.SMTPPass, s.config.FromEmail)
    res, err := service.SendNotification(ctx, req, emailSender)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}
	return res, nil
}


func main() {
    err := godotenv.Load()
    if err != nil {
         logger.Error("Error loading .env file", err)
     }
     cfg, err := config.LoadConfig()
   if err != nil {
         logger.Error("Error loading config", err)
    }
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
	  	 logger.Error("Error create listener", err)
	}
  log.Println(fmt.Sprintf("notification service starting at %s", cfg.Port ))

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, ¬ificationServer{config: cfg})
	if err := grpcServer.Serve(lis); err != nil {
        logger.Error("Error running grpc server", err)
         os.Exit(1)
    }
}