package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"your-repo/gateway-service/config"
	"your-repo/gateway-service/handlers"
	"your-repo/gateway-service/pb"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRegisterUserHandler(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file %v", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config %v", err)
	}

	authConn, err := grpc.Dial(cfg.AuthServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Did not connect: %v", err)
	}
	defer authConn.Close()
	r := mux.NewRouter()
	r.HandleFunc("/auth/register", handlers.RegisterUser(authConn)).Methods("POST")

	reqBody := `{"username":"test","password":"password", "email":"test@test.com"}`
	req, err := http.NewRequest("POST", "/auth/register", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler return wrong status code: got %v, want %v", status, http.StatusCreated)
	}
	var resp map[string]string
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if _, ok := resp["message"]; !ok {
		t.Error("Handler did not return message")
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

	conn, err := grpc.Dial(cfg.AuthServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
