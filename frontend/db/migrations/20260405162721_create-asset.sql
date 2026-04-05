-- migrate:up
CREATE TABLE
  assets (
    id BIGSERIAL PRIMARY KEY,
    album_id BIGINT NOT NULL REFERENCES albums (id),
    filename text NOT NULL,
    size BIGINT NOT NULL,
    checksum text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    deleted_at TIMESTAMPTZ
  );

-- migrate:down
DROP TABLE assets;
