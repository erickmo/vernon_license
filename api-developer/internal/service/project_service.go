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

// CreateProjectRequest berisi data yang dibutuhkan untuk membuat project baru.
type CreateProjectRequest struct {
	CompanyID   uuid.UUID `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}

// UpdateProjectRequest berisi data yang dapat diperbarui pada project.
type UpdateProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	// Status adalah salah satu dari: "active" | "completed" | "cancelled"
	Status string `json:"status"`
}

// ProjectService mengelola business logic untuk entitas Project.
type ProjectService struct {
	repo      domain.ProjectRepository
	auditRepo domain.AuditLogRepository
	logger    *zap.Logger
}

// NewProjectService membuat instance ProjectService baru.
func NewProjectService(repo domain.ProjectRepository, audit domain.AuditLogRepository, logger *zap.Logger) *ProjectService {
	return &ProjectService{
		repo:      repo,
		auditRepo: audit,
		logger:    logger,
	}
}

// ListByCompany mengembalikan semua project milik sebuah company. Dapat diakses oleh semua role.
func (s *ProjectService) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*domain.Project, error) {
	projects, err := s.repo.FindByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("ProjectService.ListByCompany: %w", err)
	}
	return projects, nil
}

// GetByID mengambil satu project berdasarkan UUID.
func (s *ProjectService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProjectService.GetByID: %w", err)
	}
	return project, nil
}

// Create membuat project baru dengan status awal "active". Dapat dilakukan oleh semua role.
// Audit log dibuat dengan action "project_created" setelah operasi berhasil.
func (s *ProjectService) Create(ctx context.Context, req CreateProjectRequest, actorID uuid.UUID, actorName string) (*domain.Project, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("ProjectService.Create: %w", domain.ErrValidationFailed)
	}
	if req.CompanyID == uuid.Nil {
		return nil, fmt.Errorf("ProjectService.Create: company_id is required: %w", domain.ErrValidationFailed)
	}

	project := &domain.Project{
		ID:          uuid.New(),
		CompanyID:   req.CompanyID,
		Name:        req.Name,
		Description: req.Description,
		Status:      "active",
		CreatedBy:   actorID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("ProjectService.Create: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{
		"company_id":  project.CompanyID.String(),
		"name":        project.Name,
		"description": project.Description,
		"status":      project.Status,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "project", project.ID, "project_created", actorID, actorName, changes)

	return project, nil
}

// Update memperbarui data project yang ada.
// Audit log dibuat dengan action "project_updated" setelah operasi berhasil.
func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, req UpdateProjectRequest, actorID uuid.UUID, actorName string) (*domain.Project, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("ProjectService.Update: %w", domain.ErrValidationFailed)
	}

	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProjectService.Update: %w", err)
	}

	project.Name = req.Name
	project.Description = req.Description
	if req.Status != "" {
		project.Status = req.Status
	}

	if err := s.repo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("ProjectService.Update: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{
		"name":        project.Name,
		"description": project.Description,
		"status":      project.Status,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "project", project.ID, "project_updated", actorID, actorName, changes)

	return project, nil
}

// Delete melakukan soft-delete pada project.
// Audit log dibuat dengan action "project_deleted" setelah operasi berhasil.
func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorName string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("ProjectService.Delete: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{"id": id.String()})
	LogAudit(ctx, s.auditRepo, s.logger, "project", id, "project_deleted", actorID, actorName, changes)

	return nil
}
