-- migrate:up
CREATE TABLE
  albums (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    deleted_at TIMESTAMPTZ
  );

-- migrate:down
DROP TABLE albums;
