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

// ProposalRepo adalah implementasi domain.ProposalRepository berbasis PostgreSQL.
type ProposalRepo struct {
	db *sqlx.DB
}

// NewProposalRepo membuat instance ProposalRepo baru.
func NewProposalRepo(db *sqlx.DB) *ProposalRepo {
	return &ProposalRepo{db: db}
}

// proposalRow adalah intermediate struct untuk scan TEXT[] fields dari PostgreSQL.
type proposalRow struct {
	ID        uuid.UUID `db:"id"`
	ProjectID uuid.UUID `db:"project_id"`
	CompanyID uuid.UUID `db:"company_id"`
	ProductID uuid.UUID `db:"product_id"`
	Version   int       `db:"version"`
	Status    string    `db:"status"`
	Modules   pq.StringArray `db:"modules"`
	Apps      pq.StringArray `db:"apps"`
	Plan      string    `db:"plan"`
	MaxUsers        *int            `db:"max_users"`
	MaxTransPerMonth *int           `db:"max_trans_per_month"`
	MaxTransPerDay  *int            `db:"max_trans_per_day"`
	MaxItems        *int            `db:"max_items"`
	MaxCustomers    *int            `db:"max_customers"`
	MaxBranches     *int            `db:"max_branches"`
	MaxStorage      *int            `db:"max_storage"`
	ContractAmount  *float64        `db:"contract_amount"`
	ExpiresAt       *time.Time      `db:"expires_at"`
	Notes           *string         `db:"notes"`
	OwnerNotes      *string         `db:"owner_notes"`
	RejectionReason *string         `db:"rejection_reason"`
	Changelog       []byte          `db:"changelog"`
	PDFPath         *string         `db:"pdf_path"`
	PDFGeneratedAt  *time.Time      `db:"pdf_generated_at"`
	SubmittedBy     uuid.UUID       `db:"submitted_by"`
	ReviewedBy      *uuid.UUID      `db:"reviewed_by"`
	ReviewedAt      *time.Time      `db:"reviewed_at"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}

func (pr *proposalRow) toDomain() *domain.Proposal {
	return &domain.Proposal{
		ID:               pr.ID,
		ProjectID:        pr.ProjectID,
		CompanyID:        pr.CompanyID,
		ProductID:        pr.ProductID,
		Version:          pr.Version,
		Status:           pr.Status,
		Modules:          []string(pr.Modules),
		Apps:             []string(pr.Apps),
		Plan:             pr.Plan,
		MaxUsers:         pr.MaxUsers,
		MaxTransPerMonth: pr.MaxTransPerMonth,
		MaxTransPerDay:   pr.MaxTransPerDay,
		MaxItems:         pr.MaxItems,
		MaxCustomers:     pr.MaxCustomers,
		MaxBranches:      pr.MaxBranches,
		MaxStorage:       pr.MaxStorage,
		ContractAmount:   pr.ContractAmount,
		ExpiresAt:        pr.ExpiresAt,
		Notes:            pr.Notes,
		OwnerNotes:       pr.OwnerNotes,
		RejectionReason:  pr.RejectionReason,
		Changelog:        pr.Changelog,
		PDFPath:          pr.PDFPath,
		PDFGeneratedAt:   pr.PDFGeneratedAt,
		SubmittedBy:      pr.SubmittedBy,
		ReviewedBy:       pr.ReviewedBy,
		ReviewedAt:       pr.ReviewedAt,
		CreatedAt:        pr.CreatedAt,
		UpdatedAt:        pr.UpdatedAt,
	}
}

const proposalSelectCols = `
	id, project_id, company_id, product_id, version, status,
	modules, apps, plan,
	max_users, max_trans_per_month, max_trans_per_day,
	max_items, max_customers, max_branches, max_storage,
	contract_amount, expires_at, notes, owner_notes, rejection_reason,
	changelog, pdf_path, pdf_generated_at,
	submitted_by, reviewed_by, reviewed_at,
	created_at, updated_at`

// FindByID mencari proposal berdasarkan UUID.
func (r *ProposalRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Proposal, error) {
	var row proposalRow
	q := `SELECT ` + proposalSelectCols + ` FROM proposals WHERE id = $1`
	if err := r.db.GetContext(ctx, &row, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProposalNotFound
		}
		return nil, fmt.Errorf("ProposalRepo.FindByID: %w", err)
	}
	return row.toDomain(), nil
}

// FindByProject mengembalikan semua proposal untuk sebuah project.
func (r *ProposalRepo) FindByProject(ctx context.Context, projectID uuid.UUID) ([]*domain.Proposal, error) {
	q := `SELECT ` + proposalSelectCols + `
		FROM proposals
		WHERE project_id = $1
		ORDER BY version DESC`
	var rows []proposalRow
	if err := r.db.SelectContext(ctx, &rows, q, projectID); err != nil {
		return nil, fmt.Errorf("ProposalRepo.FindByProject: %w", err)
	}
	return proposalRowsToDomain(rows), nil
}

// FindLatestByProjectProduct mencari proposal terbaru untuk kombinasi project + product.
func (r *ProposalRepo) FindLatestByProjectProduct(ctx context.Context, projectID, productID uuid.UUID) (*domain.Proposal, error) {
	var row proposalRow
	q := `SELECT ` + proposalSelectCols + `
		FROM proposals
		WHERE project_id = $1 AND product_id = $2
		ORDER BY version DESC
		LIMIT 1`
	if err := r.db.GetContext(ctx, &row, q, projectID, productID); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProposalNotFound
		}
		return nil, fmt.Errorf("ProposalRepo.FindLatestByProjectProduct: %w", err)
	}
	return row.toDomain(), nil
}

// FindAll mengembalikan semua proposal.
func (r *ProposalRepo) FindAll(ctx context.Context) ([]*domain.Proposal, error) {
	q := `SELECT ` + proposalSelectCols + `
		FROM proposals
		ORDER BY created_at DESC`
	var rows []proposalRow
	if err := r.db.SelectContext(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("ProposalRepo.FindAll: %w", err)
	}
	return proposalRowsToDomain(rows), nil
}

// Create menyimpan proposal baru ke database.
func (r *ProposalRepo) Create(ctx context.Context, p *domain.Proposal) error {
	const q = `
		INSERT INTO proposals
		    (id, project_id, company_id, product_id, version, status,
		     modules, apps, plan,
		     max_users, max_trans_per_month, max_trans_per_day,
		     max_items, max_customers, max_branches, max_storage,
		     contract_amount, expires_at, notes, owner_notes, rejection_reason,
		     changelog, pdf_path, pdf_generated_at,
		     submitted_by, reviewed_by, reviewed_at,
		     created_at, updated_at)
		VALUES
		    ($1, $2, $3, $4, $5, $6,
		     $7, $8, $9,
		     $10, $11, $12,
		     $13, $14, $15, $16,
		     $17, $18, $19, $20, $21,
		     $22, $23, $24,
		     $25, $26, $27,
		     NOW(), NOW())
		RETURNING created_at, updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		p.ID, p.ProjectID, p.CompanyID, p.ProductID, p.Version, p.Status,
		pq.Array(p.Modules), pq.Array(p.Apps), p.Plan,
		p.MaxUsers, p.MaxTransPerMonth, p.MaxTransPerDay,
		p.MaxItems, p.MaxCustomers, p.MaxBranches, p.MaxStorage,
		p.ContractAmount, p.ExpiresAt, p.Notes, p.OwnerNotes, p.RejectionReason,
		p.Changelog, p.PDFPath, p.PDFGeneratedAt,
		p.SubmittedBy, p.ReviewedBy, p.ReviewedAt,
	).Scan(&p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("ProposalRepo.Create: %w", err)
	}
	return nil
}

// Update memperbarui data proposal.
func (r *ProposalRepo) Update(ctx context.Context, p *domain.Proposal) error {
	const q = `
		UPDATE proposals
		SET status = $1, modules = $2, apps = $3, plan = $4,
		    max_users = $5, max_trans_per_month = $6, max_trans_per_day = $7,
		    max_items = $8, max_customers = $9, max_branches = $10, max_storage = $11,
		    contract_amount = $12, expires_at = $13,
		    notes = $14, owner_notes = $15, rejection_reason = $16,
		    changelog = $17, pdf_path = $18, pdf_generated_at = $19,
		    reviewed_by = $20, reviewed_at = $21,
		    updated_at = NOW()
		WHERE id = $22
		RETURNING updated_at`
	if err := r.db.QueryRowContext(ctx, q,
		p.Status, pq.Array(p.Modules), pq.Array(p.Apps), p.Plan,
		p.MaxUsers, p.MaxTransPerMonth, p.MaxTransPerDay,
		p.MaxItems, p.MaxCustomers, p.MaxBranches, p.MaxStorage,
		p.ContractAmount, p.ExpiresAt,
		p.Notes, p.OwnerNotes, p.RejectionReason,
		p.Changelog, p.PDFPath, p.PDFGeneratedAt,
		p.ReviewedBy, p.ReviewedAt,
		p.ID,
	).Scan(&p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrProposalNotFound
		}
		return fmt.Errorf("ProposalRepo.Update: %w", err)
	}
	return nil
}

// NextVersion mengembalikan version berikutnya untuk kombinasi project + product.
func (r *ProposalRepo) NextVersion(ctx context.Context, projectID, productID uuid.UUID) (int, error) {
	var maxVersion sql.NullInt64
	const q = `SELECT MAX(version) FROM proposals WHERE project_id = $1 AND product_id = $2`
	if err := r.db.QueryRowContext(ctx, q, projectID, productID).Scan(&maxVersion); err != nil {
		return 0, fmt.Errorf("ProposalRepo.NextVersion: %w", err)
	}
	if !maxVersion.Valid {
		return 1, nil
	}
	return int(maxVersion.Int64) + 1, nil
}

func proposalRowsToDomain(rows []proposalRow) []*domain.Proposal {
	result := make([]*domain.Proposal, len(rows))
	for i, row := range rows {
		r := row
		result[i] = r.toDomain()
	}
	return result
}
