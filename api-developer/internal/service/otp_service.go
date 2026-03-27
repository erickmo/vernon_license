//go:build !wasm

// Package service menyediakan business logic untuk OTP.
package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// OTPService menangani One-Time Password generation dan management.
type OTPService struct {
	db  *sqlx.DB
	log *zap.Logger
}

// NewOTPService membuat instance OTPService baru.
func NewOTPService(db *sqlx.DB, log *zap.Logger) *OTPService {
	return &OTPService{
		db:  db,
		log: log,
	}
}

// GenerateOTP menghasilkan OTP baru (16 karakter hex).
// OTP berlaku selama 30 menit dan dapat digunakan berkali-kali dalam periode tersebut.
func (s *OTPService) GenerateOTP(ctx context.Context) (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("GenerateOTP: %w", err)
	}
	code := hex.EncodeToString(b)

	expiresAt := time.Now().Add(30 * time.Minute)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO otp (id, code, created_at, expires_at)
		VALUES ($1, $2, NOW(), $3)
	`, uuid.New(), code, expiresAt)
	if err != nil {
		return "", fmt.Errorf("GenerateOTP: %w", err)
	}

	s.log.Debug("OTP generated", zap.String("code", code))
	return code, nil
}

// GetCurrentOTP mengambil OTP yang aktif sekarang.
// Jika belum ada atau expired, generate OTP baru.
func (s *OTPService) GetCurrentOTP(ctx context.Context) (string, time.Time, error) {
	const q = `
		SELECT code, expires_at
		FROM otp
		WHERE expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	var code string
	var expiresAt time.Time
	err := s.db.QueryRowContext(ctx, q).Scan(&code, &expiresAt)
	if err == nil {
		return code, expiresAt, nil
	}

	// Tidak ada OTP aktif, generate yang baru
	newCode, err := s.GenerateOTP(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("GetCurrentOTP: %w", err)
	}

	var newExpiresAt time.Time
	err = s.db.QueryRowContext(ctx, `
		SELECT expires_at FROM otp WHERE code = $1
	`, newCode).Scan(&newExpiresAt)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("GetCurrentOTP: scan expiry: %w", err)
	}

	return newCode, newExpiresAt, nil
}

// CleanupExpired menghapus OTP yang sudah expired.
// Dipanggil oleh scheduler secara berkala.
func (s *OTPService) CleanupExpired(ctx context.Context) error {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM otp
		WHERE expires_at < NOW() - INTERVAL '1 hour'
	`)
	if err != nil {
		return fmt.Errorf("CleanupExpired: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("CleanupExpired: rows affected: %w", err)
	}

	s.log.Debug("Expired OTPs cleaned up", zap.Int64("count", rows))
	return nil
}
