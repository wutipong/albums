-- migrate:up
-- Do nothing
SELECT 1;

-- migrate:down
DROP TABLE IF EXISTS neoq_schema_migrations;
DROP TABLE IF EXISTS neoq_jobs;
DROP TABLE IF EXISTS neoq_dead_jobs;