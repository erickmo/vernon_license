package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
)

// SetupService menangani first-run setup (create superuser pertama).
type SetupService struct {
	userRepo domain.UserRepository
	authSvc  *AuthService
	logger   *zap.Logger
}

// NewSetupService membuat instance SetupService baru.
func NewSetupService(userRepo domain.UserRepository, authSvc *AuthService, logger *zap.Logger) *SetupService {
	return &SetupService{
		userRepo: userRepo,
		authSvc:  authSvc,
		logger:   logger,
	}
}

// IsSetup mengecek apakah sudah ada user di database.
// Mengembalikan true jika sudah ada minimal satu user.
func (s *SetupService) IsSetup(ctx context.Context) (bool, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return false, fmt.Errorf("SetupService.IsSetup: %w", err)
	}
	return len(users) > 0, nil
}

// Install membuat superuser pertama.
// Hanya bisa dipanggil jika belum ada user di database.
// Mengembalikan JWT token dan user yang dibuat.
func (s *SetupService) Install(ctx context.Context, name, email, password string) (string, *domain.User, error) {
	isSetup, err := s.IsSetup(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("SetupService.Install: %w", err)
	}
	if isSetup {
		return "", nil, fmt.Errorf("SetupService.Install: %w", domain.ErrValidationFailed)
	}

	// Gunakan UserService sementara tanpa auditRepo (setup pertama kali)
	userSvc := &UserService{
		repo:      s.userRepo,
		auditRepo: nil,
		logger:    s.logger,
	}

	req := CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     "superuser",
	}

	user, err := userSvc.createWithoutAudit(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("SetupService.Install: %w", err)
	}

	token, err := s.authSvc.signToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("SetupService.Install: %w", err)
	}

	s.logger.Info("setup completed: superuser created",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	return token, user, nil
}
