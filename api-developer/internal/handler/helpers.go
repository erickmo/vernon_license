//go:build !wasm

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// apiError merepresentasikan error terstruktur.
type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// errorResponse adalah format JSON untuk semua response error.
type errorResponse struct {
	Error apiError `json:"error"`
}

// writeError menulis error response JSON ke ResponseWriter.
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: apiError{
			Code:    code,
			Message: message,
		},
	})
}

// writeJSON menulis success response JSON ke ResponseWriter.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// parseUUID mem-parse string UUID. Mengembalikan error jika tidak valid.
func parseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parseUUID: %w", err)
	}
	return id, nil
}
