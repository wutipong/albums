-- migrate:up
ALTER TYPE process_status_t ADD VALUE 'uploading';

-- migrate:down
SELECT 1;
