-- migrate:up
CREATE EXTENSION IF NOT EXISTS vector;
ALTER TABLE assets ADD COLUMN image_embedding vector(512) DEFAULT NULL;
CREATE INDEX image_embedding_idx ON assets USING hnsw (image_embedding vector_cosine_ops);

-- migrate:down
DROP INDEX IF EXISTS image_embedding_idx;
ALTER TABLE assets DROP COLUMN image_embedding;
DROP EXTENSION vector;