-- migrate:up
ALTER TABLE assets ADD COLUMN image_frames INT NOT NULL DEFAULT 0;
ALTER TABLE assets ADD COLUMN video_duration INTERVAL NOT NULL DEFAULT '0 seconds';

-- migrate:down
ALTER TABLE assets DROP COLUMN video_duration;
ALTER TABLE assets DROP COLUMN image_frames;
