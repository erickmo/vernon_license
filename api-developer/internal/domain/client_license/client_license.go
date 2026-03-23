package client_license

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	StatusTrial     = "trial"
	StatusActive    = "active"
	StatusSuspended = "suspended"
	StatusExpired   = "expired"

	PlanSaaS      = "saas"
	PlanDedicated = "dedicated"
)

var (
	ErrNotFound           = errors.New("lisensi tidak ditemukan")
	ErrAlreadyProvisioned = errors.New("lisensi sudah di-provisioning")
	ErrMissingFlashERPURL = errors.New("URL FlashERP belum dikonfigurasi")
)

type ClientLicense struct {
	ID               uuid.UUID  `db:"id"                json:"id"`
	LicenseKey       string     `db:"license_key"       json:"license_key"`
	ClientName       string     `db:"client_name"       json:"client_name"`
	ClientEmail      string     `db:"client_email"      json:"client_email"`
	Product          string     `db:"product"           json:"product"`
	Plan             string     `db:"plan"              json:"plan"`
	Status           string     `db:"status"            json:"status"`
	MaxUsers         *int       `db:"max_users"         json:"max_users"`
	MaxTransPerMonth *int       `db:"max_trans_per_month" json:"max_trans_per_month"`
	MaxTransPerDay   *int       `db:"max_trans_per_day" json:"max_trans_per_day"`
	MaxItems         *int       `db:"max_items"         json:"max_items"`
	MaxCustomers     *int       `db:"max_customers"     json:"max_customers"`
	MaxBranches      *int       `db:"max_branches"      json:"max_branches"`
	ExpiresAt        *time.Time `db:"expires_at"        json:"expires_at"`
	FlashERPURL      *string    `db:"flasherp_url"      json:"flasherp_url"`
	ProvisionAPIKey  *string    `db:"provision_api_key" json:"-"`
	IsProvisioned    bool       `db:"is_provisioned"    json:"is_provisioned"`
	CreatedBy        uuid.UUID  `db:"created_by"        json:"created_by"`
	CreatedAt        time.Time  `db:"created_at"        json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"        json:"updated_at"`
}

func NewClientLicense(clientName, clientEmail, product, plan string, createdBy uuid.UUID) *ClientLicense {
	return &ClientLicense{
		ID:          uuid.New(),
		ClientName:  clientName,
		ClientEmail: clientEmail,
		Product:     product,
		Plan:        plan,
		Status:      StatusActive,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func (l *ClientLicense) Suspend() {
	l.Status = StatusSuspended
	l.UpdatedAt = time.Now().UTC()
}

func (l *ClientLicense) Activate() {
	l.Status = StatusActive
	l.UpdatedAt = time.Now().UTC()
}

func (l *ClientLicense) Expire() {
	l.Status = StatusExpired
	l.UpdatedAt = time.Now().UTC()
}

type ListFilter struct {
	Status   string
	Product  string
	Page     int
	PageSize int
	Search   string
}

type WriteRepository interface {
	Save(ctx context.Context, l *ClientLicense) error
	Update(ctx context.Context, l *ClientLicense) error
	GetByID(ctx context.Context, id uuid.UUID) (*ClientLicense, error)
	GetByKey(ctx context.Context, key string) (*ClientLicense, error)
}

type ReadRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*ClientLicense, error)
	GetByKey(ctx context.Context, key string) (*ClientLicense, error)
	List(ctx context.Context, filter ListFilter) ([]*ClientLicense, int, error)
}
