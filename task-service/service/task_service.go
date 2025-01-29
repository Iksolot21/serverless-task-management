package tests

import (
	"context"
	"fmt"
	"testing"

	"your-repo/task-service/config"
	"your-repo/task-service/db"
	"your-repo/task-service/pb"
	"your-repo/task-service/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCreateTask(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file %v", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config %v", err)
	}
	database, err := db.OpenDB(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("Error opening database %v", err)
	}
	defer database.Close()
	task := &pb.CreateTaskRequest{
		UserId:      1,
		Title:       "test_task",
		Description: "test_desc",
	}
	result, err := service.CreateTask(database, task)
	if err != nil {
		t.Fatalf("Create task failed: %v", err)
	}
	if result == nil {
		t.Error("Create task response should not be empty")
	}
}
func TestGetTasks(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file %v", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config %v", err)
	}
	database, err := db.OpenDB(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("Error opening database %v", err)
	}
	defer database.Close()
	task := &pb.CreateTaskRequest{
		UserId:      1,
		Title:       "test_task",
		Description: "test_desc",
	}
	_, err = service.CreateTask(database, task)
	if err != nil {
		t.Fatalf("Create task failed: %v", err)
	}
	tasks, err := service.GetTasks(database, &pb.GetTasksRequest{UserId: 1})
	if err != nil {
		t.Fatalf("Get tasks failed: %v", err)
	}
	if len(tasks.Tasks) == 0 {
		t.Error("Tasks should not be empty")
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
	c := pb.NewTaskServiceClient(conn)
	_, err = c.CreateTask(context.Background(), &pb.CreateTaskRequest{
		UserId:      1,
		Title:       "grpc_test_task",
		Description: "grpc_test_desc",
	})
	if err != nil {
		t.Fatalf("Could not create task: %v", err)
	}
}
