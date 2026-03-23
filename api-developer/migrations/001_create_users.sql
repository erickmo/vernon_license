-- +migrate Up
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(200) NOT NULL,
    email         VARCHAR(200) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(50) NOT NULL DEFAULT 'developer_sales',
    status        VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- Seed initial superuser (password: FlashAdmin2024!)
-- Jalankan perintah berikut untuk insert superuser setelah setup:
-- INSERT INTO users (name, email, password_hash, role, status)
-- VALUES ('Super Admin', 'admin@flashlab.id', '<bcrypt_hash>', 'superuser', 'active');
-- Buat hash dengan: htpasswd -bnBC 12 "" FlashAdmin2024! | tr -d ':\n'

-- +migrate Down
DROP TABLE IF EXISTS users;
