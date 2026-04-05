CREATE TABLE
  albums (
    id BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    deleted_at TIMESTAMPTZ
  );