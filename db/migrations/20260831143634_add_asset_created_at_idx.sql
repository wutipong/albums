-- migrate:up
CREATE INDEX idx_assets_created_at_object_partial 
    ON assets (created_at) 
    WHERE view <> ''
        OR original <> ''
        OR thumbnail <> ''
        OR preview <> '';

-- migrate:down
DROP INDEX idx_assets_created_at_object_partial;