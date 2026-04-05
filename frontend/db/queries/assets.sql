-- name: GetAsset :one
SELECT * FROM assets WHERE id = $1 AND deleted_at IS NULL LIMIT 1; 

-- name: ListAssetsByAlbum :many
SELECT * FROM assets WHERE album_id = $1 AND deleted_at IS NULL ORDER BY filename ASC;

-- name: CreateAsset :one
INSERT INTO assets (
  album_id,
  filename,
  checksum
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAlbumAssetsByFilename :many
SELECT * 
    FROM assets 
    WHERE album_id = $1 AND filename = $2 AND deleted_at IS NULL 
    ORDER BY id ASC;

-- name: GetAlbumAssetByChecksum :one
SELECT * 
    FROM assets 
    WHERE album_id = $1 AND checksum = $2 AND deleted_at IS NULL 
    LIMIT 1;