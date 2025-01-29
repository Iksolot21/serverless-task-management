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

	"github.com/Iksolot21/serverless-task-management/user-service/config"
	"github.com/Iksolot21/serverless-task-management/user-service/db"
	"github.com/Iksolot21/serverless-task-management/user-service/logger"
	"github.com/Iksolot21/serverless-task-management/user-service/migrations"
	"github.com/Iksolot21/serverless-task-management/user-service/pb"
	"github.com/Iksolot21/serverless-task-management/user-service/service"
)

type userServer struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func (s *userServer) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.User, error) {
	user, err := service.GetUserById(s.db, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
func (s *userServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, err := service.GetUsers(s.db, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return users, nil
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
	log.Println(fmt.Sprintf("user service starting at %s", cfg.Port))
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userServer{db: database})
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Error running grpc server", err)
		os.Exit(1)
	}
}
