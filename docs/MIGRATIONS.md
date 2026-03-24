# Database Migrations

Format: sql-migrate. PostgreSQL port: **5433**.

## 001_create_users.sql
```sql
-- +migrate Up
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role          VARCHAR(50) NOT NULL DEFAULT 'sales',
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +migrate Down
DROP TABLE IF EXISTS users;
```

## 002_create_companies.sql
```sql
-- +migrate Up
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL, email VARCHAR(255), phone VARCHAR(50), address TEXT,
    pic_name VARCHAR(255), pic_email VARCHAR(255), pic_phone VARCHAR(50), notes TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_companies_name ON companies(name) WHERE deleted_at IS NULL;
-- +migrate Down
DROP TABLE IF EXISTS companies;
```

## 003_create_projects.sql
```sql
-- +migrate Up
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL, description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_projects_company ON projects(company_id) WHERE deleted_at IS NULL;
-- +migrate Down
DROP TABLE IF EXISTS projects;
```

## 004_create_products.sql
```sql
-- +migrate Up
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL, slug VARCHAR(100) NOT NULL UNIQUE, description TEXT,
    available_modules JSONB NOT NULL DEFAULT '[]',
    available_apps JSONB NOT NULL DEFAULT '[]',
    available_plans TEXT[] NOT NULL DEFAULT '{"saas","dedicated"}',
    base_pricing JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_products_slug ON products(slug) WHERE deleted_at IS NULL;
-- +migrate Down
DROP TABLE IF EXISTS products;
```

## 005_create_client_licenses.sql
```sql
-- +migrate Up
CREATE TABLE client_licenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_key VARCHAR(20) NOT NULL UNIQUE,
    project_id UUID NOT NULL REFERENCES projects(id),
    company_id UUID NOT NULL REFERENCES companies(id),
    product_id UUID NOT NULL REFERENCES products(id),
    plan VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    modules TEXT[] NOT NULL DEFAULT '{}',
    apps TEXT[] NOT NULL DEFAULT '{}',
    contract_amount DECIMAL(15,2),
    description TEXT,
    max_users INT, max_trans_per_month INT, max_trans_per_day INT,
    max_items INT, max_customers INT, max_branches INT, max_storage INT,
    expires_at TIMESTAMPTZ,
    instance_url TEXT,
    instance_name TEXT,
    provision_api_key TEXT,
    check_interval VARCHAR(10) NOT NULL DEFAULT '6h',
    last_pull_at TIMESTAMPTZ,
    is_registered BOOLEAN NOT NULL DEFAULT FALSE,
    proposal_id UUID,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ, archived_at TIMESTAMPTZ
);
CREATE INDEX idx_licenses_key ON client_licenses(license_key);
CREATE INDEX idx_licenses_project ON client_licenses(project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_licenses_status ON client_licenses(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_licenses_provision_key ON client_licenses(provision_api_key) WHERE provision_api_key IS NOT NULL;
-- +migrate Down
DROP TABLE IF EXISTS client_licenses;
```

## 006_create_proposals.sql
```sql
-- +migrate Up
CREATE TABLE proposals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id),
    company_id UUID NOT NULL REFERENCES companies(id),
    product_id UUID NOT NULL REFERENCES products(id),
    version INT NOT NULL, status VARCHAR(20) NOT NULL DEFAULT 'draft',
    modules TEXT[] NOT NULL DEFAULT '{}', apps TEXT[] NOT NULL DEFAULT '{}',
    plan VARCHAR(50) NOT NULL,
    max_users INT, max_trans_per_month INT, max_trans_per_day INT,
    max_items INT, max_customers INT, max_branches INT, max_storage INT,
    contract_amount DECIMAL(15,2), expires_at TIMESTAMPTZ,
    notes TEXT, owner_notes TEXT, rejection_reason TEXT,
    changelog JSONB, pdf_path TEXT, pdf_generated_at TIMESTAMPTZ,
    submitted_by UUID NOT NULL REFERENCES users(id),
    reviewed_by UUID REFERENCES users(id), reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id, product_id, version)
);
CREATE INDEX idx_proposals_project ON proposals(project_id, version DESC);
ALTER TABLE client_licenses ADD CONSTRAINT fk_license_proposal FOREIGN KEY (proposal_id) REFERENCES proposals(id);
-- +migrate Down
ALTER TABLE client_licenses DROP CONSTRAINT IF EXISTS fk_license_proposal;
DROP TABLE IF EXISTS proposals;
```

## 007_create_audit_logs.sql
```sql
-- +migrate Up
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL, entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    actor_id UUID NOT NULL REFERENCES users(id), actor_name VARCHAR(255) NOT NULL,
    changes JSONB, metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_time ON audit_logs(created_at DESC);
-- +migrate Down
DROP TABLE IF EXISTS audit_logs;
```

## 008_create_notifications.sql
```sql
-- +migrate Up
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL, title VARCHAR(255) NOT NULL, body TEXT NOT NULL,
    data JSONB, is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notif_user_unread ON notifications(user_id, is_read) WHERE NOT is_read;
-- +migrate Down
DROP TABLE IF EXISTS notifications;
```

## Schema Notes
- `provision_api_key` indexed for register lookup, never exposed via public API response
- `is_registered` = client app has called register endpoint
- `last_pull_at` updated on every validate call (monitoring)
- `status: pending` = created but client hasn't registered or PO hasn't approved yet
- `modules`, `apps` are tracked internally — NOT returned in validate response
