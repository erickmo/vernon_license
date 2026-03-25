package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/flashlab/vernon-license/internal/domain"
)

// ProductRepo adalah implementasi domain.ProductRepository berbasis PostgreSQL.
type ProductRepo struct {
	db *sqlx.DB
}

// NewProductRepo membuat instance ProductRepo baru.
func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

// FindByID mencari product berdasarkan UUID (soft-delete aware).
func (r *ProductRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var p productRow
	const q = `
		SELECT id, name, slug, description, available_modules, available_apps,
		       available_plans, base_pricing, is_active, created_at, updated_at, deleted_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &p, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("ProductRepo.FindByID: %w", err)
	}
	return p.toDomain(), nil
}

// FindBySlug mencari product berdasarkan slug.
func (r *ProductRepo) FindBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var p productRow
	const q = `
		SELECT id, name, slug, description, available_modules, available_apps,
		       available_plans, base_pricing, is_active, created_at, updated_at, deleted_at
		FROM products
		WHERE slug = $1 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &p, q, slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("ProductRepo.FindBySlug: %w", err)
	}
	return p.toDomain(), nil
}

// FindAll mengembalikan semua product. Jika includeInactive false, hanya yang aktif.
func (r *ProductRepo) FindAll(ctx context.Context, includeInactive bool) ([]*domain.Product, error) {
	q := `
		SELECT id, name, slug, description, available_modules, available_apps,
		       available_plans, base_pricing, is_active, created_at, updated_at, deleted_at
		FROM products
		WHERE deleted_at IS NULL`
	if !includeInactive {
		q += ` AND is_active = TRUE`
	}
	q += ` ORDER BY name ASC`

	var rows []productRow
	if err := r.db.SelectContext(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("ProductRepo.FindAll: %w", err)
	}
	products := make([]*domain.Product, len(rows))
	for i, row := range rows {
		r2 := row
		products[i] = r2.toDomain()
	}
	return products, nil
}

// Create menyimpan product baru ke database.
func (r *ProductRepo) Create(ctx context.Context, p *domain.Product) error {
	const q = `
		INSERT INTO products
		    (id, name, slug, description, available_modules, available_apps,
		     available_plans, base_pricing, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		p.ID, p.Name, p.Slug, p.Description,
		p.AvailableModules, p.AvailableApps,
		pq.Array(p.AvailablePlans), p.BasePricing,
		p.IsActive,
	).Scan(&p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("ProductRepo.Create: %w", err)
	}
	return nil
}

// Update memperbarui data product.
func (r *ProductRepo) Update(ctx context.Context, p *domain.Product) error {
	const q = `
		UPDATE products
		SET name = $1, slug = $2, description = $3,
		    available_modules = $4, available_apps = $5,
		    available_plans = $6, base_pricing = $7,
		    is_active = $8, updated_at = NOW()
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		p.Name, p.Slug, p.Description,
		p.AvailableModules, p.AvailableApps,
		pq.Array(p.AvailablePlans), p.BasePricing,
		p.IsActive, p.ID,
	).Scan(&p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrProductNotFound
		}
		return fmt.Errorf("ProductRepo.Update: %w", err)
	}
	return nil
}

// Delete melakukan soft delete pada product.
func (r *ProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE products SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("ProductRepo.Delete: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ProductRepo.Delete: rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

// productRow adalah intermediate struct untuk scan TEXT[] dari PostgreSQL.
// Menggunakan field eksplisit (bukan embedding) untuk menghindari konflik tag db.
type productRow struct {
	ID               uuid.UUID      `db:"id"`
	Name             string         `db:"name"`
	Slug             string         `db:"slug"`
	Description      *string        `db:"description"`
	AvailableModules []byte         `db:"available_modules"`
	AvailableApps    []byte         `db:"available_apps"`
	AvailablePlans   pq.StringArray `db:"available_plans"`
	BasePricing      []byte         `db:"base_pricing"`
	IsActive         bool           `db:"is_active"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
	DeletedAt        *time.Time     `db:"deleted_at"`
}

func (pr *productRow) toDomain() *domain.Product {
	return &domain.Product{
		ID:               pr.ID,
		Name:             pr.Name,
		Slug:             pr.Slug,
		Description:      pr.Description,
		AvailableModules: pr.AvailableModules,
		AvailableApps:    pr.AvailableApps,
		AvailablePlans:   []string(pr.AvailablePlans),
		BasePricing:      pr.BasePricing,
		IsActive:         pr.IsActive,
		CreatedAt:        pr.CreatedAt,
		UpdatedAt:        pr.UpdatedAt,
		DeletedAt:        pr.DeletedAt,
	}
}
