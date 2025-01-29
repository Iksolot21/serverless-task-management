package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Iksolot21/serverless-task-management/gateway-service/internal/errors"
	"github.com/Iksolot21/serverless-task-management/gateway-service/pb"

	"google.golang.org/grpc"
)

func SendNotification(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.SendNotificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		client := pb.NewNotificationServiceClient(conn)
		res, err := client.SendNotification(context.Background(), &req)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, "Could not send notification")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
