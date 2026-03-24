//go:build !wasm

package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
	appmiddleware "github.com/flashlab/vernon-license/internal/middleware"
	"github.com/flashlab/vernon-license/internal/service"
)

// LicenseHandler menangani HTTP request untuk manajemen license internal.
type LicenseHandler struct {
	licenseSvc  *service.LicenseService
	companySvc  *service.CompanyService
	projectSvc  *service.ProjectService
	productSvc  *service.ProductService
	auditSvc    *service.AuditService
	logger      *zap.Logger
}

// NewLicenseHandler membuat instance LicenseHandler baru.
func NewLicenseHandler(
	licenseSvc *service.LicenseService,
	companySvc *service.CompanyService,
	projectSvc *service.ProjectService,
	productSvc *service.ProductService,
	auditSvc *service.AuditService,
	logger *zap.Logger,
) *LicenseHandler {
	return &LicenseHandler{
		licenseSvc: licenseSvc,
		companySvc: companySvc,
		projectSvc: projectSvc,
		productSvc: productSvc,
		auditSvc:   auditSvc,
		logger:     logger,
	}
}

// licenseListItemDTO adalah representasi ringkas license untuk tampilan daftar.
type licenseListItemDTO struct {
	ID           string  `json:"id"`
	LicenseKey   string  `json:"license_key"`
	CompanyName  string  `json:"company_name"`
	ProjectName  string  `json:"project_name"`
	ProductName  string  `json:"product_name"`
	Plan         string  `json:"plan"`
	Status       string  `json:"status"`
	IsRegistered bool    `json:"is_registered"`
	ExpiresAt    *string `json:"expires_at"`
}

// licenseDetailDTO adalah representasi lengkap license untuk tampilan detail.
type licenseDetailDTO struct {
	ID               string   `json:"id"`
	LicenseKey       string   `json:"license_key"`
	CompanyID        string   `json:"company_id"`
	CompanyName      string   `json:"company_name"`
	ProjectID        string   `json:"project_id"`
	ProjectName      string   `json:"project_name"`
	ProductName      string   `json:"product_name"`
	Plan             string   `json:"plan"`
	Status           string   `json:"status"`
	Modules          []string `json:"modules"`
	Apps             []string `json:"apps"`
	ContractAmount   *float64 `json:"contract_amount"`
	Description      string   `json:"description"`
	MaxUsers         *int     `json:"max_users"`
	MaxTransPerMonth *int     `json:"max_trans_per_month"`
	MaxTransPerDay   *int     `json:"max_trans_per_day"`
	MaxItems         *int     `json:"max_items"`
	MaxCustomers     *int     `json:"max_customers"`
	MaxBranches      *int     `json:"max_branches"`
	MaxStorage       *int     `json:"max_storage"`
	ExpiresAt        *string  `json:"expires_at"`
	IsRegistered     bool     `json:"is_registered"`
	InstanceURL      string   `json:"instance_url"`
	InstanceName     string   `json:"instance_name"`
	ProvisionAPIKey  string   `json:"provision_api_key"`
	CheckInterval    string   `json:"check_interval"`
	LastPullAt       *string  `json:"last_pull_at"`
}

// createLicenseRequest adalah body JSON untuk POST /api/internal/licenses.
type createLicenseRequest struct {
	ProjectID        string   `json:"project_id"`
	CompanyID        string   `json:"company_id"`
	ProductID        string   `json:"product_id"`
	Plan             string   `json:"plan"`
	Modules          []string `json:"modules"`
	Apps             []string `json:"apps"`
	ContractAmount   *float64 `json:"contract_amount"`
	Description      *string  `json:"description"`
	MaxUsers         *int     `json:"max_users"`
	MaxTransPerMonth *int     `json:"max_trans_per_month"`
	MaxTransPerDay   *int     `json:"max_trans_per_day"`
	MaxItems         *int     `json:"max_items"`
	MaxCustomers     *int     `json:"max_customers"`
	MaxBranches      *int     `json:"max_branches"`
	MaxStorage       *int     `json:"max_storage"`
	ExpiresAt        *string  `json:"expires_at"`
	CheckInterval    string   `json:"check_interval"`
}

// renewRequest adalah body JSON untuk PUT /api/internal/licenses/{id}/renew.
type renewRequest struct {
	NewExpiresAt *string `json:"new_expires_at"`
}

// updateConstraintsRequest adalah body JSON untuk PUT /api/internal/licenses/{id}/constraints.
type updateConstraintsRequest struct {
	Modules          []string `json:"modules"`
	Apps             []string `json:"apps"`
	MaxUsers         *int     `json:"max_users"`
	MaxTransPerMonth *int     `json:"max_trans_per_month"`
	MaxTransPerDay   *int     `json:"max_trans_per_day"`
	MaxItems         *int     `json:"max_items"`
	MaxCustomers     *int     `json:"max_customers"`
	MaxBranches      *int     `json:"max_branches"`
	MaxStorage       *int     `json:"max_storage"`
	ExpiresAt        *string  `json:"expires_at"`
	CheckInterval    string   `json:"check_interval"`
}

// auditLogDTO adalah representasi audit log untuk response API.
type auditLogDTO struct {
	ID         string `json:"id"`
	Action     string `json:"action"`
	ActorID    string `json:"actor_id"`
	ActorName  string `json:"actor_name"`
	Changes    any    `json:"changes"`
	Metadata   any    `json:"metadata"`
	CreatedAt  string `json:"created_at"`
}

// List menangani GET /api/internal/licenses.
// Mengembalikan semua license dengan info company, project, dan product.
func (h *LicenseHandler) List(w http.ResponseWriter, r *http.Request) {
	licenses, err := h.licenseSvc.List(r.Context())
	if err != nil {
		h.logger.Error("LicenseHandler.List", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	items := make([]licenseListItemDTO, 0, len(licenses))
	for _, l := range licenses {
		item := licenseListItemDTO{
			ID:           l.ID.String(),
			LicenseKey:   l.LicenseKey,
			Plan:         l.Plan,
			Status:       l.Status,
			IsRegistered: l.IsRegistered,
		}

		if company, err := h.companySvc.GetByID(r.Context(), l.CompanyID); err == nil {
			item.CompanyName = company.Name
		}
		if project, err := h.projectSvc.GetByID(r.Context(), l.ProjectID); err == nil {
			item.ProjectName = project.Name
		}
		if product, err := h.productSvc.GetByID(r.Context(), l.ProductID); err == nil {
			item.ProductName = product.Name
		}

		if l.ExpiresAt != nil {
			s := l.ExpiresAt.UTC().Format(time.RFC3339)
			item.ExpiresAt = &s
		}

		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, items)
}

// ListByProject menangani GET /api/internal/projects/{projectID}/licenses.
// Mengembalikan semua license milik project tertentu.
func (h *LicenseHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := chi.URLParam(r, "projectID")
	projectID, err := parseUUID(projectIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid project ID")
		return
	}

	licenses, err := h.licenseSvc.ListByProject(r.Context(), projectID)
	if err != nil {
		h.logger.Error("LicenseHandler.ListByProject", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	items := make([]licenseListItemDTO, 0, len(licenses))
	for _, l := range licenses {
		item := licenseListItemDTO{
			ID:           l.ID.String(),
			LicenseKey:   l.LicenseKey,
			Plan:         l.Plan,
			Status:       l.Status,
			IsRegistered: l.IsRegistered,
		}

		if company, err := h.companySvc.GetByID(r.Context(), l.CompanyID); err == nil {
			item.CompanyName = company.Name
		}
		if project, err := h.projectSvc.GetByID(r.Context(), l.ProjectID); err == nil {
			item.ProjectName = project.Name
		}
		if product, err := h.productSvc.GetByID(r.Context(), l.ProductID); err == nil {
			item.ProductName = product.Name
		}

		if l.ExpiresAt != nil {
			s := l.ExpiresAt.UTC().Format(time.RFC3339)
			item.ExpiresAt = &s
		}

		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, items)
}

// GetByID menangani GET /api/internal/licenses/{id}.
// Mengembalikan detail lengkap license termasuk provision_api_key.
func (h *LicenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := parseUUID(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid license ID")
		return
	}

	l, err := h.licenseSvc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrLicenseNotFound) {
			writeError(w, http.StatusNotFound, "LICENSE_NOT_FOUND", "License tidak ditemukan")
			return
		}
		h.logger.Error("LicenseHandler.GetByID", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dto := h.toLicenseDetailDTO(r, l)
	writeJSON(w, http.StatusOK, dto)
}

// Create menangani POST /api/internal/licenses.
// Membuat license secara langsung (direct create). Hanya project_owner dan superuser.
func (h *LicenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	var req createLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.ProjectID == "" || req.CompanyID == "" || req.ProductID == "" || req.Plan == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "project_id, company_id, product_id, dan plan wajib diisi")
		return
	}

	projectID, err := parseUUID(req.ProjectID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid project_id")
		return
	}
	companyID, err := parseUUID(req.CompanyID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid company_id")
		return
	}
	productID, err := parseUUID(req.ProductID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid product_id")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	svcReq := service.CreateLicenseRequest{
		ProjectID:        projectID,
		CompanyID:        companyID,
		ProductID:        productID,
		Plan:             req.Plan,
		Modules:          req.Modules,
		Apps:             req.Apps,
		ContractAmount:   req.ContractAmount,
		Description:      req.Description,
		MaxUsers:         req.MaxUsers,
		MaxTransPerMonth: req.MaxTransPerMonth,
		MaxTransPerDay:   req.MaxTransPerDay,
		MaxItems:         req.MaxItems,
		MaxCustomers:     req.MaxCustomers,
		MaxBranches:      req.MaxBranches,
		MaxStorage:       req.MaxStorage,
		CheckInterval:    req.CheckInterval,
	}

	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, parseErr := time.Parse(time.RFC3339, *req.ExpiresAt)
		if parseErr != nil {
			// coba format date only
			t, parseErr = time.Parse("2006-01-02", *req.ExpiresAt)
			if parseErr != nil {
				writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Format expires_at tidak valid, gunakan RFC3339 atau YYYY-MM-DD")
				return
			}
		}
		svcReq.ExpiresAt = &t
	}

	license, err := h.licenseSvc.DirectCreate(r.Context(), svcReq, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			writeError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrProductInactive) {
			writeError(w, http.StatusBadRequest, "PRODUCT_INACTIVE", "Product tidak aktif")
			return
		}
		h.logger.Error("LicenseHandler.Create", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dto := h.toLicenseDetailDTO(r, license)
	writeJSON(w, http.StatusCreated, dto)
}

// Activate menangani PUT /api/internal/licenses/{id}/activate.
// Mengubah status license ke "active". Hanya project_owner dan superuser.
func (h *LicenseHandler) Activate(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := parseUUID(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid license ID")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	if err := h.licenseSvc.Activate(r.Context(), id, actorID, claims.Name); err != nil {
		if errors.Is(err, domain.ErrLicenseNotFound) {
			writeError(w, http.StatusNotFound, "LICENSE_NOT_FOUND", "License tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrLicenseInvalidTransition) {
			writeError(w, http.StatusBadRequest, "INVALID_TRANSITION", err.Error())
			return
		}
		h.logger.Error("LicenseHandler.Activate", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "active"})
}

// Suspend menangani PUT /api/internal/licenses/{id}/suspend.
// Mengubah status license ke "suspended". Hanya project_owner dan superuser.
func (h *LicenseHandler) Suspend(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := parseUUID(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid license ID")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	if err := h.licenseSvc.Suspend(r.Context(), id, actorID, claims.Name); err != nil {
		if errors.Is(err, domain.ErrLicenseNotFound) {
			writeError(w, http.StatusNotFound, "LICENSE_NOT_FOUND", "License tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrLicenseInvalidTransition) {
			writeError(w, http.StatusBadRequest, "INVALID_TRANSITION", err.Error())
			return
		}
		h.logger.Error("LicenseHandler.Suspend", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "suspended"})
}

// Renew menangani PUT /api/internal/licenses/{id}/renew.
// Memperbarui license yang expired ke "active". Hanya project_owner dan superuser.
func (h *LicenseHandler) Renew(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := parseUUID(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid license ID")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	var req renewRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	var newExpiresAt *time.Time
	if req.NewExpiresAt != nil && *req.NewExpiresAt != "" {
		t, parseErr := time.Parse(time.RFC3339, *req.NewExpiresAt)
		if parseErr != nil {
			t, parseErr = time.Parse("2006-01-02", *req.NewExpiresAt)
			if parseErr != nil {
				writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Format new_expires_at tidak valid")
				return
			}
		}
		newExpiresAt = &t
	}

	if err := h.licenseSvc.Renew(r.Context(), id, newExpiresAt, actorID, claims.Name); err != nil {
		if errors.Is(err, domain.ErrLicenseNotFound) {
			writeError(w, http.StatusNotFound, "LICENSE_NOT_FOUND", "License tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrLicenseInvalidTransition) {
			writeError(w, http.StatusBadRequest, "INVALID_TRANSITION", err.Error())
			return
		}
		h.logger.Error("LicenseHandler.Renew", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "active"})
}

// UpdateConstraints menangani PUT /api/internal/licenses/{id}/constraints.
// Memperbarui constraint license. Hanya project_owner dan superuser.
func (h *LicenseHandler) UpdateConstraints(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := parseUUID(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid license ID")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	var req updateConstraintsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	svcReq := service.UpdateConstraintsRequest{
		Modules:          req.Modules,
		Apps:             req.Apps,
		MaxUsers:         req.MaxUsers,
		MaxTransPerMonth: req.MaxTransPerMonth,
		MaxTransPerDay:   req.MaxTransPerDay,
		MaxItems:         req.MaxItems,
		MaxCustomers:     req.MaxCustomers,
		MaxBranches:      req.MaxBranches,
		MaxStorage:       req.MaxStorage,
		CheckInterval:    req.CheckInterval,
	}

	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, parseErr := time.Parse(time.RFC3339, *req.ExpiresAt)
		if parseErr != nil {
			t, parseErr = time.Parse("2006-01-02", *req.ExpiresAt)
			if parseErr != nil {
				writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Format expires_at tidak valid")
				return
			}
		}
		svcReq.ExpiresAt = &t
	}

	if err := h.licenseSvc.UpdateConstraints(r.Context(), id, svcReq, actorID, claims.Name); err != nil {
		if errors.Is(err, domain.ErrLicenseNotFound) {
			writeError(w, http.StatusNotFound, "LICENSE_NOT_FOUND", "License tidak ditemukan")
			return
		}
		h.logger.Error("LicenseHandler.UpdateConstraints", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"result": "ok"})
}

// GetAuditLogs menangani GET /api/internal/licenses/{id}/audit.
// Mengembalikan audit log untuk license tertentu.
func (h *LicenseHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := parseUUID(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid license ID")
		return
	}

	logs, err := h.auditSvc.ListByEntity(r.Context(), "license", id)
	if err != nil {
		h.logger.Error("LicenseHandler.GetAuditLogs", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dtos := make([]auditLogDTO, 0, len(logs))
	for _, l := range logs {
		dto := auditLogDTO{
			ID:        l.ID.String(),
			Action:    l.Action,
			ActorID:   l.ActorID.String(),
			ActorName: l.ActorName,
			CreatedAt: l.CreatedAt.UTC().Format(time.RFC3339),
		}
		// Changes dan Metadata adalah json.RawMessage — decode ke any supaya tidak double-escaped.
		var changes any
		if err := json.Unmarshal(l.Changes, &changes); err == nil {
			dto.Changes = changes
		} else {
			dto.Changes = string(l.Changes)
		}
		var metadata any
		if err := json.Unmarshal(l.Metadata, &metadata); err == nil {
			dto.Metadata = metadata
		} else {
			dto.Metadata = string(l.Metadata)
		}
		dtos = append(dtos, dto)
	}

	writeJSON(w, http.StatusOK, dtos)
}

// toLicenseDetailDTO mengonversi domain.ClientLicense ke licenseDetailDTO dengan lookup nama.
// Lookup error diabaikan — nama akan kosong jika lookup gagal.
func (h *LicenseHandler) toLicenseDetailDTO(r *http.Request, l *domain.ClientLicense) licenseDetailDTO {
	dto := licenseDetailDTO{
		ID:               l.ID.String(),
		LicenseKey:       l.LicenseKey,
		CompanyID:        l.CompanyID.String(),
		ProjectID:        l.ProjectID.String(),
		Plan:             l.Plan,
		Status:           l.Status,
		Modules:          l.Modules,
		Apps:             l.Apps,
		ContractAmount:   l.ContractAmount,
		MaxUsers:         l.MaxUsers,
		MaxTransPerMonth: l.MaxTransPerMonth,
		MaxTransPerDay:   l.MaxTransPerDay,
		MaxItems:         l.MaxItems,
		MaxCustomers:     l.MaxCustomers,
		MaxBranches:      l.MaxBranches,
		MaxStorage:       l.MaxStorage,
		IsRegistered:     l.IsRegistered,
		CheckInterval:    l.CheckInterval,
	}

	if l.Modules == nil {
		dto.Modules = []string{}
	}
	if l.Apps == nil {
		dto.Apps = []string{}
	}

	if l.Description != nil {
		dto.Description = *l.Description
	}
	if l.InstanceURL != nil {
		dto.InstanceURL = *l.InstanceURL
	}
	if l.InstanceName != nil {
		dto.InstanceName = *l.InstanceName
	}
	if l.ProvisionAPIKey != nil {
		dto.ProvisionAPIKey = *l.ProvisionAPIKey
	}

	if l.ExpiresAt != nil {
		s := l.ExpiresAt.UTC().Format(time.RFC3339)
		dto.ExpiresAt = &s
	}
	if l.LastPullAt != nil {
		s := l.LastPullAt.UTC().Format(time.RFC3339)
		dto.LastPullAt = &s
	}

	if company, err := h.companySvc.GetByID(r.Context(), l.CompanyID); err == nil {
		dto.CompanyName = company.Name
	}
	if project, err := h.projectSvc.GetByID(r.Context(), l.ProjectID); err == nil {
		dto.ProjectName = project.Name
	}
	if product, err := h.productSvc.GetByID(r.Context(), l.ProductID); err == nil {
		dto.ProductName = product.Name
	}

	return dto
}

// GetProvisionKey mengembalikan provision API key untuk sebuah license.
// Hanya superuser yang bisa akses endpoint ini.
func (h *LicenseHandler) GetProvisionKey(w http.ResponseWriter, r *http.Request) {
	// Check superuser role
	user, ok := appmiddleware.UserFromContext(r.Context())
	if !ok || user.Role != "superuser" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Only superuser can access provision key")
		return
	}

	licenseIDStr := chi.URLParam(r, "id")
	if licenseIDStr == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "License ID is required")
		return
	}

	licenseID, err := uuid.Parse(licenseIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid license ID format")
		return
	}

	license, err := h.licenseSvc.GetByID(r.Context(), licenseID)
	if err != nil {
		if errors.Is(err, domain.ErrLicenseNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "License not found")
			return
		}
		h.logger.Error("GetProvisionKey: GetByID failed", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response := struct {
		LicenseID           string `json:"license_id"`
		LicenseKey          string `json:"license_key"`
		ProvisionAPIKey     string `json:"provision_api_key"`
		GeneratedAt         string `json:"generated_at,omitempty"`
		NextRotationIn      string `json:"next_rotation_in,omitempty"`
	}{
		LicenseID:       license.ID.String(),
		LicenseKey:      license.LicenseKey,
		ProvisionAPIKey: "",
	}

	// Expose provision key hanya jika ada
	if license.ProvisionAPIKey != nil {
		response.ProvisionAPIKey = *license.ProvisionAPIKey
	}

	// Tampilkan kapan key di-generate
	if license.ProvisionAPIKeyGeneratedAt != nil {
		response.GeneratedAt = license.ProvisionAPIKeyGeneratedAt.UTC().Format(time.RFC3339)
		// Hitung next rotation (30 menit dari generated_at)
		nextRotation := license.ProvisionAPIKeyGeneratedAt.Add(30 * time.Minute)
		if nextRotation.After(time.Now()) {
			response.NextRotationIn = nextRotation.Sub(time.Now()).String()
		}
	}

	writeJSON(w, http.StatusOK, response)
}
