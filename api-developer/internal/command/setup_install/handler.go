package setup_install

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	userdomain "github.com/flashlab/flasherp-developer-api/internal/domain/user"
)

var ErrAlreadyInstalled = errors.New("sistem sudah terinstal")

type SetupInstall struct {
	Name     string
	Email    string
	Password string
}

type Handler struct {
	userWriteRepo userdomain.WriteRepository
	userReadRepo  userdomain.ReadRepository
}

func NewHandler(writeRepo userdomain.WriteRepository, readRepo userdomain.ReadRepository) *Handler {
	return &Handler{userWriteRepo: writeRepo, userReadRepo: readRepo}
}

func (h *Handler) Handle(ctx context.Context, cmd any) error {
	c, ok := cmd.(SetupInstall)
	if !ok {
		return fmt.Errorf("command tidak valid")
	}

	count, err := h.userReadRepo.CountByRole(ctx, userdomain.RoleSuperuser)
	if err != nil {
		return fmt.Errorf("gagal cek status instalasi: %w", err)
	}
	if count > 0 {
		return ErrAlreadyInstalled
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password: %w", err)
	}

	u, err := userdomain.NewUser(c.Name, c.Email, string(hash), userdomain.RoleSuperuser)
	if err != nil {
		return err
	}

	return h.userWriteRepo.Save(ctx, u)
}
