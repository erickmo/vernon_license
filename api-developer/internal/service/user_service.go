package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/flashlab/vernon-license/internal/domain"
)

// validRoles berisi daftar role yang diperbolehkan dalam sistem Vernon.
var validRoles = map[string]bool{
	"superuser":     true,
	"project_owner": true,
	"sales":         true,
}

// CreateUserRequest berisi data yang dibutuhkan untuk membuat user baru.
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // "superuser" | "project_owner" | "sales"
}

// UserService menyediakan business logic untuk manajemen user.
type UserService struct {
	repo      domain.UserRepository
	auditRepo domain.AuditLogRepository
	logger    *zap.Logger
}

// NewUserService membuat instance UserService baru.
func NewUserService(repo domain.UserRepository, auditRepo domain.AuditLogRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo:      repo,
		auditRepo: auditRepo,
		logger:    logger,
	}
}

// GetByID mengambil user berdasarkan ID.
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetByID: %w", err)
	}
	return user, nil
}

// List mengembalikan semua user. Validasi role dilakukan oleh caller.
func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserService.List: %w", err)
	}
	return users, nil
}

// Create membuat user baru.
// Password di-hash dengan bcrypt cost 12.
// Mengembalikan ErrUserEmailExists jika email sudah digunakan.
// Mengembalikan ErrValidationFailed jika role tidak valid.
// Audit log dibuat dengan action "user_created" setelah operasi berhasil.
func (s *UserService) Create(ctx context.Context, req CreateUserRequest, actorID uuid.UUID, actorName string) (*domain.User, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("UserService.Create: %w", domain.ErrValidationFailed)
	}

	if !validRoles[req.Role] {
		return nil, fmt.Errorf("UserService.Create: %w", domain.ErrValidationFailed)
	}

	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, fmt.Errorf("UserService.Create: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("UserService.Create: %w", domain.ErrUserEmailExists)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("UserService.Create: failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("UserService.Create: %w", err)
	}

	changes, _ := json.Marshal(map[string]any{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
	LogAudit(ctx, s.auditRepo, s.logger, "user", user.ID, "user_created", actorID, actorName, changes)

	return user, nil
}

// createWithoutAudit membuat user baru tanpa audit log.
// Digunakan oleh SetupService saat membuat superuser pertama.
func (s *UserService) createWithoutAudit(ctx context.Context, req CreateUserRequest) (*domain.User, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("UserService.createWithoutAudit: %w", domain.ErrValidationFailed)
	}

	if !validRoles[req.Role] {
		return nil, fmt.Errorf("UserService.createWithoutAudit: %w", domain.ErrValidationFailed)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("UserService.createWithoutAudit: failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("UserService.createWithoutAudit: %w", err)
	}

	return user, nil
}

// SetActive mengaktifkan atau menonaktifkan user.
// Audit log dibuat dengan action "user_activated" atau "user_deactivated".
func (s *UserService) SetActive(ctx context.Context, id uuid.UUID, active bool, actorID uuid.UUID, actorName string) error {
	if err := s.repo.SetActive(ctx, id, active); err != nil {
		return fmt.Errorf("UserService.SetActive: %w", err)
	}

	action := "user_deactivated"
	if active {
		action = "user_activated"
	}
	changes, _ := json.Marshal(map[string]any{"is_active": active})
	LogAudit(ctx, s.auditRepo, s.logger, "user", id, action, actorID, actorName, changes)

	return nil
}
