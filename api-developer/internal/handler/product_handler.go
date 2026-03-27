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

// ProductHandler menangani HTTP requests untuk Products.
type ProductHandler struct {
	svc    *service.ProductService
	logger *zap.Logger
}

// NewProductHandler membuat instance ProductHandler baru.
func NewProductHandler(svc *service.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{svc: svc, logger: logger}
}

// productDTO adalah representasi Product yang dikembalikan ke client.
type productDTO struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Description      *string         `json:"description"`
	AvailableModules json.RawMessage `json:"available_modules"`
	AvailableApps    json.RawMessage `json:"available_apps"`
	AvailablePlans   []string        `json:"available_plans"`
	BasePricing      json.RawMessage `json:"base_pricing"`
	IsActive         bool            `json:"is_active"`
}

// createProductRequest adalah body JSON untuk POST /api/internal/products.
type createProductRequest struct {
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Description      *string         `json:"description"`
	AvailableModules json.RawMessage `json:"available_modules"`
	AvailableApps    json.RawMessage `json:"available_apps"`
	AvailablePlans   []string        `json:"available_plans"`
	BasePricing      json.RawMessage `json:"base_pricing"`
	IsActive         bool            `json:"is_active"`
}

// updateProductRequest adalah body JSON untuk PUT /api/internal/products/{id}.
type updateProductRequest struct {
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Description      *string         `json:"description"`
	AvailableModules json.RawMessage `json:"available_modules"`
	AvailableApps    json.RawMessage `json:"available_apps"`
	AvailablePlans   []string        `json:"available_plans"`
	BasePricing      json.RawMessage `json:"base_pricing"`
	IsActive         bool            `json:"is_active"`
}

// toProductDTO mengkonversi domain.Product ke productDTO.
func toProductDTO(p *domain.Product) productDTO {
	return productDTO{
		ID:               p.ID.String(),
		Name:             p.Name,
		Slug:             p.Slug,
		Description:      p.Description,
		AvailableModules: p.AvailableModules,
		AvailableApps:    p.AvailableApps,
		AvailablePlans:   p.AvailablePlans,
		BasePricing:      p.BasePricing,
		IsActive:         p.IsActive,
	}
}

// List menangani GET /api/internal/products.
// Superuser mendapatkan semua (termasuk inactive), role lain hanya aktif.
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	includeInactive := ok && claims.Role == "superuser"

	products, err := h.svc.List(r.Context(), includeInactive)
	if err != nil {
		h.logger.Error("ProductHandler.List", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dtos := make([]productDTO, 0, len(products))
	for _, p := range products {
		dtos = append(dtos, toProductDTO(p))
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": dtos})
}

// GetByID menangani GET /api/internal/products/{id}.
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	product, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			writeError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product tidak ditemukan")
			return
		}
		h.logger.Error("ProductHandler.GetByID", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toProductDTO(product))
}

// Create menangani POST /api/internal/products.
// Hanya superuser yang dapat membuat product.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	if claims.Role != "superuser" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Hanya superuser yang dapat membuat product")
		return
	}

	actorID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	var req createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Name == "" || req.Slug == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "name dan slug wajib diisi")
		return
	}

	product, err := h.svc.Create(r.Context(), service.CreateProductRequest{
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		AvailableModules: req.AvailableModules,
		AvailableApps:    req.AvailableApps,
		AvailablePlans:   req.AvailablePlans,
		BasePricing:      req.BasePricing,
		IsActive:         req.IsActive,
	}, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrProductSlugExists) {
			writeError(w, http.StatusConflict, "PRODUCT_SLUG_EXISTS", "Slug sudah digunakan")
			return
		}
		if errors.Is(err, domain.ErrValidationFailed) {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Data tidak valid")
			return
		}
		h.logger.Error("ProductHandler.Create", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, toProductDTO(product))
}

// Update menangani PUT /api/internal/products/{id}.
// Hanya superuser yang dapat mengubah product.
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	if claims.Role != "superuser" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Hanya superuser yang dapat mengubah product")
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

	var req updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid request body")
		return
	}

	if req.Name == "" || req.Slug == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "name dan slug wajib diisi")
		return
	}

	product, err := h.svc.Update(r.Context(), id, service.UpdateProductRequest{
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		AvailableModules: req.AvailableModules,
		AvailableApps:    req.AvailableApps,
		AvailablePlans:   req.AvailablePlans,
		BasePricing:      req.BasePricing,
		IsActive:         req.IsActive,
	}, actorID, claims.Name)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			writeError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product tidak ditemukan")
			return
		}
		if errors.Is(err, domain.ErrProductSlugExists) {
			writeError(w, http.StatusConflict, "PRODUCT_SLUG_EXISTS", "Slug sudah digunakan")
			return
		}
		if errors.Is(err, domain.ErrProductHasLicense) {
			writeError(w, http.StatusConflict, "PRODUCT_HAS_LICENSE", "Name dan slug tidak dapat diubah karena produk sudah memiliki license")
			return
		}
		if errors.Is(err, domain.ErrValidationFailed) {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Data tidak valid")
			return
		}
		h.logger.Error("ProductHandler.Update", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toProductDTO(product))
}

// Delete menangani DELETE /api/internal/products/{id}.
// Hanya superuser yang dapat menghapus product.
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	if claims.Role != "superuser" {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "Hanya superuser yang dapat menghapus product")
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
		if errors.Is(err, domain.ErrProductNotFound) {
			writeError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product tidak ditemukan")
			return
		}
		h.logger.Error("ProductHandler.Delete", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Product berhasil dihapus"})
}
