package update_license_constraints

import (
	"context"
	"fmt"
	"time"

	clientlicense "github.com/flashlab/flasherp-developer-api/internal/domain/client_license"
)

type UpdateLicenseConstraints struct {
	LicenseKey       string
	MaxUsers         *int
	MaxTransPerMonth *int
	MaxTransPerDay   *int
	MaxItems         *int
	MaxCustomers     *int
	MaxBranches      *int
	ExpiresAt        *time.Time
}

type Handler struct {
	licenseRepo clientlicense.WriteRepository
}

func NewHandler(licenseRepo clientlicense.WriteRepository) *Handler {
	return &Handler{licenseRepo: licenseRepo}
}

func (h *Handler) Handle(ctx context.Context, cmd any) error {
	c, ok := cmd.(UpdateLicenseConstraints)
	if !ok {
		return fmt.Errorf("perintah tidak valid")
	}
	license, err := h.licenseRepo.GetByKey(ctx, c.LicenseKey)
	if err != nil {
		return err
	}
	license.MaxUsers = c.MaxUsers
	license.MaxTransPerMonth = c.MaxTransPerMonth
	license.MaxTransPerDay = c.MaxTransPerDay
	license.MaxItems = c.MaxItems
	license.MaxCustomers = c.MaxCustomers
	license.MaxBranches = c.MaxBranches
	license.ExpiresAt = c.ExpiresAt
	license.UpdatedAt = time.Now().UTC()
	return h.licenseRepo.Update(ctx, license)
}
