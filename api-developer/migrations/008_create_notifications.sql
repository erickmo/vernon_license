-- +migrate Up
CREATE TABLE notifications (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id),
    type       VARCHAR(50) NOT NULL,
    title      VARCHAR(255) NOT NULL,
    body       TEXT NOT NULL,
    data       JSONB,
    is_read    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notif_user_unread ON notifications(user_id, is_read) WHERE NOT is_read;

-- +migrate Down
DROP TABLE IF EXISTS notifications;
