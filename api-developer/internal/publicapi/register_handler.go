package publicapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/config"
	"github.com/flashlab/vernon-license/internal/domain"
)

// registerRequest adalah body JSON dari POST /api/v1/register.
type registerRequest struct {
	ProductSlug     string `json:"product_slug"`
	InstanceURL     string `json:"instance_url"`
	InstanceName    string `json:"instance_name"`
	ProvisionAPIKey string `json:"provision_api_key"`
}

// registerResponse adalah response 201 dari POST /api/v1/register.
type registerResponse struct {
	LicenseKey    string `json:"license_key"`
	Product       string `json:"product"`
	CheckInterval string `json:"check_interval"`
	Valid         bool   `json:"valid"`
	Message       string `json:"message"`
}

// RegisterHandler menangani POST /api/v1/register.
type RegisterHandler struct {
	licenses  domain.LicenseRepository
	products  domain.ProductRepository
	auditLogs domain.AuditLogRepository
	cfg       *config.Config
	log       *zap.Logger
}

// NewRegisterHandler membuat instance RegisterHandler baru dengan dependencies yang diperlukan.
func NewRegisterHandler(
	licenses domain.LicenseRepository,
	products domain.ProductRepository,
	auditLogs domain.AuditLogRepository,
	cfg *config.Config,
	log *zap.Logger,
) *RegisterHandler {
	return &RegisterHandler{
		licenses:  licenses,
		products:  products,
		auditLogs: auditLogs,
		cfg:       cfg,
		log:       log,
	}
}

// Handle memproses request POST /api/v1/register.
//
// Alur:
//  1. Parse dan validasi request body (semua field wajib).
//  2. Cari license berdasarkan provision_api_key + product_slug.
//  3. Jika tidak ditemukan → 403 INVALID_API_KEY.
//  4. Jika sudah registered → 409 ALREADY_REGISTERED.
//  5. Update registration (instance_url, instance_name, is_registered = true).
//  6. Update last_pull_at.
//  7. Buat audit log.
//  8. Response 201 dengan license_key dan validity.
func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.ProductSlug == "" || req.InstanceURL == "" || req.InstanceName == "" || req.ProvisionAPIKey == "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION_FAILED", "All fields are required: product_slug, instance_url, instance_name, provision_api_key")
		return
	}

	ctx := r.Context()

	// Step 2: cari license berdasarkan provision_api_key + product_slug
	license, err := h.licenses.FindByProvisionKey(ctx, req.ProvisionAPIKey, req.ProductSlug)
	if err != nil {
		if errors.Is(err, domain.ErrLicenseNotFound) {
			WriteError(w, http.StatusForbidden, "INVALID_API_KEY", "Invalid provision API key")
			return
		}
		h.log.Error("register: FindByProvisionKey", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	// Step 4: cek apakah sudah registered
	if license.IsRegistered {
		WriteError(w, http.StatusConflict, "ALREADY_REGISTERED", "Instance URL already registered")
		return
	}

	// Step 5: update registration
	if err := h.licenses.UpdateRegistration(ctx, license.ID, req.InstanceURL, req.InstanceName); err != nil {
		h.log.Error("register: UpdateRegistration", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	// Step 6: update last_pull_at
	if err := h.licenses.UpdateLastPullAt(ctx, license.ID); err != nil {
		// Non-fatal: log saja, tidak gagalkan request
		h.log.Warn("register: UpdateLastPullAt", zap.Error(err))
	}

	// Step 7: buat audit log — system actor (zero UUID) untuk operasi yang tidak melalui App
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		EntityType: "license",
		EntityID:   license.ID,
		Action:     "client_registered",
		ActorID:    uuid.Nil,
		ActorName:  "system",
		Changes:    nil,
		Metadata:   nil,
	}
	if err := h.auditLogs.Create(ctx, auditLog); err != nil {
		// Non-fatal: log saja
		h.log.Warn("register: Create audit log", zap.Error(err))
	}

	// Step 8: tentukan product name untuk response
	product, err := h.products.FindByID(ctx, license.ProductID)
	productName := req.ProductSlug
	if err == nil {
		productName = product.Name
	}

	// Tentukan validity berdasarkan status license
	valid := license.IsValid()

	// Gunakan check_interval dari license, fallback ke config
	checkInterval := license.CheckInterval
	if checkInterval == "" {
		checkInterval = h.cfg.LicenseCheckInterval
	}

	WriteJSON(w, http.StatusCreated, registerResponse{
		LicenseKey:    license.LicenseKey,
		Product:       productName,
		CheckInterval: checkInterval,
		Valid:          valid,
		Message:       "License registered successfully",
	})
}
