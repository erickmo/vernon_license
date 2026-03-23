package create_client_license

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	clientlicense "github.com/flashlab/flasherp-developer-api/internal/domain/client_license"
	"github.com/flashlab/flasherp-developer-api/pkg/licenseutil"
	"github.com/google/uuid"
)

type CreateClientLicense struct {
	ClientName          string
	ClientEmail         string
	Product             string
	Plan                string
	MaxUsers            *int
	MaxTransPerMonth    *int
	MaxTransPerDay      *int
	MaxItems            *int
	MaxCustomers        *int
	MaxBranches         *int
	ExpiresAt           *time.Time
	FlashERPURL         *string
	ProvisionAPIKey     *string
	CreatedBy           uuid.UUID
}

type CreateClientLicenseResult struct {
	LicenseKey      string
	ClientEmail     string
	InitialPassword string
}

type Handler struct {
	licenseRepo clientlicense.WriteRepository
}

func NewHandler(licenseRepo clientlicense.WriteRepository) *Handler {
	return &Handler{licenseRepo: licenseRepo}
}

func (h *Handler) Handle(ctx context.Context, cmd any) error {
	return fmt.Errorf("gunakan HandleCreate untuk mendapatkan result")
}

func (h *Handler) HandleCreate(ctx context.Context, cmd CreateClientLicense) (*CreateClientLicenseResult, error) {
	licenseKey, err := licenseutil.GenerateLicenseKey()
	if err != nil {
		return nil, fmt.Errorf("gagal generate license key: %w", err)
	}

	plainPassword, err := licenseutil.GenerateSecurePassword()
	if err != nil {
		return nil, fmt.Errorf("gagal generate password: %w", err)
	}

	_, err = bcrypt.GenerateFromPassword([]byte(plainPassword), 12)
	if err != nil {
		return nil, fmt.Errorf("gagal hash password: %w", err)
	}

	product := cmd.Product
	if product == "" {
		product = "flasherp"
	}

	license := clientlicense.NewClientLicense(cmd.ClientName, cmd.ClientEmail, product, cmd.Plan, cmd.CreatedBy)
	license.LicenseKey = licenseKey
	license.MaxUsers = cmd.MaxUsers
	license.MaxTransPerMonth = cmd.MaxTransPerMonth
	license.MaxTransPerDay = cmd.MaxTransPerDay
	license.MaxItems = cmd.MaxItems
	license.MaxCustomers = cmd.MaxCustomers
	license.MaxBranches = cmd.MaxBranches
	license.ExpiresAt = cmd.ExpiresAt
	license.FlashERPURL = cmd.FlashERPURL
	license.ProvisionAPIKey = cmd.ProvisionAPIKey

	if err := h.licenseRepo.Save(ctx, license); err != nil {
		return nil, fmt.Errorf("gagal simpan lisensi: %w", err)
	}

	return &CreateClientLicenseResult{
		LicenseKey:      licenseKey,
		ClientEmail:     cmd.ClientEmail,
		InitialPassword: plainPassword,
	}, nil
}
