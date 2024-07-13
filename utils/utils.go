package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the structure of the error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// RespondWithError sends an error response in a standardized format.
func RespondWithError(w http.ResponseWriter, statusCode int, errMsg string, errDetail string) {
	response := ErrorResponse{
		Error:   errMsg,
		Message: errDetail,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
