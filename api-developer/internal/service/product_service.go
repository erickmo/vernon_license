package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
)

// CreateProductRequest berisi data yang dibutuhkan untuk membuat product baru.
// Hanya dapat dipanggil oleh superuser — caller wajib memvalidasi role sebelum memanggil service ini.
type CreateProductRequest struct {
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Description      *string         `json:"description"`
	AvailableModules json.RawMessage `json:"available_modules"`
	AvailableApps    json.RawMessage `json:"available_apps"`
	AvailablePlans   []string        `json:"available_plans"`
	BasePricing      json.RawMessage `json:"base_pricing"`
	IsActive         bool            `json:"is_active"`
}

// UpdateProductRequest berisi data yang dapat diperbarui pada product.
// Hanya dapat dipanggil oleh superuser — caller wajib memvalidasi role sebelum memanggil service ini.
type UpdateProductRequest struct {
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Description      *string         `json:"description"`
	AvailableModules json.RawMessage `json:"available_modules"`
	AvailableApps    json.RawMessage `json:"available_apps"`
	AvailablePlans   []string        `json:"available_plans"`
	BasePricing      json.RawMessage `json:"base_pricing"`
	IsActive         bool            `json:"is_active"`
}

// ProductService mengelola business logic untuk entitas Product.
type ProductService struct {
	repo        domain.ProductRepository
	licenseRepo domain.LicenseRepository
	auditRepo   domain.AuditLogRepository
	logger      *zap.Logger
}

// NewProductService membuat instance ProductService baru.
func NewProductService(repo domain.ProductRepository, licenseRepo domain.LicenseRepository, audit domain.AuditLogRepository, logger *zap.Logger) *ProductService {
	return &ProductService{
		repo:        repo,
		licenseRepo: licenseRepo,
		auditRepo:   audit,
		logger:      logger,
	}
}

// List mengembalikan semua product. Jika includeInactive false, hanya yang aktif dikembalikan.
// Dapat diakses oleh semua role.
func (s *ProductService) List(ctx context.Context, includeInactive bool) ([]*domain.Product, error) {
	products, err := s.repo.FindAll(ctx, includeInactive)
	if err != nil {
		return nil, fmt.Errorf("ProductService.List: %w", err)
	}
	return products, nil
}

// GetByID mengambil satu product berdasarkan UUID.
func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.GetByID: %w", err)
	}
	return product, nil
}

// GetBySlug mengambil product berdasarkan slug.
func (s *ProductService) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	product, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("ProductService.GetBySlug: %w", err)
	}
	return product, nil
}

// Create membuat product baru. Hanya untuk superuser — caller wajib memvalidasi role terlebih dahulu.
// Mengembalikan ErrProductSlugExists jika slug sudah digunakan.
// Audit log dibuat dengan action "product_created" setelah operasi berhasil.
func (s *ProductService) Create(ctx context.Context, req CreateProductRequest, actorID uuid.UUID, actorName string) (*domain.Product, error) {
	if req.Name == "" || req.Slug == "" {
		return nil, fmt.Errorf("ProductService.Create: name and slug are required: %w", domain.ErrValidationFailed)
	}

	existing, err := s.repo.FindBySlug(ctx, req.Slug)
	if err != nil && err != domain.ErrProductNotFound {
		return nil, fmt.Errorf("ProductService.Create: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("ProductService.Create: %w", domain.ErrProductSlugExists)
	}

	availableModules := req.AvailableModules
	if len(availableModules) == 0 {
		availableModules = json.RawMessage("[]")
	}
	availableApps := req.AvailableApps
	if len(availableApps) == 0 {
		availableApps = json.RawMessage("[]")
	}
	basePricing := req.BasePricing
	if len(basePricing) == 0 {
		basePricing = json.RawMessage("{}")
	}
	availablePlans := req.AvailablePlans
	if availablePlans == nil {
		availablePlans = []string{}
	}

	product := &domain.Product{
		ID:               uuid.New(),
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		AvailableModules: availableModules,
		AvailableApps:    availableApps,
		AvailablePlans:   availablePlans,
		BasePricing:      basePricing,
		IsActive:         req.IsActive,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("ProductService.Create: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{
		"name":        product.Name,
		"slug":        product.Slug,
		"description": product.Description,
		"is_active":   product.IsActive,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "product", product.ID, "product_created", actorID, actorName, changes)

	return product, nil
}

// Update memperbarui data product. Hanya untuk superuser — caller wajib memvalidasi role terlebih dahulu.
// Audit log dibuat dengan action "product_updated" setelah operasi berhasil.
func (s *ProductService) Update(ctx context.Context, id uuid.UUID, req UpdateProductRequest, actorID uuid.UUID, actorName string) (*domain.Product, error) {
	if req.Name == "" || req.Slug == "" {
		return nil, fmt.Errorf("ProductService.Update: name and slug are required: %w", domain.ErrValidationFailed)
	}

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProductService.Update: %w", err)
	}

	// Blokir perubahan name/slug jika sudah ada license yang menggunakan produk ini
	if product.Name != req.Name || product.Slug != req.Slug {
		hasLicense, err := s.licenseRepo.ExistsByProductID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("ProductService.Update: %w", err)
		}
		if hasLicense {
			return nil, fmt.Errorf("ProductService.Update: %w", domain.ErrProductHasLicense)
		}
	}

	// Cek slug uniqueness hanya jika slug berubah
	if product.Slug != req.Slug {
		existing, err := s.repo.FindBySlug(ctx, req.Slug)
		if err != nil && err != domain.ErrProductNotFound {
			return nil, fmt.Errorf("ProductService.Update: %w", err)
		}
		if existing != nil {
			return nil, fmt.Errorf("ProductService.Update: %w", domain.ErrProductSlugExists)
		}
	}

	product.Name = req.Name
	product.Slug = req.Slug
	product.Description = req.Description
	product.IsActive = req.IsActive

	if len(req.AvailableModules) > 0 {
		product.AvailableModules = req.AvailableModules
	}
	if len(req.AvailableApps) > 0 {
		product.AvailableApps = req.AvailableApps
	}
	if len(req.AvailablePlans) > 0 {
		product.AvailablePlans = req.AvailablePlans
	}
	if len(req.BasePricing) > 0 {
		product.BasePricing = req.BasePricing
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("ProductService.Update: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{
		"name":        product.Name,
		"slug":        product.Slug,
		"description": product.Description,
		"is_active":   product.IsActive,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "product", product.ID, "product_updated", actorID, actorName, changes)

	return product, nil
}

// Delete melakukan soft-delete pada product. Hanya untuk superuser — caller wajib memvalidasi role terlebih dahulu.
// Audit log dibuat dengan action "product_deleted" setelah operasi berhasil.
func (s *ProductService) Delete(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorName string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("ProductService.Delete: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{"id": id.String()})
	LogAudit(ctx, s.auditRepo, s.logger, "product", id, "product_deleted", actorID, actorName, changes)

	return nil
}
