package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL    string
	Port           string
	UserServiceURL string
	JWTSecret      string
}

func LoadConfig() (Config, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	databaseURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	return Config{
		DatabaseURL:    databaseURL,
		Port:           port,
		UserServiceURL: userServiceURL,
		JWTSecret:      jwtSecret,
	}, nil
}
