package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/flashlab/vernon-license/internal/domain"
)

// CompanyRepo adalah implementasi domain.CompanyRepository berbasis PostgreSQL.
type CompanyRepo struct {
	db *sqlx.DB
}

// NewCompanyRepo membuat instance CompanyRepo baru.
func NewCompanyRepo(db *sqlx.DB) *CompanyRepo {
	return &CompanyRepo{db: db}
}

// FindByID mencari company berdasarkan UUID (soft-delete aware).
func (r *CompanyRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	var c domain.Company
	const q = `
		SELECT id, name, email, phone, address, pic_name, pic_email, pic_phone,
		       notes, created_by, created_at, updated_at, deleted_at
		FROM companies
		WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &c, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrCompanyNotFound
		}
		return nil, fmt.Errorf("CompanyRepo.FindByID: %w", err)
	}
	return &c, nil
}

// FindAll mengembalikan semua company yang belum dihapus.
func (r *CompanyRepo) FindAll(ctx context.Context) ([]*domain.Company, error) {
	const q = `
		SELECT id, name, email, phone, address, pic_name, pic_email, pic_phone,
		       notes, created_by, created_at, updated_at, deleted_at
		FROM companies
		WHERE deleted_at IS NULL
		ORDER BY name ASC`
	var companies []*domain.Company
	if err := r.db.SelectContext(ctx, &companies, q); err != nil {
		return nil, fmt.Errorf("CompanyRepo.FindAll: %w", err)
	}
	return companies, nil
}

// Create menyimpan company baru ke database.
func (r *CompanyRepo) Create(ctx context.Context, c *domain.Company) error {
	const q = `
		INSERT INTO companies
		    (id, name, email, phone, address, pic_name, pic_email, pic_phone,
		     notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING created_at, updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		c.ID, c.Name, c.Email, c.Phone, c.Address,
		c.PICName, c.PICEmail, c.PICPhone,
		c.Notes, c.CreatedBy,
	).Scan(&c.CreatedAt, &c.UpdatedAt); err != nil {
		return fmt.Errorf("CompanyRepo.Create: %w", err)
	}
	return nil
}

// Update memperbarui data company.
func (r *CompanyRepo) Update(ctx context.Context, c *domain.Company) error {
	const q = `
		UPDATE companies
		SET name = $1, email = $2, phone = $3, address = $4,
		    pic_name = $5, pic_email = $6, pic_phone = $7,
		    notes = $8, updated_at = NOW()
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		c.Name, c.Email, c.Phone, c.Address,
		c.PICName, c.PICEmail, c.PICPhone,
		c.Notes, c.ID,
	).Scan(&c.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrCompanyNotFound
		}
		return fmt.Errorf("CompanyRepo.Update: %w", err)
	}
	return nil
}

// Delete melakukan soft delete pada company.
func (r *CompanyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE companies SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("CompanyRepo.Delete: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("CompanyRepo.Delete: rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrCompanyNotFound
	}
	return nil
}
