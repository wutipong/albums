-- migrate:up
ALTER TABLE assets ADD COLUMN size BIGINT NOT NULL DEFAULT 0;

CREATE TYPE asset_type AS ENUM ('image', 'animated', 'video', 'audio');
ALTER TABLE assets ADD COLUMN type asset_type NOT NULL DEFAULT 'image';

ALTER TABLE assets ADD COLUMN original TEXT NOT NULL DEFAULT '';
ALTER TABLE assets ADD COLUMN preview TEXT NOT NULL DEFAULT '';
ALTER TABLE assets ADD COLUMN thumbnail TEXT NOT NULL DEFAULT '';
ALTER TABLE assets ADD COLUMN view TEXT NOT NULL DEFAULT '';

-- migrate:down
ALTER TABLE assets DROP COLUMN view;
ALTER TABLE assets DROP COLUMN thumbnail;
ALTER TABLE assets DROP COLUMN preview;
ALTER TABLE assets DROP COLUMN original;

ALTER TABLE assets DROP COLUMN type;
DROP TYPE asset_type;

ALTER TABLE assets DROP COLUMN size;