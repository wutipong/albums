-- name: GetAsset :one
SELECT * FROM assets WHERE id = $1 AND deleted_at IS NULL LIMIT 1; 

-- name: ListAssetsByAlbum :many
SELECT * FROM assets WHERE album_id = $1 AND deleted_at IS NULL ORDER BY filename ASC;

-- name: CreateAsset :one
INSERT INTO assets (
  album_id,
  filename,
  size,
  checksum
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;