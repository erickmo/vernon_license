//go:build !wasm

package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/flashlab/vernon-license/internal/domain"
	appmiddleware "github.com/flashlab/vernon-license/internal/middleware"
	"github.com/flashlab/vernon-license/internal/service"
)

// NotificationHandler menangani HTTP requests untuk Notifications.
type NotificationHandler struct {
	svc    *service.NotificationService
	logger *zap.Logger
}

// NewNotificationHandler membuat instance NotificationHandler baru.
func NewNotificationHandler(svc *service.NotificationService, logger *zap.Logger) *NotificationHandler {
	return &NotificationHandler{svc: svc, logger: logger}
}

// notificationDTO adalah representasi Notification yang dikembalikan ke client.
type notificationDTO struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

// toNotificationDTO mengkonversi domain.Notification ke notificationDTO.
func toNotificationDTO(n *domain.Notification) notificationDTO {
	return notificationDTO{
		ID:        n.ID.String(),
		Type:      n.Type,
		Title:     n.Title,
		Body:      n.Body,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// List menangani GET /api/internal/notifications.
// Mengembalikan 50 notifikasi terbaru milik current user.
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	userID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	notifications, err := h.svc.ListForUser(r.Context(), userID)
	if err != nil {
		h.logger.Error("NotificationHandler.List", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	dtos := make([]notificationDTO, 0, len(notifications))
	for _, n := range notifications {
		dtos = append(dtos, toNotificationDTO(n))
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": dtos})
}

// CountUnread menangani GET /api/internal/notifications/unread-count.
// Mengembalikan jumlah notifikasi yang belum dibaca oleh current user.
func (h *NotificationHandler) CountUnread(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	userID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	count, err := h.svc.CountUnread(r.Context(), userID)
	if err != nil {
		h.logger.Error("NotificationHandler.CountUnread", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"count": count})
}

// MarkRead menangani PUT /api/internal/notifications/{id}/read.
// Menandai satu notifikasi sebagai sudah dibaca.
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	_, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "ID tidak valid")
		return
	}

	if err := h.svc.MarkRead(r.Context(), id); err != nil {
		h.logger.Error("NotificationHandler.MarkRead", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Notifikasi ditandai sudah dibaca"})
}

// MarkAllRead menangani PUT /api/internal/notifications/read-all.
// Menandai semua notifikasi current user sebagai sudah dibaca.
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	claims, ok := appmiddleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	userID, err := parseUUID(claims.Sub)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token subject")
		return
	}

	if err := h.svc.MarkAllRead(r.Context(), userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "USER_NOT_FOUND", "User tidak ditemukan")
			return
		}
		h.logger.Error("NotificationHandler.MarkAllRead", zap.Error(err))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Semua notifikasi ditandai sudah dibaca"})
}
