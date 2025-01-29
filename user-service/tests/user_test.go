package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Iksolot21/serverless-task-management/user-service/config"
	"github.com/Iksolot21/serverless-task-management/user-service/db"
	"github.com/Iksolot21/serverless-task-management/user-service/pb"
	"github.com/Iksolot21/serverless-task-management/user-service/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGetUserById(t *testing.T) {
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

	user := &pb.User{
		Username: "test_user",
		Email:    "test@test.com",
	}
	query := `INSERT INTO users (username, password, email, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err = database.QueryRow(query, user.Username, "password", user.Email, time.Now()).Scan(&id)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	result, err := service.GetUserById(database, &pb.GetUserByIdRequest{Id: int32(id)})
	if err != nil {
		t.Fatalf("Get user by id failed: %v", err)
	}
	if result == nil {
		t.Error("User should not be nil")
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

	c := pb.NewUserServiceClient(conn)
	resp, err := c.GetUserById(context.Background(), &pb.GetUserByIdRequest{Id: 1})
	if err != nil {
		t.Fatalf("Could not get user: %v", err)
	}
	if resp == nil {
		t.Error("User should not be nil")
	}
}
