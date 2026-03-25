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

// CompanyHandler menangani HTTP requests untuk Companies.
type CompanyHandler struct {
	svc    *service.CompanyService
	logger *zap.Logger
}

// NewCompanyHandler membuat instance CompanyHandler baru.
func NewCompanyHandler(svc *service.CompanyService, logger *zap.Logger) *CompanyHandler {
	return &CompanyHandler{svc: svc, logger: logger}
}

// companyDTO adalah representasi Company yang dikembalikan ke client.
type companyDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	PICName   *string `json:"pic_name"`
	PICEmail  *string `json:"pic_email"`
	PICPhone  *string `json:"pic_phone"`
	Notes     *string `json:"notes"`
	CreatedAt string  `json:"created_at"`
}

// createCompanyRequest adalah body JSON untuk POST /api/internal/companies.
type createCompanyRequest struct {
	Name     string  `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	PICName  *string `json:"pic_name"`
	PICEmail *string `json:"pic_email"`
	PICPhone *string `json:"pic_phone"`
	Notes    *string `json:"notes"`
}

// updateCompanyRequest adalah body JSON untuk PUT /api/internal/companies/{id}.
type updateCompanyRequest struct {
	Name     string  `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	PICName  *string `json:"pic_name"`
	PICEmail *string `json:"pic_email"`
	PICPhone *string `json:"pic_phone"`
	Notes    *string `json:"notes"`
}

// toCompanyDTO mengkonversi domain.Company ke companyDTO.
func toCompanyDTO(c *domain.Company) companyDTO {
	return companyDTO{
		ID:        c.ID.String(),
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		PICName:   c.PICName,
		PICEmail:  c.PICEmail,
		PICPhone:  c.PICPhone,
		Notes:     c.Notes,
		CreatedAt: c.CreatedAt.Format("2006-01-02"),
	}
}

// List menangani GET /api/internal/companies.
// Mengembalikan semua companies yang belum dihapus.
func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	companies, err := h.svc.List(r.Context())
	if err != nil {
		h.logger.Error("CompanyHandler.List", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dtos := make([]companyDTO, 0, len(companies))
	for _, c := range companies {
		dtos = append(dtos, toCompanyDTO(c))
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": dtos})
}

// GetByID menangani GET /api/internal/companies/{id}.
// Mengembalikan satu company berdasarkan UUID.
func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	company, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			writeError(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company tidak ditemukan")
			return
		}
		h.logger.Error("CompanyHandler.GetByID", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toCompanyDTO(company))
}

// Create menangani POST /api/internal/companies.
// Body: {name, email, phone, address, pic_name, pic_email, pic_phone, notes}
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	var req createCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "name wajib diisi")
		return
	}

	company, err := h.svc.Create(r.Context(), service.CreateCompanyRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Address:  req.Address,
		PICName:  req.PICName,
		PICEmail: req.PICEmail,
		PICPhone: req.PICPhone,
		Notes:    req.Notes,
	}, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrValidationFailed) {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Data tidak valid")
			return
		}
		h.logger.Error("CompanyHandler.Create", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, toCompanyDTO(company))
}

// Update menangani PUT /api/internal/companies/{id}.
// Body: {name, email, phone, address, pic_name, pic_email, pic_phone, notes}
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
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

	var req updateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "name wajib diisi")
		return
	}

	company, err := h.svc.Update(r.Context(), id, service.UpdateCompanyRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Address:  req.Address,
		PICName:  req.PICName,
		PICEmail: req.PICEmail,
		PICPhone: req.PICPhone,
		Notes:    req.Notes,
	}, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			writeError(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrValidationFailed) {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Data tidak valid")
			return
		}
		h.logger.Error("CompanyHandler.Update", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toCompanyDTO(company))
}

// Delete menangani DELETE /api/internal/companies/{id}.
// Melakukan soft-delete pada company.
func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
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

	if err := h.svc.Delete(r.Context(), id, actorID, claims.Name); err != nil {
		if errors.Is(err, domain.ErrCompanyNotFound) {
			writeError(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company tidak ditemukan")
			return
		}
		h.logger.Error("CompanyHandler.Delete", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Company berhasil dihapus"})
}
