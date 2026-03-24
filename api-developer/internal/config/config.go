// Package config mengelola konfigurasi aplikasi dari environment variables.
package config

import (
	"fmt"
	"os"
)

// Config menyimpan semua konfigurasi aplikasi Vernon License.
type Config struct {
	// DatabaseURL adalah PostgreSQL connection string.
	DatabaseURL string

	// JWTSecret adalah secret key HS256 untuk Vernon App login.
	JWTSecret string

	// Port adalah port HTTP server (default: "8081").
	Port string

	// LogLevel adalah level logging (default: "info").
	LogLevel string

	// StoragePath adalah path penyimpanan file lokal.
	StoragePath string

	// LicenseCheckInterval adalah interval cek license dari client app (default: "6h").
	LicenseCheckInterval string

	// CompanyName adalah nama perusahaan Vernon.
	CompanyName string

	// CompanyAddress adalah alamat perusahaan Vernon.
	CompanyAddress string

	// CompanyPhone adalah nomor telepon perusahaan Vernon.
	CompanyPhone string

	// CompanyEmail adalah email perusahaan Vernon.
	CompanyEmail string

	// CompanyLogoPath adalah path ke file logo perusahaan.
	CompanyLogoPath string
}

// Load membaca konfigurasi dari environment variables.
// Mengembalikan error jika konfigurasi wajib tidak tersedia.
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("Load: DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("Load: JWT_SECRET is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}

	checkInterval := os.Getenv("LICENSE_CHECK_INTERVAL")
	if checkInterval == "" {
		checkInterval = "6h"
	}

	return &Config{
		DatabaseURL:          dbURL,
		JWTSecret:            jwtSecret,
		Port:                 port,
		LogLevel:             logLevel,
		StoragePath:          storagePath,
		LicenseCheckInterval: checkInterval,
		CompanyName:          os.Getenv("COMPANY_NAME"),
		CompanyAddress:       os.Getenv("COMPANY_ADDRESS"),
		CompanyPhone:         os.Getenv("COMPANY_PHONE"),
		CompanyEmail:         os.Getenv("COMPANY_EMAIL"),
		CompanyLogoPath:      os.Getenv("COMPANY_LOGO_PATH"),
	}, nil
}
