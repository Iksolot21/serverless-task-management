package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Iksolot21/serverless-task-management/auth-service/config"
	"github.com/Iksolot21/serverless-task-management/auth-service/db"
	"github.com/Iksolot21/serverless-task-management/auth-service/pb"
	"github.com/Iksolot21/serverless-task-management/auth-service/service"
	"github.com/Iksolot21/serverless-task-management/auth-service/utils"

	_ "github.com/lib/pq"
)

func TestRegisterUser(t *testing.T) {
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

	user := &pb.RegisterRequest{
		Username: "test",
		Password: "password",
		Email:    "test@test.com",
	}
	result, err := service.RegisterUser(database, user)
	if err != nil {
		t.Fatalf("Register user failed: %v", err)
	}
	if result == "" {
		t.Error("Register user response should not be empty")
	}
}

func TestLoginUser(t *testing.T) {
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

	user := &pb.RegisterRequest{
		Username: "test_login",
		Password: "password",
		Email:    "test_login@test.com",
	}
	_, err = service.RegisterUser(database, user)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	login := &pb.LoginRequest{
		Username: "test_login",
		Password: "password",
	}

	token, err := service.LoginUser(database, login, cfg.JWTSecret)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if token == "" {
		t.Error("Token should not be empty")
	}
}
func TestGenerateShortURL(t *testing.T) {
	shortUrl, err := utils.GenerateRandomString(10)
	if err != nil {
		t.Fatalf("Error generating short url: %v", err)
	}
	if len(shortUrl) == 0 {
		t.Error("Short url should not be empty")
	}
}
func TestValidateToken(t *testing.T) {
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

	user := &pb.RegisterRequest{
		Username: "test_token",
		Password: "password",
		Email:    "test_token@test.com",
	}
	_, err = service.RegisterUser(database, user)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	login := &pb.LoginRequest{
		Username: "test_token",
		Password: "password",
	}

	token, err := service.LoginUser(database, login, cfg.JWTSecret)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if token == "" {
		t.Error("Token should not be empty")
	}

	isValid, _, err := service.ValidateToken(token, cfg.JWTSecret)
	if err != nil {
		t.Fatalf("Validate token failed: %v", err)
	}
	if !isValid {
		t.Error("Token should be valid")
	}

	isValid, _, err = service.ValidateToken("invalid_token", cfg.JWTSecret)
	if err == nil {
		t.Error("Invalid token should return error")
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

	c := pb.NewAuthServiceClient(conn)

	_, err = c.Register(context.Background(), &pb.RegisterRequest{
		Username: "grpc_test",
		Password: "password",
		Email:    "grpc@test.com",
	})
	if err != nil {
		t.Fatalf("Could not register user: %v", err)
	}
}
