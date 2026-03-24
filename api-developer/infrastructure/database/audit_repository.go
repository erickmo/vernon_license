package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/flashlab/vernon-license/internal/domain"
)

// AuditRepo adalah implementasi domain.AuditLogRepository berbasis PostgreSQL.
type AuditRepo struct {
	db *sqlx.DB
}

// NewAuditRepo membuat instance AuditRepo baru.
func NewAuditRepo(db *sqlx.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

// Create menyimpan audit log baru ke database.
func (r *AuditRepo) Create(ctx context.Context, log *domain.AuditLog) error {
	const q = `
		INSERT INTO audit_logs
		    (id, entity_type, entity_id, action, actor_id, actor_name, changes, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING created_at`
	if err := r.db.QueryRowContext(ctx, q,
		log.ID, log.EntityType, log.EntityID, log.Action,
		log.ActorID, log.ActorName, log.Changes, log.Metadata,
	).Scan(&log.CreatedAt); err != nil {
		return fmt.Errorf("AuditRepo.Create: %w", err)
	}
	return nil
}

// FindByEntity mengembalikan semua audit log untuk entity tertentu, diurutkan terbaru.
func (r *AuditRepo) FindByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.AuditLog, error) {
	const q = `
		SELECT id, entity_type, entity_id, action, actor_id, actor_name, changes, metadata, created_at
		FROM audit_logs
		WHERE entity_type = $1 AND entity_id = $2
		ORDER BY created_at DESC`
	var logs []*domain.AuditLog
	if err := r.db.SelectContext(ctx, &logs, q, entityType, entityID); err != nil {
		return nil, fmt.Errorf("AuditRepo.FindByEntity: %w", err)
	}
	return logs, nil
}

// FindAll mengembalikan audit log dengan pagination, diurutkan terbaru.
func (r *AuditRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.AuditLog, error) {
	const q = `
		SELECT id, entity_type, entity_id, action, actor_id, actor_name, changes, metadata, created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	var logs []*domain.AuditLog
	if err := r.db.SelectContext(ctx, &logs, q, limit, offset); err != nil {
		return nil, fmt.Errorf("AuditRepo.FindAll: %w", err)
	}
	return logs, nil
}
