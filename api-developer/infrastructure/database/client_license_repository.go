package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	clientlicense "github.com/flashlab/flasherp-developer-api/internal/domain/client_license"
)

type licenseRow struct {
	ID               uuid.UUID      `db:"id"`
	LicenseKey       string         `db:"license_key"`
	ClientName       string         `db:"client_name"`
	ClientEmail      string         `db:"client_email"`
	Product          string         `db:"product"`
	Plan             string         `db:"plan"`
	Status           string         `db:"status"`
	MaxUsers         sql.NullInt64  `db:"max_users"`
	MaxTransPerMonth sql.NullInt64  `db:"max_trans_per_month"`
	MaxTransPerDay   sql.NullInt64  `db:"max_trans_per_day"`
	MaxItems         sql.NullInt64  `db:"max_items"`
	MaxCustomers     sql.NullInt64  `db:"max_customers"`
	MaxBranches      sql.NullInt64  `db:"max_branches"`
	ExpiresAt        sql.NullTime   `db:"expires_at"`
	FlashERPURL      sql.NullString `db:"flasherp_url"`
	ProvisionAPIKey  sql.NullString `db:"provision_api_key"`
	IsProvisioned    bool           `db:"is_provisioned"`
	CreatedBy        uuid.UUID      `db:"created_by"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}

func nullInt(v *int) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*v), Valid: true}
}

func nullTime(v *time.Time) sql.NullTime {
	if v == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *v, Valid: true}
}

func nullStr(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}

func intPtr(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}
	i := int(v.Int64)
	return &i
}

func timePtr(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	return &v.Time
}

func strPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}

func toLicenseDomain(row *licenseRow) *clientlicense.ClientLicense {
	return &clientlicense.ClientLicense{
		ID:               row.ID,
		LicenseKey:       row.LicenseKey,
		ClientName:       row.ClientName,
		ClientEmail:      row.ClientEmail,
		Product:          row.Product,
		Plan:             row.Plan,
		Status:           row.Status,
		MaxUsers:         intPtr(row.MaxUsers),
		MaxTransPerMonth: intPtr(row.MaxTransPerMonth),
		MaxTransPerDay:   intPtr(row.MaxTransPerDay),
		MaxItems:         intPtr(row.MaxItems),
		MaxCustomers:     intPtr(row.MaxCustomers),
		MaxBranches:      intPtr(row.MaxBranches),
		ExpiresAt:        timePtr(row.ExpiresAt),
		FlashERPURL:      strPtr(row.FlashERPURL),
		ProvisionAPIKey:  strPtr(row.ProvisionAPIKey),
		IsProvisioned:    row.IsProvisioned,
		CreatedBy:        row.CreatedBy,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}

type ClientLicenseRepository struct {
	db *sqlx.DB
}

func NewClientLicenseRepository(db *sqlx.DB) *ClientLicenseRepository {
	return &ClientLicenseRepository{db: db}
}

func (r *ClientLicenseRepository) Save(ctx context.Context, l *clientlicense.ClientLicense) error {
	query := `
		INSERT INTO client_licenses (
			id, license_key, client_name, client_email, product, plan, status,
			max_users, max_trans_per_month, max_trans_per_day, max_items, max_customers, max_branches,
			expires_at, flasherp_url, provision_api_key, is_provisioned, created_by, created_at, updated_at
		) VALUES (
			:id, :license_key, :client_name, :client_email, :product, :plan, :status,
			:max_users, :max_trans_per_month, :max_trans_per_day, :max_items, :max_customers, :max_branches,
			:expires_at, :flasherp_url, :provision_api_key, :is_provisioned, :created_by, :created_at, :updated_at
		)`
	row := &licenseRow{
		ID:               l.ID,
		LicenseKey:       l.LicenseKey,
		ClientName:       l.ClientName,
		ClientEmail:      l.ClientEmail,
		Product:          l.Product,
		Plan:             l.Plan,
		Status:           l.Status,
		MaxUsers:         nullInt(l.MaxUsers),
		MaxTransPerMonth: nullInt(l.MaxTransPerMonth),
		MaxTransPerDay:   nullInt(l.MaxTransPerDay),
		MaxItems:         nullInt(l.MaxItems),
		MaxCustomers:     nullInt(l.MaxCustomers),
		MaxBranches:      nullInt(l.MaxBranches),
		ExpiresAt:        nullTime(l.ExpiresAt),
		FlashERPURL:      nullStr(l.FlashERPURL),
		ProvisionAPIKey:  nullStr(l.ProvisionAPIKey),
		IsProvisioned:    l.IsProvisioned,
		CreatedBy:        l.CreatedBy,
		CreatedAt:        l.CreatedAt,
		UpdatedAt:        l.UpdatedAt,
	}
	_, err := r.db.NamedExecContext(ctx, query, row)
	return err
}

func (r *ClientLicenseRepository) Update(ctx context.Context, l *clientlicense.ClientLicense) error {
	query := `
		UPDATE client_licenses SET
			client_name = :client_name,
			client_email = :client_email,
			product = :product,
			plan = :plan,
			status = :status,
			max_users = :max_users,
			max_trans_per_month = :max_trans_per_month,
			max_trans_per_day = :max_trans_per_day,
			max_items = :max_items,
			max_customers = :max_customers,
			max_branches = :max_branches,
			expires_at = :expires_at,
			flasherp_url = :flasherp_url,
			provision_api_key = :provision_api_key,
			is_provisioned = :is_provisioned,
			updated_at = :updated_at
		WHERE id = :id`
	row := &licenseRow{
		ID:               l.ID,
		ClientName:       l.ClientName,
		ClientEmail:      l.ClientEmail,
		Product:          l.Product,
		Plan:             l.Plan,
		Status:           l.Status,
		MaxUsers:         nullInt(l.MaxUsers),
		MaxTransPerMonth: nullInt(l.MaxTransPerMonth),
		MaxTransPerDay:   nullInt(l.MaxTransPerDay),
		MaxItems:         nullInt(l.MaxItems),
		MaxCustomers:     nullInt(l.MaxCustomers),
		MaxBranches:      nullInt(l.MaxBranches),
		ExpiresAt:        nullTime(l.ExpiresAt),
		FlashERPURL:      nullStr(l.FlashERPURL),
		ProvisionAPIKey:  nullStr(l.ProvisionAPIKey),
		IsProvisioned:    l.IsProvisioned,
		UpdatedAt:        time.Now().UTC(),
	}
	_, err := r.db.NamedExecContext(ctx, query, row)
	return err
}

func (r *ClientLicenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*clientlicense.ClientLicense, error) {
	var row licenseRow
	err := r.db.GetContext(ctx, &row, `SELECT * FROM client_licenses WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, clientlicense.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toLicenseDomain(&row), nil
}

func (r *ClientLicenseRepository) GetByKey(ctx context.Context, key string) (*clientlicense.ClientLicense, error) {
	var row licenseRow
	err := r.db.GetContext(ctx, &row, `SELECT * FROM client_licenses WHERE license_key = $1`, key)
	if err == sql.ErrNoRows {
		return nil, clientlicense.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toLicenseDomain(&row), nil
}

func (r *ClientLicenseRepository) List(ctx context.Context, filter clientlicense.ListFilter) ([]*clientlicense.ClientLicense, int, error) {
	conditions := []string{}
	args := []any{}
	argIdx := 1

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Product != "" {
		conditions = append(conditions, fmt.Sprintf("product = $%d", argIdx))
		args = append(args, filter.Product)
		argIdx++
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(client_name ILIKE $%d OR client_email ILIKE $%d OR license_key ILIKE $%d)", argIdx, argIdx+1, argIdx+2))
		search := "%" + filter.Search + "%"
		args = append(args, search, search, search)
		argIdx += 3
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM client_licenses %s", where)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	listArgs := append(args, pageSize, offset)
	listQuery := fmt.Sprintf(
		"SELECT * FROM client_licenses %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, argIdx, argIdx+1,
	)

	var rows []licenseRow
	if err := r.db.SelectContext(ctx, &rows, listQuery, listArgs...); err != nil {
		return nil, 0, err
	}

	licenses := make([]*clientlicense.ClientLicense, len(rows))
	for i := range rows {
		licenses[i] = toLicenseDomain(&rows[i])
	}
	return licenses, total, nil
}
