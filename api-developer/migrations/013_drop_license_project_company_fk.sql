-- +migrate Up
-- Allow self-registered licenses without project/company assignment.
-- Admin can assign project/company later via the App.
ALTER TABLE client_licenses DROP CONSTRAINT IF EXISTS client_licenses_project_id_fkey;
ALTER TABLE client_licenses DROP CONSTRAINT IF EXISTS client_licenses_company_id_fkey;
ALTER TABLE client_licenses ALTER COLUMN project_id DROP NOT NULL;
ALTER TABLE client_licenses ALTER COLUMN company_id DROP NOT NULL;

-- +migrate Down
ALTER TABLE client_licenses ALTER COLUMN project_id SET NOT NULL;
ALTER TABLE client_licenses ALTER COLUMN company_id SET NOT NULL;
