-- migrate:up
ALTER TABLE albums DROP COLUMN cover;
ALTER TABLE albums ADD COLUMN cover TEXT NOT NULL DEFAULT '';

-- migrate:down
ALTER TABLE albums DROP COLUMN cover;
ALTER TABLE albums ADD COLUMN cover uuid DEFAULT NULL;