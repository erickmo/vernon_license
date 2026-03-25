-- +migrate Up
ALTER TABLE client_licenses DROP CONSTRAINT IF EXISTS client_licenses_created_by_fkey;
ALTER TABLE client_licenses ALTER COLUMN created_by DROP NOT NULL;

-- +migrate Down
ALTER TABLE client_licenses ALTER COLUMN created_by SET NOT NULL;
ALTER TABLE client_licenses ADD CONSTRAINT client_licenses_created_by_fkey FOREIGN KEY (created_by) REFERENCES users(id);
