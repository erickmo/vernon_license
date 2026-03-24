package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/flashlab/vernon-license/internal/domain"
)

// UserRepo adalah implementasi domain.UserRepository berbasis PostgreSQL.
type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo membuat instance UserRepo baru.
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// FindByID mencari user berdasarkan UUID.
func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var u domain.User
	const q = `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1`
	if err := r.db.GetContext(ctx, &u, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("UserRepo.FindByID: %w", err)
	}
	return &u, nil
}

// FindByEmail mencari user berdasarkan email.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	const q = `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1`
	if err := r.db.GetContext(ctx, &u, q, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("UserRepo.FindByEmail: %w", err)
	}
	return &u, nil
}

// FindAll mengembalikan semua user.
func (r *UserRepo) FindAll(ctx context.Context) ([]*domain.User, error) {
	const q = `
		SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`
	var users []*domain.User
	if err := r.db.SelectContext(ctx, &users, q); err != nil {
		return nil, fmt.Errorf("UserRepo.FindAll: %w", err)
	}
	return users, nil
}

// Create menyimpan user baru ke database.
func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	const q = `
		INSERT INTO users (id, name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		u.ID, u.Name, u.Email, u.PasswordHash, u.Role, u.IsActive,
	).Scan(&u.CreatedAt, &u.UpdatedAt); err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}
	return nil
}

// Update memperbarui data user.
func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	const q = `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, role = $4, is_active = $5,
		    updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		u.Name, u.Email, u.PasswordHash, u.Role, u.IsActive, u.ID,
	).Scan(&u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("UserRepo.Update: %w", err)
	}
	return nil
}

// SetActive mengaktifkan atau menonaktifkan user.
func (r *UserRepo) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	const q = `UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.ExecContext(ctx, q, active, id)
	if err != nil {
		return fmt.Errorf("UserRepo.SetActive: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UserRepo.SetActive: rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
