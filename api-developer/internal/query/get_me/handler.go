package get_me

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/flashlab/flasherp-developer-api/internal/domain/user"
)

type GetMe struct {
	UserID string
}

type Handler struct {
	userRepo user.ReadRepository
}

func NewHandler(userRepo user.ReadRepository) *Handler {
	return &Handler{userRepo: userRepo}
}

func (h *Handler) Handle(ctx context.Context, q any) (any, error) {
	query, ok := q.(GetMe)
	if !ok {
		return nil, fmt.Errorf("query tidak valid")
	}

	id, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, fmt.Errorf("user ID tidak valid: %w", err)
	}

	u, err := h.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}
