-- migrate:up
CREATE TYPE process_status_t AS ENUM ('pending', 'processing', 'processed');
ALTER TABLE assets ADD COLUMN process_status process_status_t NOT NULL DEFAULT 'pending';

-- migrate:down
ALTER TABLE assets DROP COLUMN process_status;
DROP TYPE process_status_t;