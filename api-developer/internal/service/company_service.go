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

// CreateCompanyRequest berisi data yang dibutuhkan untuk membuat company baru.
type CreateCompanyRequest struct {
	Name     string  `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	PICName  *string `json:"pic_name"`
	PICEmail *string `json:"pic_email"`
	PICPhone *string `json:"pic_phone"`
	Notes    *string `json:"notes"`
}

// UpdateCompanyRequest berisi data yang dapat diperbarui pada company.
type UpdateCompanyRequest struct {
	Name     string  `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	PICName  *string `json:"pic_name"`
	PICEmail *string `json:"pic_email"`
	PICPhone *string `json:"pic_phone"`
	Notes    *string `json:"notes"`
}

// CompanyService mengelola business logic untuk entitas Company.
type CompanyService struct {
	repo      domain.CompanyRepository
	auditRepo domain.AuditLogRepository
	logger    *zap.Logger
}

// NewCompanyService membuat instance CompanyService baru.
func NewCompanyService(repo domain.CompanyRepository, audit domain.AuditLogRepository, logger *zap.Logger) *CompanyService {
	return &CompanyService{
		repo:      repo,
		auditRepo: audit,
		logger:    logger,
	}
}

// List mengembalikan semua company yang belum dihapus. Dapat diakses oleh semua role.
func (s *CompanyService) List(ctx context.Context) ([]*domain.Company, error) {
	companies, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.List: %w", err)
	}
	return companies, nil
}

// GetByID mengambil satu company berdasarkan UUID.
func (s *CompanyService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	company, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.GetByID: %w", err)
	}
	return company, nil
}

// Create membuat company baru. Dapat dilakukan oleh semua role.
// Audit log dibuat dengan action "company_created" setelah operasi berhasil.
func (s *CompanyService) Create(ctx context.Context, req CreateCompanyRequest, actorID uuid.UUID, actorName string) (*domain.Company, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("CompanyService.Create: %w", domain.ErrValidationFailed)
	}

	company := &domain.Company{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		PICName:   req.PICName,
		PICEmail:  req.PICEmail,
		PICPhone:  req.PICPhone,
		Notes:     req.Notes,
		CreatedBy: actorID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, company); err != nil {
		return nil, fmt.Errorf("CompanyService.Create: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{
		"name":      company.Name,
		"email":     company.Email,
		"phone":     company.Phone,
		"address":   company.Address,
		"pic_name":  company.PICName,
		"pic_email": company.PICEmail,
		"pic_phone": company.PICPhone,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "company", company.ID, "company_created", actorID, actorName, changes)

	return company, nil
}

// Update memperbarui data company yang ada.
// Audit log dibuat dengan action "company_updated" setelah operasi berhasil.
func (s *CompanyService) Update(ctx context.Context, id uuid.UUID, req UpdateCompanyRequest, actorID uuid.UUID, actorName string) (*domain.Company, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("CompanyService.Update: %w", domain.ErrValidationFailed)
	}

	company, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CompanyService.Update: %w", err)
	}

	company.Name = req.Name
	company.Email = req.Email
	company.Phone = req.Phone
	company.Address = req.Address
	company.PICName = req.PICName
	company.PICEmail = req.PICEmail
	company.PICPhone = req.PICPhone
	company.Notes = req.Notes

	if err := s.repo.Update(ctx, company); err != nil {
		return nil, fmt.Errorf("CompanyService.Update: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{
		"name":      company.Name,
		"email":     company.Email,
		"phone":     company.Phone,
		"address":   company.Address,
		"pic_name":  company.PICName,
		"pic_email": company.PICEmail,
		"pic_phone": company.PICPhone,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "company", company.ID, "company_updated", actorID, actorName, changes)

	return company, nil
}

// Delete melakukan soft-delete pada company.
// Audit log dibuat dengan action "company_deleted" setelah operasi berhasil.
func (s *CompanyService) Delete(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorName string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("CompanyService.Delete: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{"id": id.String()})
	LogAudit(ctx, s.auditRepo, s.logger, "company", id, "company_deleted", actorID, actorName, changes)

	return nil
}
