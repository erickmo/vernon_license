package publicapi

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
)

// validateOTPRequest adalah body JSON dari POST /api/v1/validate_otp.
type validateOTPRequest struct {
	OTP string `json:"otp"`
}

// validateOTPResponse adalah response dari POST /api/v1/validate_otp.
type validateOTPResponse struct {
	Status bool `json:"status"`
}

// ValidateOTPHandler menangani POST /api/v1/validate_otp.
type ValidateOTPHandler struct {
	otpRepo domain.OTPRepository
	log     *zap.Logger
}

// NewValidateOTPHandler membuat instance ValidateOTPHandler baru.
func NewValidateOTPHandler(otpRepo domain.OTPRepository, log *zap.Logger) *ValidateOTPHandler {
	return &ValidateOTPHandler{otpRepo: otpRepo, log: log}
}

// Handle memproses request POST /api/v1/validate_otp.
// Client app memanggil endpoint ini untuk memvalidasi OTP yang dikirim oleh license app.
func (h *ValidateOTPHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req validateOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.OTP == "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION_FAILED", "otp is required")
		return
	}

	if err := h.otpRepo.IsActive(r.Context(), req.OTP); err != nil {
		WriteJSON(w, http.StatusOK, validateOTPResponse{Status: false})
		return
	}

	WriteJSON(w, http.StatusOK, validateOTPResponse{Status: true})
}
