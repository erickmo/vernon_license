package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/flashlab/vernon-license/internal/domain"
)

// LicenseRepo adalah implementasi domain.LicenseRepository berbasis PostgreSQL.
type LicenseRepo struct {
	db *sqlx.DB
}

// NewLicenseRepo membuat instance LicenseRepo baru.
func NewLicenseRepo(db *sqlx.DB) *LicenseRepo {
	return &LicenseRepo{db: db}
}

// licenseRow adalah intermediate struct untuk scan TEXT[] fields dari PostgreSQL.
type licenseRow struct {
	ID         uuid.UUID      `db:"id"`
	LicenseKey string         `db:"license_key"`
	ProjectID  uuid.UUID      `db:"project_id"`
	CompanyID  uuid.UUID      `db:"company_id"`
	ProductID  uuid.UUID      `db:"product_id"`
	Plan       string         `db:"plan"`
	Status     string         `db:"status"`
	Modules    pq.StringArray `db:"modules"`
	Apps       pq.StringArray `db:"apps"`
	ContractAmount  *float64   `db:"contract_amount"`
	Description     *string    `db:"description"`
	MaxUsers        *int       `db:"max_users"`
	MaxTransPerMonth *int      `db:"max_trans_per_month"`
	MaxTransPerDay  *int       `db:"max_trans_per_day"`
	MaxItems        *int       `db:"max_items"`
	MaxCustomers    *int       `db:"max_customers"`
	MaxBranches     *int       `db:"max_branches"`
	MaxStorage      *int       `db:"max_storage"`
	ExpiresAt       *time.Time `db:"expires_at"`
	InstanceURL     *string    `db:"instance_url"`
	InstanceName    *string    `db:"instance_name"`
	ProvisionAPIKey *string    `db:"provision_api_key"`
	CheckInterval   string     `db:"check_interval"`
	LastPullAt      *time.Time `db:"last_pull_at"`
	IsRegistered    bool       `db:"is_registered"`
	ProposalID      *uuid.UUID `db:"proposal_id"`
	CreatedBy       uuid.UUID  `db:"created_by"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	ArchivedAt      *time.Time `db:"archived_at"`
}

func (lr *licenseRow) toDomain() *domain.ClientLicense {
	return &domain.ClientLicense{
		ID:               lr.ID,
		LicenseKey:       lr.LicenseKey,
		ProjectID:        lr.ProjectID,
		CompanyID:        lr.CompanyID,
		ProductID:        lr.ProductID,
		Plan:             lr.Plan,
		Status:           lr.Status,
		Modules:          []string(lr.Modules),
		Apps:             []string(lr.Apps),
		ContractAmount:   lr.ContractAmount,
		Description:      lr.Description,
		MaxUsers:         lr.MaxUsers,
		MaxTransPerMonth: lr.MaxTransPerMonth,
		MaxTransPerDay:   lr.MaxTransPerDay,
		MaxItems:         lr.MaxItems,
		MaxCustomers:     lr.MaxCustomers,
		MaxBranches:      lr.MaxBranches,
		MaxStorage:       lr.MaxStorage,
		ExpiresAt:        lr.ExpiresAt,
		InstanceURL:      lr.InstanceURL,
		InstanceName:     lr.InstanceName,
		ProvisionAPIKey:  lr.ProvisionAPIKey,
		CheckInterval:    lr.CheckInterval,
		LastPullAt:       lr.LastPullAt,
		IsRegistered:     lr.IsRegistered,
		ProposalID:       lr.ProposalID,
		CreatedBy:        lr.CreatedBy,
		CreatedAt:        lr.CreatedAt,
		UpdatedAt:        lr.UpdatedAt,
		DeletedAt:        lr.DeletedAt,
		ArchivedAt:       lr.ArchivedAt,
	}
}

const licenseSelectCols = `
	id, license_key, project_id, company_id, product_id, plan, status,
	modules, apps, contract_amount, description,
	max_users, max_trans_per_month, max_trans_per_day,
	max_items, max_customers, max_branches, max_storage,
	expires_at, instance_url, instance_name, provision_api_key,
	check_interval, last_pull_at, is_registered, proposal_id,
	created_by, created_at, updated_at, deleted_at, archived_at`

// FindByID mencari license berdasarkan UUID.
func (r *LicenseRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.ClientLicense, error) {
	var row licenseRow
	q := `SELECT ` + licenseSelectCols + ` FROM client_licenses WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &row, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrLicenseNotFound
		}
		return nil, fmt.Errorf("LicenseRepo.FindByID: %w", err)
	}
	return row.toDomain(), nil
}

// FindByKey mencari license berdasarkan license_key.
func (r *LicenseRepo) FindByKey(ctx context.Context, key string) (*domain.ClientLicense, error) {
	var row licenseRow
	q := `SELECT ` + licenseSelectCols + ` FROM client_licenses WHERE license_key = $1 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &row, q, key); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrLicenseNotFound
		}
		return nil, fmt.Errorf("LicenseRepo.FindByKey: %w", err)
	}
	return row.toDomain(), nil
}

// FindByProvisionKey mencari license berdasarkan provision_api_key dan product slug.
func (r *LicenseRepo) FindByProvisionKey(ctx context.Context, provisionKey, productSlug string) (*domain.ClientLicense, error) {
	var row licenseRow
	q := `SELECT cl.` + licenseSelectCols + `
		FROM client_licenses cl
		JOIN products p ON p.id = cl.product_id
		WHERE cl.provision_api_key = $1
		  AND p.slug = $2
		  AND cl.deleted_at IS NULL`
	// Re-qualify columns since we use alias
	q = `SELECT
		cl.id, cl.license_key, cl.project_id, cl.company_id, cl.product_id, cl.plan, cl.status,
		cl.modules, cl.apps, cl.contract_amount, cl.description,
		cl.max_users, cl.max_trans_per_month, cl.max_trans_per_day,
		cl.max_items, cl.max_customers, cl.max_branches, cl.max_storage,
		cl.expires_at, cl.instance_url, cl.instance_name, cl.provision_api_key,
		cl.check_interval, cl.last_pull_at, cl.is_registered, cl.proposal_id,
		cl.created_by, cl.created_at, cl.updated_at, cl.deleted_at, cl.archived_at
	FROM client_licenses cl
	JOIN products p ON p.id = cl.product_id
	WHERE cl.provision_api_key = $1
	  AND p.slug = $2
	  AND cl.deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &row, q, provisionKey, productSlug); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrLicenseNotFound
		}
		return nil, fmt.Errorf("LicenseRepo.FindByProvisionKey: %w", err)
	}
	return row.toDomain(), nil
}

// FindByProject mengembalikan semua license untuk sebuah project.
func (r *LicenseRepo) FindByProject(ctx context.Context, projectID uuid.UUID) ([]*domain.ClientLicense, error) {
	q := `SELECT ` + licenseSelectCols + `
		FROM client_licenses
		WHERE project_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`
	var rows []licenseRow
	if err := r.db.SelectContext(ctx, &rows, q, projectID); err != nil {
		return nil, fmt.Errorf("LicenseRepo.FindByProject: %w", err)
	}
	return licenseRowsToDomain(rows), nil
}

// FindAll mengembalikan semua license yang belum dihapus.
func (r *LicenseRepo) FindAll(ctx context.Context) ([]*domain.ClientLicense, error) {
	q := `SELECT ` + licenseSelectCols + `
		FROM client_licenses
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`
	var rows []licenseRow
	if err := r.db.SelectContext(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("LicenseRepo.FindAll: %w", err)
	}
	return licenseRowsToDomain(rows), nil
}

// FindExpiring mengembalikan license yang akan expired dalam withinDays hari.
func (r *LicenseRepo) FindExpiring(ctx context.Context, withinDays int) ([]*domain.ClientLicense, error) {
	q := `SELECT ` + licenseSelectCols + `
		FROM client_licenses
		WHERE deleted_at IS NULL
		  AND status = 'active'
		  AND expires_at IS NOT NULL
		  AND expires_at <= NOW() + ($1 || ' days')::INTERVAL
		  AND expires_at > NOW()
		ORDER BY expires_at ASC`
	var rows []licenseRow
	if err := r.db.SelectContext(ctx, &rows, q, withinDays); err != nil {
		return nil, fmt.Errorf("LicenseRepo.FindExpiring: %w", err)
	}
	return licenseRowsToDomain(rows), nil
}

// Create menyimpan license baru ke database.
func (r *LicenseRepo) Create(ctx context.Context, l *domain.ClientLicense) error {
	const q = `
		INSERT INTO client_licenses
		    (id, license_key, project_id, company_id, product_id, plan, status,
		     modules, apps, contract_amount, description,
		     max_users, max_trans_per_month, max_trans_per_day,
		     max_items, max_customers, max_branches, max_storage,
		     expires_at, instance_url, instance_name, provision_api_key,
		     check_interval, last_pull_at, is_registered, proposal_id,
		     created_by, created_at, updated_at)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7,
		     $8, $9, $10, $11,
		     $12, $13, $14,
		     $15, $16, $17, $18,
		     $19, $20, $21, $22,
		     $23, $24, $25, $26,
		     $27, NOW(), NOW())
		RETURNING created_at, updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		l.ID, l.LicenseKey, l.ProjectID, l.CompanyID, l.ProductID, l.Plan, l.Status,
		pq.Array(l.Modules), pq.Array(l.Apps), l.ContractAmount, l.Description,
		l.MaxUsers, l.MaxTransPerMonth, l.MaxTransPerDay,
		l.MaxItems, l.MaxCustomers, l.MaxBranches, l.MaxStorage,
		l.ExpiresAt, l.InstanceURL, l.InstanceName, l.ProvisionAPIKey,
		l.CheckInterval, l.LastPullAt, l.IsRegistered, l.ProposalID,
		l.CreatedBy,
	).Scan(&l.CreatedAt, &l.UpdatedAt); err != nil {
		return fmt.Errorf("LicenseRepo.Create: %w", err)
	}
	return nil
}

// Update memperbarui data license.
func (r *LicenseRepo) Update(ctx context.Context, l *domain.ClientLicense) error {
	const q = `
		UPDATE client_licenses
		SET plan = $1, status = $2, modules = $3, apps = $4,
		    contract_amount = $5, description = $6,
		    max_users = $7, max_trans_per_month = $8, max_trans_per_day = $9,
		    max_items = $10, max_customers = $11, max_branches = $12, max_storage = $13,
		    expires_at = $14, check_interval = $15,
		    archived_at = $16, updated_at = NOW()
		WHERE id = $17 AND deleted_at IS NULL
		RETURNING updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		l.Plan, l.Status, pq.Array(l.Modules), pq.Array(l.Apps),
		l.ContractAmount, l.Description,
		l.MaxUsers, l.MaxTransPerMonth, l.MaxTransPerDay,
		l.MaxItems, l.MaxCustomers, l.MaxBranches, l.MaxStorage,
		l.ExpiresAt, l.CheckInterval,
		l.ArchivedAt, l.ID,
	).Scan(&l.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrLicenseNotFound
		}
		return fmt.Errorf("LicenseRepo.Update: %w", err)
	}
	return nil
}

// UpdateRegistration memperbarui instance_url, instance_name, dan is_registered.
func (r *LicenseRepo) UpdateRegistration(ctx context.Context, id uuid.UUID, instanceURL, instanceName string) error {
	const q = `
		UPDATE client_licenses
		SET instance_url = $1, instance_name = $2, is_registered = TRUE, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, instanceURL, instanceName, id)
	if err != nil {
		return fmt.Errorf("LicenseRepo.UpdateRegistration: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("LicenseRepo.UpdateRegistration: rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrLicenseNotFound
	}
	return nil
}

// UpdateLastPullAt memperbarui last_pull_at ke waktu sekarang.
func (r *LicenseRepo) UpdateLastPullAt(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE client_licenses SET last_pull_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("LicenseRepo.UpdateLastPullAt: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("LicenseRepo.UpdateLastPullAt: rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrLicenseNotFound
	}
	return nil
}

// UpdateStatus memperbarui status license.
func (r *LicenseRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	const q = `UPDATE client_licenses SET status = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, status, id)
	if err != nil {
		return fmt.Errorf("LicenseRepo.UpdateStatus: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("LicenseRepo.UpdateStatus: rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrLicenseNotFound
	}
	return nil
}

func licenseRowsToDomain(rows []licenseRow) []*domain.ClientLicense {
	result := make([]*domain.ClientLicense, len(rows))
	for i, row := range rows {
		r := row
		result[i] = r.toDomain()
	}
	return result
}
