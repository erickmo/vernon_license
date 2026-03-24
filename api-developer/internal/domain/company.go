package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Company merepresentasikan perusahaan klien Vernon.
type Company struct {
	ID        uuid.UUID  `db:"id"`
	Name      string     `db:"name"`
	Email     *string    `db:"email"`
	Phone     *string    `db:"phone"`
	Address   *string    `db:"address"`
	PICName   *string    `db:"pic_name"`
	PICEmail  *string    `db:"pic_email"`
	PICPhone  *string    `db:"pic_phone"`
	Notes     *string    `db:"notes"`
	CreatedBy uuid.UUID  `db:"created_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// CompanyRepository mendefinisikan operasi persistence untuk entitas Company.
type CompanyRepository interface {
	// FindByID mencari company berdasarkan UUID (soft-delete aware).
	FindByID(ctx context.Context, id uuid.UUID) (*Company, error)

	// FindAll mengembalikan semua company yang belum dihapus.
	FindAll(ctx context.Context) ([]*Company, error)

	// Create menyimpan company baru ke database.
	Create(ctx context.Context, c *Company) error

	// Update memperbarui data company.
	Update(ctx context.Context, c *Company) error

	// Delete melakukan soft delete pada company.
	Delete(ctx context.Context, id uuid.UUID) error
}
