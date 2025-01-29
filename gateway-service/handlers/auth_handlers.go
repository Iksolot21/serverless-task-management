package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Iksolot21/serverless-task-management/gateway-service/config"
	"github.com/Iksolot21/serverless-task-management/gateway-service/internal/errors"
	"github.com/Iksolot21/serverless-task-management/gateway-service/pb"

	"google.golang.org/grpc"
)

func RegisterUser(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		client := pb.NewAuthServiceClient(conn)
		res, err := client.Register(context.Background(), &req)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, "Failed to register user")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": res.Message})

	}
}

func LoginUser(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		client := pb.NewAuthServiceClient(conn)
		res, err := client.Login(context.Background(), &req)
		if err != nil {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": res.Token})
	}
}

func GetCurrentUser(conn *grpc.ClientConn, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			errors.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}
		token = token[7:]

		client := pb.NewAuthServiceClient(conn)
		res, err := client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		if !res.IsValid {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		userClient := pb.NewUserServiceClient(conn)
		user, err := userClient.GetUserById(context.Background(), &pb.GetUserByIdRequest{Id: res.UserId})
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, "Could not get user")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	}
}
