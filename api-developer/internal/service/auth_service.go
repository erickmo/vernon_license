package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/flashlab/vernon-license/internal/config"
	"github.com/flashlab/vernon-license/internal/domain"
	jwtpkg "github.com/flashlab/vernon-license/pkg/jwt"
)

// AuthService menyediakan login dan token management untuk Vernon App.
// BUKAN untuk public API.
type AuthService struct {
	userRepo  domain.UserRepository
	auditRepo domain.AuditLogRepository
	jwtSecret string
	logger    *zap.Logger
}

// NewAuthService membuat instance AuthService baru.
func NewAuthService(userRepo domain.UserRepository, auditRepo domain.AuditLogRepository, cfg *config.Config, logger *zap.Logger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		auditRepo: auditRepo,
		jwtSecret: cfg.JWTSecret,
		logger:    logger,
	}
}

// Login memvalidasi email dan password lalu mengembalikan JWT token.
// Mengembalikan ErrAuthInvalidCredentials jika email atau password salah, atau user tidak aktif.
// Audit log dibuat dengan action "user_login" setelah login berhasil.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("AuthService.Login: %w", domain.ErrAuthInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("AuthService.Login: %w", domain.ErrAuthInvalidCredentials)
	}

	if !user.IsActive {
		return "", nil, fmt.Errorf("AuthService.Login: %w", domain.ErrAuthInvalidCredentials)
	}

	token, err := jwtpkg.Sign(user.ID.String(), user.Name, user.Role, s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("AuthService.Login: failed to sign token: %w", err)
	}

	meta, _ := json.Marshal(map[string]any{"email": user.Email})
	LogAudit(ctx, s.auditRepo, s.logger, "user", user.ID, "user_login", user.ID, user.Name, meta)

	return token, user, nil
}

// signToken menghasilkan JWT untuk user yang diberikan.
// Digunakan secara internal oleh SetupService.
func (s *AuthService) signToken(user *domain.User) (string, error) {
	return jwtpkg.Sign(user.ID.String(), user.Name, user.Role, s.jwtSecret)
}

// GetMe mengambil current user berdasarkan user ID dari JWT.
func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("AuthService.GetMe: %w", err)
	}
	return user, nil
}
