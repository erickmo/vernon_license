-- +migrate Up
CREATE TABLE client_licenses (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_key          VARCHAR(20) NOT NULL UNIQUE,
    project_id           UUID NOT NULL REFERENCES projects(id),
    company_id           UUID NOT NULL REFERENCES companies(id),
    product_id           UUID NOT NULL REFERENCES products(id),
    plan                 VARCHAR(50) NOT NULL,
    status               VARCHAR(20) NOT NULL DEFAULT 'pending',
    modules              TEXT[] NOT NULL DEFAULT '{}',
    apps                 TEXT[] NOT NULL DEFAULT '{}',
    contract_amount      DECIMAL(15,2),
    description          TEXT,
    max_users            INT,
    max_trans_per_month  INT,
    max_trans_per_day    INT,
    max_items            INT,
    max_customers        INT,
    max_branches         INT,
    max_storage          INT,
    expires_at           TIMESTAMPTZ,
    instance_url         TEXT,
    instance_name        TEXT,
    provision_api_key    TEXT,
    check_interval       VARCHAR(10) NOT NULL DEFAULT '6h',
    last_pull_at         TIMESTAMPTZ,
    is_registered        BOOLEAN NOT NULL DEFAULT FALSE,
    proposal_id          UUID,
    created_by           UUID NOT NULL REFERENCES users(id),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ,
    archived_at          TIMESTAMPTZ
);

CREATE INDEX idx_licenses_key ON client_licenses(license_key);
CREATE INDEX idx_licenses_project ON client_licenses(project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_licenses_status ON client_licenses(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_licenses_provision_key ON client_licenses(provision_api_key) WHERE provision_api_key IS NOT NULL;

-- +migrate Down
DROP TABLE IF EXISTS client_licenses;
