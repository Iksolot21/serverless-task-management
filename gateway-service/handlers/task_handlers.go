package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Iksolot21/serverless-task-management/gateway-service/config"
	"github.com/Iksolot21/serverless-task-management/gateway-service/internal/errors"
	"github.com/Iksolot21/serverless-task-management/gateway-service/pb"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func GetTasks(conn *grpc.ClientConn, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := pb.NewTaskServiceClient(conn)
		token := r.Header.Get("Authorization")
		if token == "" {
			errors.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}
		token = token[7:]
		authClient := pb.NewAuthServiceClient(conn)
		res, err := authClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		if !res.IsValid {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		tasks, err := client.GetTasks(context.Background(), &pb.GetTasksRequest{UserId: res.UserId})
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, "Could not get tasks")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks.Tasks)
	}
}

func CreateTask(conn *grpc.ClientConn, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		token := r.Header.Get("Authorization")
		if token == "" {
			errors.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}
		token = token[7:]
		authClient := pb.NewAuthServiceClient(conn)
		res, err := authClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		if !res.IsValid {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		req.UserId = res.UserId
		client := pb.NewTaskServiceClient(conn)
		_, err = client.CreateTask(context.Background(), &req)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, "Could not create task")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Task created successfully"})
	}
}

func GetTaskById(conn *grpc.ClientConn, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid task id")
			return
		}
		token := r.Header.Get("Authorization")
		if token == "" {
			errors.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}
		token = token[7:]
		authClient := pb.NewAuthServiceClient(conn)
		res, err := authClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		if !res.IsValid {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		client := pb.NewTaskServiceClient(conn)
		task, err := client.GetTaskById(context.Background(), &pb.GetTaskByIdRequest{Id: int32(id), UserId: res.UserId})
		if err != nil {
			errors.RespondWithError(w, http.StatusNotFound, "Could not get task")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}
func PatchTaskById(conn *grpc.ClientConn, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid task id")
			return
		}

		var req pb.UpdateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		req.Id = int32(id)

		token := r.Header.Get("Authorization")
		if token == "" {
			errors.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}
		token = token[7:]
		authClient := pb.NewAuthServiceClient(conn)
		res, err := authClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		if !res.IsValid {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		req.UserId = res.UserId

		client := pb.NewTaskServiceClient(conn)
		_, err = client.UpdateTask(context.Background(), &req)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, "Could not update task")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Task updated successfully"})
	}
}
func DeleteTaskById(conn *grpc.ClientConn, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid task id")
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			errors.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}
		token = token[7:]
		authClient := pb.NewAuthServiceClient(conn)
		res, err := authClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		if !res.IsValid {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		client := pb.NewTaskServiceClient(conn)
		_, err = client.DeleteTask(context.Background(), &pb.DeleteTaskRequest{Id: int32(id), UserId: res.UserId})
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, "Could not delete task")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})

	}
}
