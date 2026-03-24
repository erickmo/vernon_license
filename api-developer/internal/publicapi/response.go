// Package publicapi mengimplementasikan 2 public HTTP endpoint Vernon License:
// POST /api/v1/register dan GET /api/v1/validate.
package publicapi

import (
	"encoding/json"
	"net/http"
)

// APIError merepresentasikan error dengan kode dan pesan yang terstruktur.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse adalah format JSON untuk semua response error public API.
type ErrorResponse struct {
	Valid  bool     `json:"valid"`
	Error  APIError `json:"error"`
}

// WriteError menulis error response JSON ke ResponseWriter dengan HTTP status yang diberikan.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Valid: false,
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

// WriteJSON menulis success response JSON ke ResponseWriter dengan HTTP status yang diberikan.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
