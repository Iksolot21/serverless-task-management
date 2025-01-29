package errors

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	apiError := APIError{Code: code, Message: message}
	json.NewEncoder(w).Encode(apiError)
}

func NewBadRequestError(message string) APIError {
	return APIError{Code: http.StatusBadRequest, Message: message}
}
func NewUnauthorizedError(message string) APIError {
	return APIError{Code: http.StatusUnauthorized, Message: message}
}
func NewNotFoundError(message string) APIError {
	return APIError{Code: http.StatusNotFound, Message: message}
}
func NewInternalServerError(message string) APIError {
	return APIError{Code: http.StatusInternalServerError, Message: message}
}
