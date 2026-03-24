// Package middleware menyediakan HTTP middleware untuk Vernon App.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/flashlab/vernon-license/internal/domain"
	jwtpkg "github.com/flashlab/vernon-license/pkg/jwt"
	"go.uber.org/zap"
)

// contextKey adalah tipe private untuk context keys, menghindari collision.
type contextKey string

const contextKeyUser contextKey = "user_claims"

// roleWeight memetakan role ke nilai numerik untuk perbandingan hierarki.
var roleWeight = map[string]int{
	"sales":         1,
	"project_owner": 2,
	"superuser":     3,
}

// AuthMiddleware memvalidasi JWT dari Authorization: Bearer header.
// Set user claims ke request context dengan key contextKeyUser.
// Mengembalikan 401 jika token tidak ada atau tidak valid.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				writeUnauthorized(w, domain.ErrAuthInvalidCredentials.Error())
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtpkg.Verify(tokenString, jwtSecret)
			if err != nil {
				writeUnauthorized(w, domain.ErrAuthTokenExpired.Error())
				return
			}

			ctx := context.WithValue(r.Context(), contextKeyUser, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole memastikan user memiliki salah satu role yang dibutuhkan.
// Hierarki role: superuser > project_owner > sales.
// Mengembalikan 403 jika role tidak mencukupi.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := UserFromContext(r.Context())
			if !ok {
				writeUnauthorized(w, domain.ErrAuthInvalidCredentials.Error())
				return
			}

			userWeight := roleWeight[claims.Role]
			for _, role := range roles {
				if roleWeight[role] <= userWeight {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeForbidden(w, domain.ErrAuthInsufficientRole.Error())
		})
	}
}

// UserFromContext mengambil JWT claims dari request context.
// Mengembalikan claims dan true jika ada, atau nil dan false jika tidak ada.
func UserFromContext(ctx context.Context) (*jwtpkg.Claims, bool) {
	claims, ok := ctx.Value(contextKeyUser).(*jwtpkg.Claims)
	return claims, ok
}

// writeUnauthorized menulis respons HTTP 401.
func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":{"code":"` + msg + `","message":"Unauthorized"}}`))
}

// writeForbidden menulis respons HTTP 403.
func writeForbidden(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":{"code":"` + msg + `","message":"Forbidden"}}`))
}

// NewAuthLogger membuat logger untuk middleware. Digunakan secara opsional.
func NewAuthLogger(logger *zap.Logger) *zap.Logger {
	return logger.Named("auth_middleware")
}
