package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/flashlab/flasherp-developer-api/internal/command/login"
	"github.com/flashlab/flasherp-developer-api/internal/domain/user"
	"github.com/flashlab/flasherp-developer-api/internal/query/get_me"
	"github.com/flashlab/flasherp-developer-api/pkg/middleware"
	"github.com/flashlab/flasherp-developer-api/pkg/querybus"
)

type AuthHandler struct {
	loginHandler *login.Handler
	queryBus     *querybus.QueryBus
}

func NewAuthHandler(loginHandler *login.Handler, queryBus *querybus.QueryBus) *AuthHandler {
	return &AuthHandler{loginHandler: loginHandler, queryBus: queryBus}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	if req.Identifier == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "identifier dan password wajib diisi")
		return
	}

	result, err := h.loginHandler.HandleLogin(r.Context(), login.Login{
		Identifier: req.Identifier,
		Password:   req.Password,
	})
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			respondError(w, http.StatusUnauthorized, "email atau password salah")
			return
		}
		if errors.Is(err, user.ErrInactiveAccount) {
			respondError(w, http.StatusForbidden, "akun tidak aktif")
			return
		}
		respondError(w, http.StatusInternalServerError, "terjadi kesalahan sistem")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	result, err := querybus.Dispatch[*user.User](r.Context(), h.queryBus, get_me.GetMe{UserID: claims.UserID})
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			respondError(w, http.StatusNotFound, "user tidak ditemukan")
			return
		}
		respondError(w, http.StatusInternalServerError, "terjadi kesalahan sistem")
		return
	}

	type meResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	respondJSON(w, http.StatusOK, meResponse{
		ID:    result.ID.String(),
		Name:  result.Name,
		Email: result.Email,
		Role:  result.Role,
	})
}
