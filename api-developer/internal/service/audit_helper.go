// Package service berisi business logic layer untuk Vernon License.
package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
)

// LogAudit membuat audit log entry secara fire-and-forget.
// Jika operasi gagal, error hanya dicatat ke logger dan tidak dipropagasi ke caller.
func LogAudit(
	ctx context.Context,
	repo domain.AuditLogRepository,
	logger *zap.Logger,
	entityType string,
	entityID uuid.UUID,
	action string,
	actorID uuid.UUID,
	actorName string,
	changes json.RawMessage,
) {
	entry := &domain.AuditLog{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		ActorID:    actorID,
		ActorName:  actorName,
		Changes:    changes,
		Metadata:   json.RawMessage("{}"),
	}
	if err := repo.Create(ctx, entry); err != nil {
		logger.Error("audit log failed",
			zap.String("entity_type", entityType),
			zap.String("entity_id", entityID.String()),
			zap.String("action", action),
			zap.Error(err),
		)
	}
}
