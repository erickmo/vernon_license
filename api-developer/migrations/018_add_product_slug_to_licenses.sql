-- +migrate Up
ALTER TABLE client_licenses ADD COLUMN product_slug VARCHAR(100);

-- Backfill dari tabel products
UPDATE client_licenses cl
SET product_slug = p.slug
FROM products p
WHERE cl.product_id = p.id;

-- +migrate Down
ALTER TABLE client_licenses DROP COLUMN product_slug;
