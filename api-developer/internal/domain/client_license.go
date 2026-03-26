package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ClientLicense merepresentasikan lisensi yang diberikan kepada klien untuk sebuah produk.
type ClientLicense struct {
	ID         uuid.UUID `db:"id"`
	LicenseKey string    `db:"license_key"`
	ProjectID  *uuid.UUID `db:"project_id"`
	CompanyID  *uuid.UUID `db:"company_id"`
	ProductID  uuid.UUID `db:"product_id"`
	Plan       string    `db:"plan"`
	// Status adalah salah satu dari: "active" | "trial" | "suspended" | "expired" | "pending"
	Status           string     `db:"status"`
	Modules          []string   `db:"modules"`
	Apps             []string   `db:"apps"`
	ContractAmount   *float64   `db:"contract_amount"`
	Description      *string    `db:"description"`
	MaxUsers         *int       `db:"max_users"`
	MaxTransPerMonth *int       `db:"max_trans_per_month"`
	MaxTransPerDay   *int       `db:"max_trans_per_day"`
	MaxItems         *int       `db:"max_items"`
	MaxCustomers     *int       `db:"max_customers"`
	MaxBranches      *int       `db:"max_branches"`
	MaxStorage       *int       `db:"max_storage"`
	ExpiresAt        *time.Time `db:"expires_at"`
	InstanceURL      *string    `db:"instance_url"`
	InstanceName     *string    `db:"instance_name"`
	// OTP tidak pernah diekspos ke luar (json:"-").
	OTP            *string    `db:"otp" json:"-"`
	OTPGeneratedAt *time.Time `db:"otp_generated_at"`
	OTPPrevious    *string    `db:"otp_previous" json:"-"`
	OTPPreviousAt  *time.Time `db:"otp_previous_at"`
	ClientAppIP                       *string    `db:"client_app_ip"`
	SuperuserUsername                  *string    `db:"superuser_username"`
	CheckInterval                     string     `db:"check_interval"`
	LastPullAt                        *time.Time `db:"last_pull_at"`
	IsRegistered                      bool       `db:"is_registered"`
	ProposalID                        *uuid.UUID `db:"proposal_id"`
	CreatedBy                         *uuid.UUID `db:"created_by"`
	CreatedAt                         time.Time  `db:"created_at"`
	UpdatedAt                         time.Time  `db:"updated_at"`
	DeletedAt                         *time.Time `db:"deleted_at"`
	ArchivedAt                        *time.Time `db:"archived_at"`
}

// IsValid menentukan apakah license valid untuk validate endpoint.
func (l *ClientLicense) IsValid() bool {
	if l.Status == "active" {
		return l.ExpiresAt == nil || l.ExpiresAt.After(time.Now())
	}
	if l.Status == "trial" {
		return true // approved trial
	}
	return false
}

// ValidReason mengembalikan reason jika license tidak valid.
func (l *ClientLicense) ValidReason() string {
	switch l.Status {
	case "suspended":
		return "suspended"
	case "expired":
		return "expired"
	case "pending", "trial":
		return "pending_approval"
	default:
		if l.ExpiresAt != nil && !l.ExpiresAt.After(time.Now()) {
			return "expired"
		}
		return ""
	}
}

// LicenseRepository mendefinisikan operasi persistence untuk entitas ClientLicense.
type LicenseRepository interface {
	// FindByID mencari license berdasarkan UUID.
	FindByID(ctx context.Context, id uuid.UUID) (*ClientLicense, error)

	// FindByKey mencari license berdasarkan license_key.
	FindByKey(ctx context.Context, key string) (*ClientLicense, error)

	// FindByOTP mencari license berdasarkan otp dan product slug.
	// Hanya current code yang valid — tidak ada grace period dengan previous code.
	// Previous code disimpan untuk audit trail saja.
	FindByOTP(ctx context.Context, otp, productSlug string) (*ClientLicense, error)

	// FindByProject mengembalikan semua license untuk sebuah project.
	FindByProject(ctx context.Context, projectID uuid.UUID) ([]*ClientLicense, error)

	// FindAll mengembalikan semua license yang belum dihapus.
	FindAll(ctx context.Context) ([]*ClientLicense, error)

	// FindExpiring mengembalikan license yang akan expired dalam withinDays hari.
	FindExpiring(ctx context.Context, withinDays int) ([]*ClientLicense, error)

	// FindByCompany mengembalikan semua license untuk sebuah company.
	FindByCompany(ctx context.Context, companyID uuid.UUID) ([]*ClientLicense, error)

	// FindByCompanyAndProduct mencari license aktif berdasarkan company dan product.
	// Digunakan untuk cek duplikasi saat register.
	FindByCompanyAndProduct(ctx context.Context, companyID, productID uuid.UUID) (*ClientLicense, error)

	// Create menyimpan license baru ke database.
	Create(ctx context.Context, l *ClientLicense) error

	// Update memperbarui data license.
	Update(ctx context.Context, l *ClientLicense) error

	// UpdateRegistration memperbarui instance_url dan instance_name saat register.
	UpdateRegistration(ctx context.Context, id uuid.UUID, instanceURL, instanceName string) error

	// UpdateLastPullAt memperbarui last_pull_at ke waktu sekarang.
	UpdateLastPullAt(ctx context.Context, id uuid.UUID) error

	// UpdateStatus memperbarui status license.
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error

	// UpdateSuperuser memperbarui superuser_username pada license.
	UpdateSuperuser(ctx context.Context, id uuid.UUID, username string) error
}
