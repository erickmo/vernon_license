//go:build !wasm

// Package service mengimplementasikan business logic untuk provisioning licenses.
package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
)

// ProvisionService menangani provisioning API key generation dan rotation.
type ProvisionService struct {
	licenses domain.LicenseRepository
	db       *sqlx.DB
	log      *zap.Logger
}

// NewProvisionService membuat instance ProvisionService baru.
func NewProvisionService(
	licenses domain.LicenseRepository,
	db *sqlx.DB,
	log *zap.Logger,
) *ProvisionService {
	return &ProvisionService{
		licenses: licenses,
		db:       db,
		log:      log,
	}
}

// GenerateProvisionKey menghasilkan provision API key acak dengan format:
// PROV-[32 hex chars]
func (s *ProvisionService) GenerateProvisionKey() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("GenerateProvisionKey: %w", err)
	}
	return "PROV-" + hex.EncodeToString(b), nil
}

// RotateForProject mengrotasi provision key untuk semua licenses dalam sebuah project.
// Grace period: previous key masih valid selama grace period (default: 1 jam).
// Dipanggil saat membuat license dari proposal yang di-approve.
func (s *ProvisionService) RotateForProject(ctx context.Context, projectID uuid.UUID) error {
	licenses, err := s.licenses.FindByProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("RotateForProject: %w", err)
	}

	for _, lic := range licenses {
		if err := s.rotateLicense(ctx, lic); err != nil {
			s.log.Error("RotateForProject: rotate license failed",
				zap.String("license_id", lic.ID.String()),
				zap.Error(err))
			// Non-fatal: lanjut ke license berikutnya
		}
	}

	return nil
}

// RotateAll mengrotasi provision key untuk SEMUA licenses yang belum diregistrasi.
// Ini adalah scheduler job yang dijalankan setiap 30 menit.
// Hanya rotate licenses dengan status "pending" atau "trial" (belum registered).
func (s *ProvisionService) RotateAll(ctx context.Context) error {
	const q = `
		SELECT id, license_key, project_id, company_id, product_id, plan, status,
		       modules, apps, contract_amount, description,
		       max_users, max_trans_per_month, max_trans_per_day,
		       max_items, max_customers, max_branches, max_storage,
		       expires_at, instance_url, instance_name, provision_api_key,
		       check_interval, last_pull_at, is_registered, proposal_id,
		       created_by, created_at, updated_at, deleted_at, archived_at
		FROM client_licenses
		WHERE deleted_at IS NULL AND is_registered = FALSE
		  AND provision_api_key IS NOT NULL
	`

	var licenses []*domain.ClientLicense
	rows, err := s.db.QueryxContext(ctx, q)
	if err != nil {
		return fmt.Errorf("RotateAll query: %w", err)
	}
	defer rows.Close()

	// Scan ke intermediate struct
	type licenseRow struct {
		domain.ClientLicense
	}

	for rows.Next() {
		var l licenseRow
		if err := rows.StructScan(&l); err != nil {
			s.log.Error("RotateAll: scan failed", zap.Error(err))
			continue
		}
		licenses = append(licenses, &l.ClientLicense)
	}

	rotateCount := 0
	for _, lic := range licenses {
		if err := s.rotateLicense(ctx, lic); err != nil {
			s.log.Error("RotateAll: rotate license failed",
				zap.String("license_id", lic.ID.String()),
				zap.Error(err))
			continue
		}
		rotateCount++
	}

	s.log.Info("RotateAll completed", zap.Int("rotated", rotateCount), zap.Int("total", len(licenses)))
	return nil
}

// rotateLicense adalah internal function yang rotate satu license.
// Algoritma:
//  1. Generate key baru
//  2. Backup key lama ke previous + previous_at = now
//  3. Set key baru + generated_at = now
//  4. Update ke database
func (s *ProvisionService) rotateLicense(ctx context.Context, lic *domain.ClientLicense) error {
	newKey, err := s.GenerateProvisionKey()
	if err != nil {
		return fmt.Errorf("rotateLicense: generate: %w", err)
	}

	const q = `
		UPDATE client_licenses
		SET provision_api_key_previous = provision_api_key,
		    provision_api_key_previous_at = provision_api_key_generated_at,
		    provision_api_key = $1,
		    provision_api_key_generated_at = NOW(),
		    updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`

	res, err := s.db.ExecContext(ctx, q, newKey, lic.ID)
	if err != nil {
		return fmt.Errorf("rotateLicense: exec: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rotateLicense: rows affected: %w", err)
	}

	if n == 0 {
		return fmt.Errorf("rotateLicense: license not found: %s", lic.ID)
	}

	s.log.Debug("License provision key rotated",
		zap.String("license_id", lic.ID.String()),
		zap.String("license_key", lic.LicenseKey))

	return nil
}
