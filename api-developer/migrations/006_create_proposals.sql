-- +migrate Up
CREATE TABLE proposals (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id          UUID NOT NULL REFERENCES projects(id),
    company_id          UUID NOT NULL REFERENCES companies(id),
    product_id          UUID NOT NULL REFERENCES products(id),
    version             INT NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'draft',
    modules             TEXT[] NOT NULL DEFAULT '{}',
    apps                TEXT[] NOT NULL DEFAULT '{}',
    plan                VARCHAR(50) NOT NULL,
    max_users           INT,
    max_trans_per_month INT,
    max_trans_per_day   INT,
    max_items           INT,
    max_customers       INT,
    max_branches        INT,
    max_storage         INT,
    contract_amount     DECIMAL(15,2),
    expires_at          TIMESTAMPTZ,
    notes               TEXT,
    owner_notes         TEXT,
    rejection_reason    TEXT,
    changelog           JSONB,
    pdf_path            TEXT,
    pdf_generated_at    TIMESTAMPTZ,
    submitted_by        UUID NOT NULL REFERENCES users(id),
    reviewed_by         UUID REFERENCES users(id),
    reviewed_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id, product_id, version)
);

CREATE INDEX idx_proposals_project ON proposals(project_id, version DESC);

ALTER TABLE client_licenses
    ADD CONSTRAINT fk_license_proposal
    FOREIGN KEY (proposal_id) REFERENCES proposals(id);

-- +migrate Down
ALTER TABLE client_licenses DROP CONSTRAINT IF EXISTS fk_license_proposal;
DROP TABLE IF EXISTS proposals;
