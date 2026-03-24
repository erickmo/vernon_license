package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/flashlab/vernon-license/internal/domain"
)

// ProjectRepo adalah implementasi domain.ProjectRepository berbasis PostgreSQL.
type ProjectRepo struct {
	db *sqlx.DB
}

// NewProjectRepo membuat instance ProjectRepo baru.
func NewProjectRepo(db *sqlx.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

// FindByID mencari project berdasarkan UUID (soft-delete aware).
func (r *ProjectRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	var p domain.Project
	const q = `
		SELECT id, company_id, name, description, status, created_by,
		       created_at, updated_at, deleted_at
		FROM projects
		WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProjectNotFound
		}
		return nil, fmt.Errorf("ProjectRepo.FindByID: %w", err)
	}
	return &p, nil
}

// FindByCompany mengembalikan semua project yang belum dihapus untuk sebuah company.
func (r *ProjectRepo) FindByCompany(ctx context.Context, companyID uuid.UUID) ([]*domain.Project, error) {
	const q = `
		SELECT id, company_id, name, description, status, created_by,
		       created_at, updated_at, deleted_at
		FROM projects
		WHERE company_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`
	var projects []*domain.Project
	if err := r.db.SelectContext(ctx, &projects, q, companyID); err != nil {
		return nil, fmt.Errorf("ProjectRepo.FindByCompany: %w", err)
	}
	return projects, nil
}

// Create menyimpan project baru ke database.
func (r *ProjectRepo) Create(ctx context.Context, p *domain.Project) error {
	const q = `
		INSERT INTO projects (id, company_id, name, description, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		p.ID, p.CompanyID, p.Name, p.Description, p.Status, p.CreatedBy,
	).Scan(&p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("ProjectRepo.Create: %w", err)
	}
	return nil
}

// Update memperbarui data project.
func (r *ProjectRepo) Update(ctx context.Context, p *domain.Project) error {
	const q = `
		UPDATE projects
		SET name = $1, description = $2, status = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		p.Name, p.Description, p.Status, p.ID,
	).Scan(&p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrProjectNotFound
		}
		return fmt.Errorf("ProjectRepo.Update: %w", err)
	}
	return nil
}

// Delete melakukan soft delete pada project.
func (r *ProjectRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE projects SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("ProjectRepo.Delete: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ProjectRepo.Delete: rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrProjectNotFound
	}
	return nil
}
