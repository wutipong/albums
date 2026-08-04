-- migrate:up
ALTER TYPE process_status_t ADD VALUE IF NOT EXISTS 'failed';

-- migrate:down
SELECT 1;
