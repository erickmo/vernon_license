package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Proposal merepresentasikan proposal lisensi yang diajukan oleh sales.
type Proposal struct {
	ID        uuid.UUID `db:"id"`
	ProjectID uuid.UUID `db:"project_id"`
	CompanyID uuid.UUID `db:"company_id"`
	ProductID uuid.UUID `db:"product_id"`
	Version   int       `db:"version"`
	// Status adalah salah satu dari: "draft" | "submitted" | "approved" | "rejected"
	Status          string          `db:"status"`
	Modules         []string        `db:"modules"`
	Apps            []string        `db:"apps"`
	Plan            string          `db:"plan"`
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
	Changelog       json.RawMessage `db:"changelog"`
	PDFPath         *string         `db:"pdf_path"`
	PDFGeneratedAt  *time.Time      `db:"pdf_generated_at"`
	SubmittedBy     uuid.UUID       `db:"submitted_by"`
	ReviewedBy      *uuid.UUID      `db:"reviewed_by"`
	ReviewedAt      *time.Time      `db:"reviewed_at"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}

// ChangelogEntry adalah struktur diff satu field antar versi proposal.
type ChangelogEntry struct {
	Field    string `json:"field"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

// Changelog adalah struktur lengkap changelog proposal dibanding versi sebelumnya.
type Changelog struct {
	ComparedToVersion int              `json:"compared_to_version"`
	Summary           string           `json:"summary"`
	Changes           []ChangelogEntry `json:"changes"`
	Unchanged         []string         `json:"unchanged"`
}

// comparableFields mendefinisikan field-field yang dibandingkan dalam changelog.
type comparableField struct {
	name string
	prev any
	curr any
}

// ComputeChangelog membandingkan proposal curr dengan proposal prev dan mengembalikan Changelog.
func ComputeChangelog(prev, curr *Proposal) Changelog {
	fields := []comparableField{
		{"plan", prev.Plan, curr.Plan},
		{"modules", prev.Modules, curr.Modules},
		{"apps", prev.Apps, curr.Apps},
		{"contract_amount", prev.ContractAmount, curr.ContractAmount},
		{"expires_at", prev.ExpiresAt, curr.ExpiresAt},
		{"max_users", prev.MaxUsers, curr.MaxUsers},
		{"max_trans_per_month", prev.MaxTransPerMonth, curr.MaxTransPerMonth},
		{"max_trans_per_day", prev.MaxTransPerDay, curr.MaxTransPerDay},
		{"max_items", prev.MaxItems, curr.MaxItems},
		{"max_customers", prev.MaxCustomers, curr.MaxCustomers},
		{"max_branches", prev.MaxBranches, curr.MaxBranches},
		{"max_storage", prev.MaxStorage, curr.MaxStorage},
		{"notes", prev.Notes, curr.Notes},
	}

	var changes []ChangelogEntry
	var unchanged []string

	for _, f := range fields {
		if !reflect.DeepEqual(f.prev, f.curr) {
			changes = append(changes, ChangelogEntry{
				Field:    f.name,
				OldValue: f.prev,
				NewValue: f.curr,
			})
		} else {
			unchanged = append(unchanged, f.name)
		}
	}

	summary := fmt.Sprintf("%d field(s) changed from v%d to v%d", len(changes), prev.Version, curr.Version)

	return Changelog{
		ComparedToVersion: prev.Version,
		Summary:           summary,
		Changes:           changes,
		Unchanged:         unchanged,
	}
}

// ParseChangelog mem-parse json.RawMessage changelog menjadi struct Changelog.
// Mengembalikan error jika JSON tidak valid.
func ParseChangelog(raw json.RawMessage) (*Changelog, error) {
	var cl Changelog
	if err := json.Unmarshal(raw, &cl); err != nil {
		return nil, fmt.Errorf("ParseChangelog: %w", err)
	}
	return &cl, nil
}

// ProposalRepository mendefinisikan operasi persistence untuk entitas Proposal.
type ProposalRepository interface {
	// FindByID mencari proposal berdasarkan UUID.
	FindByID(ctx context.Context, id uuid.UUID) (*Proposal, error)

	// FindByProject mengembalikan semua proposal untuk sebuah project.
	FindByProject(ctx context.Context, projectID uuid.UUID) ([]*Proposal, error)

	// FindLatestByProjectProduct mencari proposal terbaru untuk kombinasi project + product.
	FindLatestByProjectProduct(ctx context.Context, projectID, productID uuid.UUID) (*Proposal, error)

	// FindAll mengembalikan semua proposal.
	FindAll(ctx context.Context) ([]*Proposal, error)

	// Create menyimpan proposal baru ke database.
	Create(ctx context.Context, p *Proposal) error

	// Update memperbarui data proposal.
	Update(ctx context.Context, p *Proposal) error

	// NextVersion mengembalikan version berikutnya untuk kombinasi project + product.
	NextVersion(ctx context.Context, projectID, productID uuid.UUID) (int, error)
}
