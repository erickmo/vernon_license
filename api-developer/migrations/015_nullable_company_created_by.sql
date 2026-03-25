-- +migrate Up
ALTER TABLE companies DROP CONSTRAINT IF EXISTS companies_created_by_fkey;
ALTER TABLE companies ALTER COLUMN created_by DROP NOT NULL;

-- +migrate Down
ALTER TABLE companies ALTER COLUMN created_by SET NOT NULL;
ALTER TABLE companies ADD CONSTRAINT companies_created_by_fkey FOREIGN KEY (created_by) REFERENCES users(id);
