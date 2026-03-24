//go:build !wasm

package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
	appmiddleware "github.com/flashlab/vernon-license/internal/middleware"
	"github.com/flashlab/vernon-license/internal/service"
)

// UserHandler menangani HTTP requests untuk User management.
type UserHandler struct {
	svc    *service.UserService
	logger *zap.Logger
}

// NewUserHandler membuat instance UserHandler baru.
func NewUserHandler(svc *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

// userListDTO adalah representasi User untuk list response (tanpa password hash).
type userListDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// createUserRequest adalah body JSON untuk POST /api/internal/users.
type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // "project_owner" | "sales"
}

// setActiveRequest adalah body JSON untuk PUT /api/internal/users/{id}/active.
type setActiveRequest struct {
	IsActive bool `json:"is_active"`
}

// toUserListDTO mengkonversi domain.User ke userListDTO.
func toUserListDTO(u *domain.User) userListDTO {
	return userListDTO{
		ID:       u.ID.String(),
		Name:     u.Name,
		Email:    u.Email,
		Role:     u.Role,
		IsActive: u.IsActive,
	}
}

// List menangani GET /api/internal/users.
// Hanya superuser yang dapat melihat daftar user.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	if claims.Role != "superuser" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Hanya superuser yang dapat melihat daftar user")
		return
	}

	users, err := h.svc.List(r.Context())
	if err != nil {
		h.logger.Error("UserHandler.List", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dtos := make([]userListDTO, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, toUserListDTO(u))
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": dtos})
}

// Create menangani POST /api/internal/users.
// Hanya superuser yang dapat membuat user baru.
// Body: {name, email, password, role}
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	if claims.Role != "superuser" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Hanya superuser yang dapat membuat user")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "name, email, password, dan role wajib diisi")
		return
	}

	user, err := h.svc.Create(r.Context(), service.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrUserEmailExists) {
			writeError(w, http.StatusConflict, "USER_EMAIL_EXISTS", "Email sudah digunakan")
			return
		}
		if errors.Is(err, domain.ErrValidationFailed) {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Data tidak valid")
			return
		}
		h.logger.Error("UserHandler.Create", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, toUserListDTO(user))
}

// SetActive menangani PUT /api/internal/users/{id}/active.
// Hanya superuser yang dapat mengaktifkan atau menonaktifkan user.
// Body: {is_active: bool}
func (h *UserHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	if claims.Role != "superuser" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Hanya superuser yang dapat mengubah status user")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	var req setActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if err := h.svc.SetActive(r.Context(), id, req.IsActive, actorID, claims.Name); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "USER_NOT_FOUND", "User tidak ditemukan")
			return
		}
		h.logger.Error("UserHandler.SetActive", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"is_active": req.IsActive})
}
