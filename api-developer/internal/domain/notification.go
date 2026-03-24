package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Notification adalah notifikasi in-app untuk user Vernon.
type Notification struct {
	ID        uuid.UUID       `db:"id"`
	UserID    uuid.UUID       `db:"user_id"`
	Type      string          `db:"type"`
	Title     string          `db:"title"`
	Body      string          `db:"body"`
	Data      json.RawMessage `db:"data"`
	IsRead    bool            `db:"is_read"`
	CreatedAt time.Time       `db:"created_at"`
}

// NotificationRepository mendefinisikan operasi persistence untuk entitas Notification.
type NotificationRepository interface {
	// Create menyimpan notifikasi baru ke database.
	Create(ctx context.Context, n *Notification) error

	// FindByUser mengembalikan notifikasi untuk user tertentu, dibatasi limit.
	FindByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*Notification, error)

	// CountUnread mengembalikan jumlah notifikasi yang belum dibaca oleh user.
	CountUnread(ctx context.Context, userID uuid.UUID) (int, error)

	// MarkRead menandai satu notifikasi sebagai sudah dibaca.
	MarkRead(ctx context.Context, id uuid.UUID) error

	// MarkAllRead menandai semua notifikasi milik user sebagai sudah dibaca.
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
}
