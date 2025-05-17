-- +migrate Up
ALTER TABLE accounts ADD COLUMN is_premium BOOLEAN NOT NULL DEFAULT false;

-- +migrate Down
ALTER TABLE accounts DROP COLUMN is_premium;
