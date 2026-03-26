-- +migrate Up
ALTER TABLE client_licenses ADD COLUMN client_app_ip VARCHAR(255);
ALTER TABLE client_licenses ADD COLUMN superuser_username VARCHAR(255);

-- +migrate Down
ALTER TABLE client_licenses DROP COLUMN IF EXISTS client_app_ip;
ALTER TABLE client_licenses DROP COLUMN IF EXISTS superuser_username;
