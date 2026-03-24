package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
)

// NotificationService menyediakan manajemen notifikasi untuk Vernon App.
type NotificationService struct {
	repo   domain.NotificationRepository
	logger *zap.Logger
}

// NewNotificationService membuat instance NotificationService baru.
func NewNotificationService(repo domain.NotificationRepository, logger *zap.Logger) *NotificationService {
	return &NotificationService{
		repo:   repo,
		logger: logger,
	}
}

// Send membuat notifikasi baru untuk user.
func (s *NotificationService) Send(ctx context.Context, userID uuid.UUID, notifType, title, body string, data json.RawMessage) error {
	if data == nil {
		data = json.RawMessage("{}")
	}

	n := &domain.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Body:      body,
		Data:      data,
		IsRead:    false,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return fmt.Errorf("NotificationService.Send: %w", err)
	}
	return nil
}

// ListForUser mengambil notifikasi user (paling baru dulu), limit 50.
func (s *NotificationService) ListForUser(ctx context.Context, userID uuid.UUID) ([]*domain.Notification, error) {
	notifications, err := s.repo.FindByUser(ctx, userID, 50)
	if err != nil {
		return nil, fmt.Errorf("NotificationService.ListForUser: %w", err)
	}
	return notifications, nil
}

// CountUnread menghitung unread notifications untuk badge.
func (s *NotificationService) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := s.repo.CountUnread(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("NotificationService.CountUnread: %w", err)
	}
	return count, nil
}

// MarkRead menandai satu notifikasi sebagai sudah dibaca.
func (s *NotificationService) MarkRead(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.MarkRead(ctx, id); err != nil {
		return fmt.Errorf("NotificationService.MarkRead: %w", err)
	}
	return nil
}

// MarkAllRead menandai semua notifikasi user sebagai sudah dibaca.
func (s *NotificationService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.MarkAllRead(ctx, userID); err != nil {
		return fmt.Errorf("NotificationService.MarkAllRead: %w", err)
	}
	return nil
}
