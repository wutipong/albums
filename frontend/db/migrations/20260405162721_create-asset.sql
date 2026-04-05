-- migrate:up
CREATE TABLE
  assets (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    album_id uuid NOT NULL REFERENCES albums (id),
    filename text NOT NULL,
    checksum text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    deleted_at TIMESTAMPTZ
  );

-- migrate:down
DROP TABLE assets;
