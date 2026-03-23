-- +migrate Up
CREATE TABLE client_licenses (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_key          VARCHAR(64) UNIQUE NOT NULL,
    client_name          VARCHAR(200) NOT NULL,
    client_email         VARCHAR(200) NOT NULL,
    product              VARCHAR(50) NOT NULL DEFAULT 'flasherp',
    plan                 VARCHAR(20) NOT NULL DEFAULT 'saas',
    status               VARCHAR(20) NOT NULL DEFAULT 'active',
    max_users            INTEGER NULL,
    max_trans_per_month  INTEGER NULL,
    max_trans_per_day    INTEGER NULL,
    max_items            INTEGER NULL,
    max_customers        INTEGER NULL,
    max_branches         INTEGER NULL,
    expires_at           TIMESTAMPTZ NULL,
    flasherp_url         TEXT NULL,
    provision_api_key    TEXT NULL,
    is_provisioned       BOOLEAN NOT NULL DEFAULT false,
    created_by           UUID REFERENCES users(id),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_client_licenses_key     ON client_licenses(license_key);
CREATE INDEX idx_client_licenses_status  ON client_licenses(status);
CREATE INDEX idx_client_licenses_product ON client_licenses(product);

-- +migrate Down
DROP TABLE IF EXISTS client_licenses;
