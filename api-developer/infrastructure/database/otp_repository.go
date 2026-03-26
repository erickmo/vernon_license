//go:build !wasm

package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// OTPRepository implementasi domain.OTPRepository menggunakan PostgreSQL.
type OTPRepository struct {
	db *sqlx.DB
}

// NewOTPRepository membuat instance OTPRepository baru.
func NewOTPRepository(db *sqlx.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

// IsActive memverifikasi bahwa code ada di tabel otp dan belum expired.
// OTP dapat digunakan berkali-kali selama masih dalam periode aktif.
func (r *OTPRepository) IsActive(ctx context.Context, code string) error {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM otp
		WHERE code = $1 AND expires_at > NOW()
	`, code)
	if err != nil {
		return fmt.Errorf("IsActive: %w", err)
	}
	if count == 0 {
		return errors.New("OTP not found or expired")
	}
	return nil
}

// GetActive mengembalikan kode OTP yang sedang aktif (belum expired).
func (r *OTPRepository) GetActive(ctx context.Context) (string, error) {
	var code string
	err := r.db.GetContext(ctx, &code, `
		SELECT code FROM otp
		WHERE expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`)
	if err != nil {
		return "", fmt.Errorf("GetActive: %w", err)
	}
	return code, nil
}
