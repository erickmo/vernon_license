//go:build !wasm

// Package handler menyediakan HTTP handler untuk internal Vernon App API.
// Handler ini dilindungi JWT middleware kecuali setup endpoints.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/flashlab/vernon-license/internal/domain"
	appmiddleware "github.com/flashlab/vernon-license/internal/middleware"
	"github.com/flashlab/vernon-license/internal/service"
	"go.uber.org/zap"
)

// AuthHandler menangani endpoint autentikasi internal Vernon App.
type AuthHandler struct {
	authSvc *service.AuthService
	logger  *zap.Logger
}

// NewAuthHandler membuat instance AuthHandler baru.
func NewAuthHandler(authSvc *service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		logger:  logger,
	}
}

// loginRequest adalah body JSON untuk POST /api/internal/auth/login.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// userDTO adalah representasi user yang aman untuk dikembalikan ke client.
type userDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

// loginResponse adalah response dari POST /api/internal/auth/login.
type loginResponse struct {
	Token string  `json:"token"`
	User  userDTO `json:"user"`
}

// Login menangani POST /api/internal/auth/login.
// Tidak memerlukan auth middleware.
// Body: {email, password} → return {token, user: {id, name, role, email}}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "email and password are required")
		return
	}

	token, user, err := h.authSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrAuthInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email atau password salah")
			return
		}
		h.logger.Error("AuthHandler.Login", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{
		Token: token,
		User: userDTO{
			ID:    user.ID.String(),
			Name:  user.Name,
			Role:  user.Role,
			Email: user.Email,
		},
	})
}

// GetMe menangani GET /api/internal/auth/me.
// Memerlukan auth middleware — mengambil user dari JWT context.
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	userUUID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	user, err := h.authSvc.GetMe(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "USER_NOT_FOUND", "User tidak ditemukan")
			return
		}
		h.logger.Error("AuthHandler.GetMe", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, userDTO{
		ID:    user.ID.String(),
		Name:  user.Name,
		Role:  user.Role,
		Email: user.Email,
	})
}
