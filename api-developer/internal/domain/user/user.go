package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	RoleDeveloperSales = "developer_sales"
	RoleSuperuser      = "superuser"

	StatusActive   = "active"
	StatusInactive = "inactive"
)

var (
	ErrNotFound            = errors.New("user tidak ditemukan")
	ErrInvalidCredentials  = errors.New("email atau password salah")
	ErrInactiveAccount     = errors.New("akun tidak aktif")
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewUser(name, email, passwordHash, role string) (*User, error) {
	return &User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		Status:       StatusActive,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}, nil
}

type WriteRepository interface {
	Save(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type ReadRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CountByRole(ctx context.Context, role string) (int, error)
}
