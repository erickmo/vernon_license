package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditLog merekam setiap aksi penting yang dilakukan di sistem Vernon.
type AuditLog struct {
	ID         uuid.UUID       `db:"id"`
	EntityType string          `db:"entity_type"`
	EntityID   uuid.UUID       `db:"entity_id"`
	Action     string          `db:"action"`
	ActorID    uuid.UUID       `db:"actor_id"`
	ActorName  string          `db:"actor_name"`
	Changes    json.RawMessage `db:"changes"`
	Metadata   json.RawMessage `db:"metadata"`
	CreatedAt  time.Time       `db:"created_at"`
}

// AuditLogRepository mendefinisikan operasi persistence untuk entitas AuditLog.
type AuditLogRepository interface {
	// Create menyimpan audit log baru ke database.
	Create(ctx context.Context, log *AuditLog) error

	// FindByEntity mengembalikan semua audit log untuk entity tertentu.
	FindByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*AuditLog, error)

	// FindAll mengembalikan audit log dengan pagination.
	FindAll(ctx context.Context, limit, offset int) ([]*AuditLog, error)
}
