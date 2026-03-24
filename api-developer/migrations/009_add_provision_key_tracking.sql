-- +migrate Up
ALTER TABLE client_licenses
ADD COLUMN provision_api_key_generated_at TIMESTAMPTZ,
ADD COLUMN provision_api_key_previous TEXT,
ADD COLUMN provision_api_key_previous_at TIMESTAMPTZ;

-- Set existing records: assume they were generated now for grace period calculation
UPDATE client_licenses
SET provision_api_key_generated_at = created_at
WHERE provision_api_key IS NOT NULL AND provision_api_key_generated_at IS NULL;

-- Index untuk query provision key dengan grace period check
CREATE INDEX idx_licenses_provision_key_time ON client_licenses(provision_api_key, provision_api_key_generated_at) WHERE provision_api_key IS NOT NULL;

-- +migrate Down
DROP INDEX IF EXISTS idx_licenses_provision_key_time;
ALTER TABLE client_licenses
DROP COLUMN provision_api_key_previous_at,
DROP COLUMN provision_api_key_previous,
DROP COLUMN provision_api_key_generated_at;
