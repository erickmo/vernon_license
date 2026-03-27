package publicapi

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/config"
	"github.com/flashlab/vernon-license/internal/domain"
	"github.com/flashlab/vernon-license/pkg/licenseutil"
)

// registerRequest adalah body JSON dari POST /api/v1/register.
type registerRequest struct {
	AppName     string `json:"app_name"`
	OTP         string `json:"otp"`
	ClientName  string `json:"client_name"`
	InstanceURL string `json:"instance_url"`
}

// registerResponse adalah response 201 dari POST /api/v1/register.
type registerResponse struct {
	LicenseKey    string `json:"license_key"`
	Product       string `json:"product"`
	CheckInterval string `json:"check_interval"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

// RegisterHandler menangani POST /api/v1/register.
type RegisterHandler struct {
	licenses  domain.LicenseRepository
	products  domain.ProductRepository
	companies domain.CompanyRepository
	auditLogs domain.AuditLogRepository
	cfg       *config.Config
	log       *zap.Logger
	otpRepo   domain.OTPRepository
}

// NewRegisterHandler membuat instance RegisterHandler baru dengan dependencies yang diperlukan.
func NewRegisterHandler(
	licenses domain.LicenseRepository,
	products domain.ProductRepository,
	companies domain.CompanyRepository,
	auditLogs domain.AuditLogRepository,
	cfg *config.Config,
	log *zap.Logger,
	otpRepo domain.OTPRepository,
) *RegisterHandler {
	return &RegisterHandler{
		licenses:  licenses,
		products:  products,
		companies: companies,
		auditLogs: auditLogs,
		cfg:       cfg,
		log:       log,
		otpRepo:   otpRepo,
	}
}

// Handle memproses request POST /api/v1/register.
//
// Alur:
//  1. Parse dan validasi request body (semua field wajib).
//  2. Cek product ada berdasarkan app_name (slug).
//  3. Cek OTP aktif di tabel otp.
//  4. Find-or-create company berdasarkan client_name.
//  5. Cek kombinasi company+product belum ada — jika ada return ALREADY_REGISTERED.
//  6. Buat license baru dengan status "pending", company_id, client_app_ip, instance_url.
//  7. Response 201 dengan license_key dan status.
func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.AppName == "" || req.OTP == "" || req.ClientName == "" || req.InstanceURL == "" {
		WriteError(w, http.StatusBadRequest, "VALIDATION_FAILED", "All fields are required: app_name, otp, client_name, instance_url")
		return
	}

	ctx := r.Context()

	// Capture client IP dari request (Chi RealIP middleware sudah handle X-Forwarded-For)
	clientIP := r.RemoteAddr
	// Strip port jika ada (e.g. "192.168.1.1:54321" → "192.168.1.1")
	if host, _, err := net.SplitHostPort(clientIP); err == nil {
		clientIP = host
	}

	// Step 2: validasi product ada berdasarkan app_name (slug)
	product, err := h.products.FindBySlug(ctx, req.AppName)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			WriteError(w, http.StatusForbidden, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}
		h.log.Error("register: FindBySlug", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	// Step 3: validasi OTP aktif
	if err := h.otpRepo.IsActive(ctx, req.OTP); err != nil {
		WriteError(w, http.StatusForbidden, "INVALID_CLIENT_CODE", "Invalid or expired OTP")
		return
	}

	// Step 4: find-or-create company berdasarkan client_name
	company, err := h.companies.FindByName(ctx, req.ClientName)
	if err != nil {
		if !errors.Is(err, domain.ErrCompanyNotFound) {
			h.log.Error("register: FindByName", zap.Error(err))
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			return
		}
		// Company belum ada — buat baru
		company = &domain.Company{
			ID:        uuid.New(),
			Name:      req.ClientName,
			CreatedBy: nil,
		}
		if err := h.companies.Create(ctx, company); err != nil {
			h.log.Error("register: Create company", zap.Error(err))
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			return
		}
	}

	// Step 5: cek kombinasi company+product sudah ada
	if _, err := h.licenses.FindByCompanyAndProduct(ctx, company.ID, product.ID); err == nil {
		WriteError(w, http.StatusConflict, "ALREADY_REGISTERED", "This company is already registered for this product")
		return
	} else if !errors.Is(err, domain.ErrLicenseNotFound) {
		h.log.Error("register: FindByCompanyAndProduct", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	// Step 6: buat license baru dengan status pending
	checkInterval := h.cfg.LicenseCheckInterval
	if checkInterval == "" {
		checkInterval = "6h"
	}

	licenseKey, err := licenseutil.GenerateLicenseKey()
	if err != nil {
		h.log.Error("register: GenerateLicenseKey", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	instanceURL := req.InstanceURL
	clientName := req.ClientName

	license := &domain.ClientLicense{
		ID:            uuid.New(),
		LicenseKey:    licenseKey,
		ProductID:     product.ID,
		ProductSlug:   product.Slug,
		CompanyID:     &company.ID,
		Plan:          "standard",
		Status:        "pending",
		Modules:       []string{},
		Apps:          []string{},
		InstanceURL:   &instanceURL,
		InstanceName:  &clientName,
		ClientAppIP:   &clientIP,
		CheckInterval: checkInterval,
		IsRegistered:  true,
		CreatedBy:     nil,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.licenses.Create(ctx, license); err != nil {
		h.log.Error("register: Create license", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	// Audit log — system actor
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
		h.log.Warn("register: Create audit log", zap.Error(err))
	}

	WriteJSON(w, http.StatusCreated, registerResponse{
		LicenseKey:    license.LicenseKey,
		Product:       product.Name,
		CheckInterval: checkInterval,
		Status:        "pending",
		Message:       "Registration received. License is pending approval.",
	})
}
