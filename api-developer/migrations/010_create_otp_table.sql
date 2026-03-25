-- +migrate Up
CREATE TABLE otp (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ
);

-- Index untuk cepat mencari OTP aktif (tidak expired, belum digunakan)
CREATE INDEX idx_otp_active ON otp(used_at, expires_at) WHERE used_at IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS idx_otp_active;
DROP TABLE IF EXISTS otp;
