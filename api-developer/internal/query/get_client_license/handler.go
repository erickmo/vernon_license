package get_client_license

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	clientlicense "github.com/flashlab/flasherp-developer-api/internal/domain/client_license"
)

type GetClientLicense struct {
	ID string
}

type Handler struct {
	licenseRepo clientlicense.ReadRepository
}

func NewHandler(licenseRepo clientlicense.ReadRepository) *Handler {
	return &Handler{licenseRepo: licenseRepo}
}

func (h *Handler) Handle(ctx context.Context, q any) (any, error) {
	query, ok := q.(GetClientLicense)
	if !ok {
		return nil, fmt.Errorf("query tidak valid")
	}

	// Coba parse sebagai UUID dulu; jika gagal, anggap sebagai license key
	if id, err := uuid.Parse(query.ID); err == nil {
		license, err := h.licenseRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return license, nil
	}

	license, err := h.licenseRepo.GetByKey(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return license, nil
}
