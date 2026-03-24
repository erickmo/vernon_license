-- +migrate Up
CREATE TABLE companies (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255),
    phone      VARCHAR(50),
    address    TEXT,
    pic_name   VARCHAR(255),
    pic_email  VARCHAR(255),
    pic_phone  VARCHAR(50),
    notes      TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_companies_name ON companies(name) WHERE deleted_at IS NULL;

-- +migrate Down
DROP TABLE IF EXISTS companies;
