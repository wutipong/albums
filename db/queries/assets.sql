-- name: GetAsset :one
SELECT * FROM assets WHERE id = $1 AND deleted_at IS NULL LIMIT 1; 

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

-- name: GetPendingAssets :many
SELECT * 
    FROM assets 
    WHERE process_status = 'pending' AND deleted_at IS NULL;

-- name: UpdateAsset :one
UPDATE assets SET
  filename = $2,
  checksum = $3,
  type = $4,
  original = $5,
  preview = $6,
  thumbnail = $7,
  view = $8,
  process_status = $9,
  modified_at = NOW ()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: GetAssetProcessStatus :one
SELECT process_status
  FROM assets
  WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateAssetProcessStatus :one
UPDATE assets SET
  process_status = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;