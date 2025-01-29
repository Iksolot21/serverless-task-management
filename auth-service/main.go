package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/Iksolot21/serverless-task-management/auth-service/config"
	"github.com/Iksolot21/serverless-task-management/auth-service/db"
	"github.com/Iksolot21/serverless-task-management/auth-service/logger"
	"github.com/Iksolot21/serverless-task-management/auth-service/migrations"
	"github.com/Iksolot21/serverless-task-management/auth-service/pb"
	"github.com/Iksolot21/serverless-task-management/auth-service/service"
)

type authServer struct {
	pb.UnimplementedAuthServiceServer
	db  *sql.DB
	cfg config.Config
}

func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	res, err := service.RegisterUser(s.db, req) //здесь
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	return &pb.RegisterResponse{Message: res}, nil
}

func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, err := service.LoginUser(s.db, req, s.cfg.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to login user: %w", err)
	}
	return &pb.LoginResponse{Token: res}, nil
}
func (s *authServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	res, userId, err := service.ValidateToken(req.Token, s.cfg.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	return &pb.ValidateTokenResponse{IsValid: res, UserId: int32(userId)}, nil
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

	database, err := db.OpenDB(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Error opening database", err)
	}
	defer database.Close()
	err = migrations.RunMigrations(database)
	if err != nil {
		logger.Error("Error running migrations", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		logger.Error("Error create listener", err)
	}
	log.Println(fmt.Sprintf("auth service starting at %s", cfg.Port))

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServer{db: database, cfg: cfg})
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Error running grpc server", err)
		os.Exit(1)
	}
}
