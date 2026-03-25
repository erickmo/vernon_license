package domain

import "context"

// OTPRepository adalah interface untuk manajemen OTP (One-Time Password).
type OTPRepository interface {
	// IsActive memverifikasi bahwa code ada di tabel otp dan belum expired.
	// OTP dapat digunakan berkali-kali selama masih dalam periode aktif.
	IsActive(ctx context.Context, code string) error
}
