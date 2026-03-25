package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Product merepresentasikan produk software yang dapat dilisensikan.
type Product struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	Description *string   `db:"description"`
	// AvailableModules adalah list modul yang tersedia, format: [{key, name, description}]
	AvailableModules json.RawMessage `db:"available_modules"`
	// AvailableApps adalah list aplikasi yang tersedia.
	AvailableApps  json.RawMessage `db:"available_apps"`
	AvailablePlans []string        `db:"available_plans"`
	BasePricing    json.RawMessage `db:"base_pricing"`
	IsActive       bool            `db:"is_active"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
	DeletedAt      *time.Time      `db:"deleted_at"`
}

// ProductRepository mendefinisikan operasi persistence untuk entitas Product.
type ProductRepository interface {
	// FindByID mencari product berdasarkan UUID (soft-delete aware).
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)

	// FindBySlug mencari product berdasarkan slug.
	FindBySlug(ctx context.Context, slug string) (*Product, error)

	// FindAll mengembalikan semua product. Jika includeInactive false, hanya yang aktif.
	FindAll(ctx context.Context, includeInactive bool) ([]*Product, error)

	// Create menyimpan product baru ke database.
	Create(ctx context.Context, p *Product) error

	// Update memperbarui data product.
	Update(ctx context.Context, p *Product) error

	// Delete melakukan soft delete pada product.
	Delete(ctx context.Context, id uuid.UUID) error
}
