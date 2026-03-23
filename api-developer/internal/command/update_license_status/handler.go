package update_license_status

import (
	"context"
	"fmt"
	"time"

	clientlicense "github.com/flashlab/flasherp-developer-api/internal/domain/client_license"
)

type UpdateLicenseStatus struct {
	LicenseKey string
	Status     string
}

type Handler struct {
	licenseRepo clientlicense.WriteRepository
}

func NewHandler(licenseRepo clientlicense.WriteRepository) *Handler {
	return &Handler{licenseRepo: licenseRepo}
}

func (h *Handler) Handle(ctx context.Context, cmd any) error {
	c, ok := cmd.(UpdateLicenseStatus)
	if !ok {
		return fmt.Errorf("perintah tidak valid")
	}
	license, err := h.licenseRepo.GetByKey(ctx, c.LicenseKey)
	if err != nil {
		return err
	}
	switch c.Status {
	case clientlicense.StatusActive:
		license.Activate()
	case clientlicense.StatusSuspended:
		license.Suspend()
	case clientlicense.StatusExpired:
		license.Expire()
	default:
		license.Status = c.Status
		license.UpdatedAt = time.Now().UTC()
	}
	return h.licenseRepo.Update(ctx, license)
}
