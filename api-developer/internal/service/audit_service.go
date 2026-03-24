package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
)

// AuditService menyediakan logging untuk semua perubahan entitas.
type AuditService struct {
	repo   domain.AuditLogRepository
	logger *zap.Logger
}

// NewAuditService membuat instance AuditService baru.
func NewAuditService(repo domain.AuditLogRepository, logger *zap.Logger) *AuditService {
	return &AuditService{
		repo:   repo,
		logger: logger,
	}
}

// Log mencatat audit event. Non-blocking — error hanya di-log, tidak di-return.
func (s *AuditService) Log(
	ctx context.Context,
	entityType string,
	entityID uuid.UUID,
	action string,
	actorID uuid.UUID,
	actorName string,
	changes json.RawMessage,
	metadata json.RawMessage,
) {
	if changes == nil {
		changes = json.RawMessage("{}")
	}
	if metadata == nil {
		metadata = json.RawMessage("{}")
	}

	entry := &domain.AuditLog{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		ActorID:    actorID,
		ActorName:  actorName,
		Changes:    changes,
		Metadata:   metadata,
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		s.logger.Error("AuditService.Log: failed to create audit log",
			zap.String("entity_type", entityType),
			zap.String("entity_id", entityID.String()),
			zap.String("action", action),
			zap.Error(err),
		)
	}
}

// ListByEntity mengambil audit logs untuk satu entitas.
func (s *AuditService) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.AuditLog, error) {
	logs, err := s.repo.FindByEntity(ctx, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("AuditService.ListByEntity: %w", err)
	}
	return logs, nil
}

// ListAll mengembalikan audit log global dengan pagination. Untuk superuser.
func (s *AuditService) ListAll(ctx context.Context, limit, offset int) ([]*domain.AuditLog, error) {
	logs, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("AuditService.ListAll: %w", err)
	}
	return logs, nil
}
