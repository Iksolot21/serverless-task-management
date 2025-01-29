package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Iksolot21/serverless-task-management/gateway-service/config"
	"github.com/Iksolot21/serverless-task-management/gateway-service/internal/errors"
	"github.com/Iksolot21/serverless-task-management/gateway-service/pb"

	"google.golang.org/grpc"
)

func AuthMiddleware(conn *grpc.ClientConn, cfg config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}
		token := parts[1]
		client := pb.NewAuthServiceClient(conn)
		res, err := client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		if !res.IsValid {
			errors.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		ctx := context.WithValue(r.Context(), "userId", res.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
}
