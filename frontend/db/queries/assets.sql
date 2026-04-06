-- name: GetAsset :one
SELECT * FROM assets WHERE id = $1 AND deleted_at IS NULL LIMIT 1; 

-- name: ListAssetsByAlbum :many
SELECT * FROM assets WHERE album_id = $1 AND deleted_at IS NULL ORDER BY filename ASC;

-- name: CreateAsset :one
INSERT INTO assets (
  album_id,
  filename,
  checksum,
  type,
  original,
  size
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetDuplicatedAlbumAsset :one
SELECT * 
    FROM assets 
    WHERE album_id = $1 AND checksum = $2 AND size = $3 AND deleted_at IS NULL 
    LIMIT 1;

-- name: GetAlbumAssetWithoutPreview :many
SELECT * 
    FROM assets 
    WHERE album_id = $1 AND preview = '' AND deleted_at IS NULL;

-- name: GetAlbumAssetWithoutThumbnail :many
SELECT * 
    FROM assets 
    WHERE album_id = $1 AND thumbnail = '' AND deleted_at IS NULL;

-- name: GetAlbumAssetWithoutView :many
SELECT * 
    FROM assets 
    WHERE album_id = $1 AND view = '' AND deleted_at IS NULL;

-- name: UpdateAsset :one
UPDATE assets SET
  filename = $2,
  checksum = $3,
  type = $4,
  original = $5,
  preview = $6,
  thumbnail = $7,
  view = $8,
  modified_at = NOW ()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;