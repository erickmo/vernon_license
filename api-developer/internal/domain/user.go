// Package domain mendefinisikan domain models dan repository interfaces Vernon License.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User merepresentasikan akun internal Vernon App.
type User struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	// Role adalah salah satu dari: "superuser" | "project_owner" | "sales"
	Role      string    `db:"role"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UserRepository mendefinisikan operasi persistence untuk entitas User.
type UserRepository interface {
	// FindByID mencari user berdasarkan UUID.
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail mencari user berdasarkan email.
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindAll mengembalikan semua user.
	FindAll(ctx context.Context) ([]*User, error)

	// Create menyimpan user baru ke database.
	Create(ctx context.Context, user *User) error

	// Update memperbarui data user.
	Update(ctx context.Context, user *User) error

	// SetActive mengaktifkan atau menonaktifkan user.
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}
