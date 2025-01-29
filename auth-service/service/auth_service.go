package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Iksolot21/serverless-task-management/auth-service/internal/jwt"
	"github.com/Iksolot21/serverless-task-management/auth-service/pb"
)

func RegisterUser(db *sql.DB, req *pb.RegisterRequest) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	query := `INSERT INTO users (username, password, email, created_at) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(query, req.Username, string(hashedPassword), req.Email, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to insert user into database: %w", err)
	}

	return "User registered successfully", nil
}

func LoginUser(db *sql.DB, req *pb.LoginRequest, jwtSecret string) (string, error) {
	var user struct {
		ID       int
		Username string
		Password string
		Email    string
	}

	query := `SELECT id, username, password, email FROM users WHERE username = $1`
	row := db.QueryRow(query, req.Username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}
	token, err := jwt.GenerateJWT(user.ID, user.Username, user.Email, jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return token, nil
}
func ValidateToken(tokenString, jwtSecret string) (bool, int, error) {
	user, err := jwt.ValidateJWT(tokenString, jwtSecret)
	if err != nil {
		return false, 0, fmt.Errorf("invalid token: %w", err)
	}
	return true, user.ID, nil
}
