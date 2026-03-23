package get_setup_status

import (
	"context"
	"fmt"

	userdomain "github.com/flashlab/flasherp-developer-api/internal/domain/user"
)

type GetSetupStatus struct{}

type SetupStatusResult struct {
	IsInstalled bool `json:"is_installed"`
}

type Handler struct {
	userReadRepo userdomain.ReadRepository
}

func NewHandler(userReadRepo userdomain.ReadRepository) *Handler {
	return &Handler{userReadRepo: userReadRepo}
}

func (h *Handler) Handle(ctx context.Context, q any) (any, error) {
	if _, ok := q.(GetSetupStatus); !ok {
		return nil, fmt.Errorf("query tidak valid")
	}

	count, err := h.userReadRepo.CountByRole(ctx, userdomain.RoleSuperuser)
	if err != nil {
		return nil, fmt.Errorf("gagal cek status: %w", err)
	}

	return &SetupStatusResult{IsInstalled: count > 0}, nil
}
