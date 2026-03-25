package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Project merepresentasikan proyek yang dimiliki oleh sebuah company.
type Project struct {
	ID          uuid.UUID `db:"id"`
	CompanyID   uuid.UUID `db:"company_id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	// Status adalah salah satu dari: "active" | "completed" | "cancelled"
	Status    string     `db:"status"`
	CreatedBy uuid.UUID  `db:"created_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// ProjectRepository mendefinisikan operasi persistence untuk entitas Project.
type ProjectRepository interface {
	// FindByID mencari project berdasarkan UUID (soft-delete aware).
	FindByID(ctx context.Context, id uuid.UUID) (*Project, error)

	// FindByCompany mengembalikan semua project milik sebuah company.
	FindByCompany(ctx context.Context, companyID uuid.UUID) ([]*Project, error)

	// Create menyimpan project baru ke database.
	Create(ctx context.Context, p *Project) error

	// Update memperbarui data project.
	Update(ctx context.Context, p *Project) error

	// Delete melakukan soft delete pada project.
	Delete(ctx context.Context, id uuid.UUID) error
}
