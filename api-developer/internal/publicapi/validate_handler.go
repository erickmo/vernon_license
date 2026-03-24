package publicapi

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/config"
	"github.com/flashlab/vernon-license/internal/domain"
)

// validateValidResponse adalah response 200 ketika license valid.
type validateValidResponse struct {
	Valid         bool   `json:"valid"`
	LicenseKey    string `json:"license_key"`
	CheckInterval string `json:"check_interval"`
}

// validateInvalidResponse adalah response 200 ketika license tidak valid.
type validateInvalidResponse struct {
	Valid         bool   `json:"valid"`
	LicenseKey    string `json:"license_key"`
	Reason        string `json:"reason"`
	CheckInterval string `json:"check_interval"`
}

// ValidateHandler menangani GET /api/v1/validate.
type ValidateHandler struct {
	licenses domain.LicenseRepository
	cfg      *config.Config
	log      *zap.Logger
}

// NewValidateHandler membuat instance ValidateHandler baru dengan dependencies yang diperlukan.
func NewValidateHandler(
	licenses domain.LicenseRepository,
	cfg *config.Config,
	log *zap.Logger,
) *ValidateHandler {
	return &ValidateHandler{
		licenses: licenses,
		cfg:      cfg,
		log:      log,
	}
}

// Handle memproses request GET /api/v1/validate?key=FL-XXXXXXXX.
//
// Alur:
//  1. Ambil query param "key".
//  2. Cari license berdasarkan key.
//  3. Jika tidak ditemukan → 404 LICENSE_NOT_FOUND.
//  4. Update last_pull_at.
//  5. Tentukan validity dari license.IsValid().
//  6. Response 200 dengan valid, license_key, check_interval, dan reason jika tidak valid.
func (h *ValidateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Query param 'key' is required")
		return
	}

	ctx := r.Context()

	license, err := h.licenses.FindByKey(ctx, key)
	if err != nil {
		if errors.Is(err, domain.ErrLicenseNotFound) {
			WriteError(w, http.StatusNotFound, "LICENSE_NOT_FOUND", "Unknown license key")
			return
		}
		h.log.Error("validate: FindByKey", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	// Update last_pull_at untuk monitoring — non-fatal
	if err := h.licenses.UpdateLastPullAt(ctx, license.ID); err != nil {
		h.log.Warn("validate: UpdateLastPullAt", zap.Error(err))
	}

	// Gunakan check_interval dari license, fallback ke config
	checkInterval := license.CheckInterval
	if checkInterval == "" {
		checkInterval = h.cfg.LicenseCheckInterval
	}

	if license.IsValid() {
		WriteJSON(w, http.StatusOK, validateValidResponse{
			Valid:         true,
			LicenseKey:    license.LicenseKey,
			CheckInterval: checkInterval,
		})
		return
	}

	WriteJSON(w, http.StatusOK, validateInvalidResponse{
		Valid:         false,
		LicenseKey:    license.LicenseKey,
		Reason:        license.ValidReason(),
		CheckInterval: checkInterval,
	})
}
