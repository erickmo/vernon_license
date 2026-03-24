// Package middleware menyediakan HTTP middleware untuk Vernon License API.
package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

// NewRateLimiter mengembalikan middleware rate limiting berbasis IP.
// Membatasi setiap IP address maksimal requestsPerMinute request per menit.
// Jika limit terlampaui, server merespons dengan HTTP 429.
func NewRateLimiter(requestsPerMinute int) func(http.Handler) http.Handler {
	return httprate.LimitByIP(requestsPerMinute, time.Minute)
}
