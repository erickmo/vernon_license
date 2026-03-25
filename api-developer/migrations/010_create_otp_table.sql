-- +migrate Up
CREATE TABLE otp (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_otp_expires ON otp(expires_at);

-- +migrate Down
DROP INDEX IF EXISTS idx_otp_expires;
DROP TABLE IF EXISTS otp;
