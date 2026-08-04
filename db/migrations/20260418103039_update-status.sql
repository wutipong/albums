-- migrate:up
ALTER TYPE process_status_t ADD VALUE IF NOT EXISTS 'uploading';

-- migrate:down
SELECT 1;
