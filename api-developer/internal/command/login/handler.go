package login

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/flashlab/flasherp-developer-api/internal/domain/user"
	jwtpkg "github.com/flashlab/flasherp-developer-api/pkg/jwt"
)

type Login struct {
	Identifier string
	Password   string
}

type LoginResult struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}

type Handler struct {
	userRepo   user.WriteRepository
	jwtService *jwtpkg.Service
}

func NewHandler(userRepo user.WriteRepository, jwtService *jwtpkg.Service) *Handler {
	return &Handler{userRepo: userRepo, jwtService: jwtService}
}

func (h *Handler) Handle(ctx context.Context, cmd any) error {
	return fmt.Errorf("gunakan HandleLogin untuk mendapatkan result")
}

func (h *Handler) HandleLogin(ctx context.Context, cmd Login) (*LoginResult, error) {
	u, err := h.userRepo.GetByEmail(ctx, cmd.Identifier)
	if err != nil {
		return nil, user.ErrInvalidCredentials
	}
	if u.Status != user.StatusActive {
		return nil, user.ErrInactiveAccount
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(cmd.Password)); err != nil {
		return nil, user.ErrInvalidCredentials
	}
	token, err := h.jwtService.GenerateToken(u.ID.String(), u.Email, u.Role)
	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %w", err)
	}
	return &LoginResult{
		AccessToken: token,
		UserID:      u.ID.String(),
		Name:        u.Name,
		Email:       u.Email,
		Role:        u.Role,
	}, nil
}
