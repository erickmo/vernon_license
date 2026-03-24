-- +migrate Up
CREATE TABLE products (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name              VARCHAR(255) NOT NULL,
    slug              VARCHAR(100) NOT NULL UNIQUE,
    description       TEXT,
    available_modules JSONB NOT NULL DEFAULT '[]',
    available_apps    JSONB NOT NULL DEFAULT '[]',
    available_plans   TEXT[] NOT NULL DEFAULT '{"saas","dedicated"}',
    base_pricing      JSONB NOT NULL DEFAULT '{}',
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_products_slug ON products(slug) WHERE deleted_at IS NULL;

-- +migrate Down
DROP TABLE IF EXISTS products;
