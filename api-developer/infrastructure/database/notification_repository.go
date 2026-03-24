package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/flashlab/vernon-license/internal/domain"
)

// NotificationRepo adalah implementasi domain.NotificationRepository berbasis PostgreSQL.
type NotificationRepo struct {
	db *sqlx.DB
}

// NewNotificationRepo membuat instance NotificationRepo baru.
func NewNotificationRepo(db *sqlx.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

// Create menyimpan notifikasi baru ke database.
func (r *NotificationRepo) Create(ctx context.Context, n *domain.Notification) error {
	const q = `
		INSERT INTO notifications (id, user_id, type, title, body, data, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING created_at`
	if err := r.db.QueryRowContext(ctx, q,
		n.ID, n.UserID, n.Type, n.Title, n.Body, n.Data, n.IsRead,
	).Scan(&n.CreatedAt); err != nil {
		return fmt.Errorf("NotificationRepo.Create: %w", err)
	}
	return nil
}

// FindByUser mengembalikan notifikasi untuk user tertentu, dibatasi limit, diurutkan terbaru.
func (r *NotificationRepo) FindByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Notification, error) {
	const q = `
		SELECT id, user_id, type, title, body, data, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`
	var notifications []*domain.Notification
	if err := r.db.SelectContext(ctx, &notifications, q, userID, limit); err != nil {
		return nil, fmt.Errorf("NotificationRepo.FindByUser: %w", err)
	}
	return notifications, nil
}

// CountUnread mengembalikan jumlah notifikasi yang belum dibaca oleh user.
func (r *NotificationRepo) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	const q = `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`
	var count int
	if err := r.db.QueryRowContext(ctx, q, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("NotificationRepo.CountUnread: %w", err)
	}
	return count, nil
}

// MarkRead menandai satu notifikasi sebagai sudah dibaca.
func (r *NotificationRepo) MarkRead(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE notifications SET is_read = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("NotificationRepo.MarkRead: %w", err)
	}
	return nil
}

// MarkAllRead menandai semua notifikasi milik user sebagai sudah dibaca.
func (r *NotificationRepo) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	const q = `UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`
	_, err := r.db.ExecContext(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("NotificationRepo.MarkAllRead: %w", err)
	}
	return nil
}
