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

	"your-repo/task-service/config"
	"your-repo/task-service/db"
	"your-repo/task-service/logger"
	"your-repo/task-service/migrations"
	"your-repo/task-service/pb"
	"your-repo/task-service/service"
)

type taskServer struct {
	pb.UnimplementedTaskServiceServer
	db  *sql.DB
	cfg config.Config
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	resp, err := service.CreateTask(s.db, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return resp, nil
}
func (s *taskServer) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	resp, err := service.GetTasks(s.db, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return resp, nil
}
func (s *taskServer) GetTaskById(ctx context.Context, req *pb.GetTaskByIdRequest) (*pb.Task, error) {
	resp, err := service.GetTaskById(s.db, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return resp, nil
}
func (s *taskServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	resp, err := service.UpdateTask(s.db, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return resp, nil
}
func (s *taskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	resp, err := service.DeleteTask(s.db, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete task: %w", err)
	}
	return resp, nil
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

	log.Println(fmt.Sprintf("task service starting at %s", cfg.Port))
	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, &taskServer{db: database, cfg: cfg})
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Error running grpc server", err)
		os.Exit(1)
	}
}
