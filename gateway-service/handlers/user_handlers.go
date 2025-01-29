package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"your-repo/gateway-service/internal/errors"
	"your-repo/gateway-service/pb"

	"google.golang.org/grpc"
)

func GetUserById(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.GetUserByIdRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		client := pb.NewUserServiceClient(conn)
		user, err := client.GetUserById(context.Background(), &req)
		if err != nil {
			errors.RespondWithError(w, http.StatusNotFound, "Could not find user")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	}
}
func GetUsers(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := pb.NewUserServiceClient(conn)
		users, err := client.GetUsers(context.Background(), &pb.GetUsersRequest{})
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, "Could not find users")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users.Users)
	}
}
