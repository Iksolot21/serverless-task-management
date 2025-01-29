package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId int, username, email string, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, jwtSecret string) (struct {
	ID       int
	Username string
	Email    string
}, error) {
	var user struct {
		ID       int
		Username string
		Email    string
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return user, fmt.Errorf("failed to parse token: %w", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["userId"].(float64)
		if !ok {
			return user, errors.New("invalid token claims")
		}
		username, ok := claims["username"].(string)
		if !ok {
			return user, errors.New("invalid token claims")
		}
		email, ok := claims["email"].(string)
		if !ok {
			return user, errors.New("invalid token claims")
		}

		user.ID = int(id)
		user.Username = username
		user.Email = email

		return user, nil
	}

	return user, errors.New("invalid token")
}
