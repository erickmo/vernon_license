package domain

import "context"

// OTPRepository adalah interface untuk manajemen OTP (One-Time Password).
type OTPRepository interface {
	// IsActive memverifikasi bahwa code ada di tabel otp dan belum expired.
	// OTP dapat digunakan berkali-kali selama masih dalam periode aktif.
	IsActive(ctx context.Context, code string) error

	// GetActive mengembalikan kode OTP yang sedang aktif (belum expired).
	// Mengembalikan error jika tidak ada OTP aktif.
	GetActive(ctx context.Context) (string, error)
}
