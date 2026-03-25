-- +migrate Up
CREATE UNIQUE INDEX idx_licenses_company_product_active
    ON client_licenses (company_id, product_id)
    WHERE deleted_at IS NULL AND company_id IS NOT NULL;

-- +migrate Down
DROP INDEX IF EXISTS idx_licenses_company_product_active;
