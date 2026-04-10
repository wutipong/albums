-- migrate:up
ALTER TABLE assets ADD COLUMN thumbnail_width INT NOT NULL DEFAULT 0;
ALTER TABLE assets ADD COLUMN thumbnail_height INT NOT NULL DEFAULT 0;

ALTER TABLE assets ADD COLUMN view_width INT NOT NULL DEFAULT 0;
ALTER TABLE assets ADD COLUMN view_height INT NOT NULL DEFAULT 0;
-- migrate:down
ALTER TABLE assets DROP COLUMN thumbnail_width;
ALTER TABLE assets DROP COLUMN thumbnail_height;

ALTER TABLE assets DROP COLUMN view_width;
ALTER TABLE assets DROP COLUMN view_height;