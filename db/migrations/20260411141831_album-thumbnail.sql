-- migrate:up
ALTER TABLE albums ADD COLUMN cover uuid DEFAULT NULL;

-- migrate:down
ALTER TABLE albums DROP COLUMN cover;
