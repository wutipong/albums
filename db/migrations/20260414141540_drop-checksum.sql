-- migrate:up
ALTER TABLE assets DROP COLUMN checksum;
ALTER TABLE assets DROP COLUMN size;

-- migrate:down
ALTER TABLE assets ADD COLUMN size BIGINT NOT NULL DEFAULT 0;
ALTER TABLE assets ADD COLUMN checksum text NOT NULL DEFAULT '';