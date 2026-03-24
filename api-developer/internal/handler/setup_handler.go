//go:build !wasm

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/flashlab/vernon-license/internal/service"
	"go.uber.org/zap"
)

// SetupHandler menangani endpoint first-run setup.
type SetupHandler struct {
	setupSvc *service.SetupService
	logger   *zap.Logger
}

// NewSetupHandler membuat instance SetupHandler baru.
func NewSetupHandler(setupSvc *service.SetupService, logger *zap.Logger) *SetupHandler {
	return &SetupHandler{
		setupSvc: setupSvc,
		logger:   logger,
	}
}

// installRequest adalah body JSON untuk POST /api/internal/setup/install.
type installRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// setupStatusResponse adalah response dari GET /api/internal/setup/status.
type setupStatusResponse struct {
	IsSetup bool `json:"is_setup"`
}

// GetStatus menangani GET /api/internal/setup/status.
// Tidak memerlukan auth. Mengembalikan {is_setup: bool}.
func (h *SetupHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	isSetup, err := h.setupSvc.IsSetup(r.Context())
	if err != nil {
		h.logger.Error("SetupHandler.GetStatus", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, setupStatusResponse{IsSetup: isSetup})
}

// Install menangani POST /api/internal/setup/install.
// Tidak memerlukan auth. Membuat superuser pertama.
// Body: {name, email, password} → return {token, user}
func (h *SetupHandler) Install(w http.ResponseWriter, r *http.Request) {
	var req installRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "name, email, and password are required")
		return
	}

	token, user, err := h.setupSvc.Install(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		h.logger.Error("SetupHandler.Install", zap.Error(err))
		writeError(w, http.StatusBadRequest, "SETUP_FAILED", "Setup gagal atau sudah pernah dilakukan")
		return
	}

	writeJSON(w, http.StatusCreated, loginResponse{
		Token: token,
		User: userDTO{
			ID:    user.ID.String(),
			Name:  user.Name,
			Role:  user.Role,
			Email: user.Email,
		},
	})
}
