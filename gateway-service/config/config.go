package config

import (
	"os"
)

type Config struct {
	Port                   string
	AuthServiceURL         string
	TaskServiceURL         string
	UserServiceURL         string
	NotificationServiceURL string
	JWTSecret              string
	FrontendURL            string
}

func LoadConfig() (Config, error) {
	port := os.Getenv("PORT")
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	taskServiceURL := os.Getenv("TASK_SERVICE_URL")
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	frontendURL := os.Getenv("FRONTEND_URL")

	return Config{
		Port:                   port,
		AuthServiceURL:         authServiceURL,
		TaskServiceURL:         taskServiceURL,
		UserServiceURL:         userServiceURL,
		NotificationServiceURL: notificationServiceURL,
		JWTSecret:              jwtSecret,
		FrontendURL:            frontendURL,
	}, nil

}
