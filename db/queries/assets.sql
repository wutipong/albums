-- name: GetAsset :one
SELECT * FROM assets WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetAssetsByType :many
SELECT *
FROM assets
WHERE
    type = $1
    AND deleted_at IS NULL;

-- name: CreateAsset :one
INSERT INTO
    assets (
        album_id,
        filename,
        type,
        original
    )
VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: GetPendingAssets :many
SELECT *
FROM assets
WHERE
    process_status = 'pending'
    AND deleted_at IS NULL;

-- name: GetPendingOrFailedAssets :many
SELECT *
FROM assets
WHERE
    (process_status = 'pending' OR process_status = 'failed')
    AND deleted_at IS NULL;

-- name: GetAssetsWithoutUploading :many
SELECT *
FROM assets
WHERE
    process_status <> 'uploading'
    AND deleted_at IS NULL;

-- name: GetImageAssetsWithoutEmbedding :many
SELECT *
FROM assets
WHERE
    type = 'image'
    AND deleted_at IS NULL
    AND image_embedding IS NULL;

-- name: GetAssetsWithObjects :many
SELECT *
    FROM assets
    WHERE (original <> '' OR view <> '' OR thumbnail <> '' OR preview <> '')
    ORDER BY assets.id
    LIMIT $1 OFFSET $2;

-- name: IsObjectInUse :one
SELECT EXISTS (
    SELECT 1
        FROM assets
        WHERE original = $1 OR view = $1 OR thumbnail = $1 OR preview = $1
);

-- name: UpdateAsset :one
UPDATE assets
SET
    filename = $2,
    type = $3,
    original = $4,
    preview = $5,
    thumbnail = $6,
    view = $7,
    process_status = $8,
    modified_at = NOW(),
    thumbnail_width = $9,
    thumbnail_height = $10,
    view_width = $11,
    view_height = $12,
    image_frames = $13,
    video_duration = $14,
    image_embedding = $15
WHERE
    id = $1
    AND deleted_at IS NULL
RETURNING
    *;

-- name: GetAssetProcessStatus :one
SELECT process_status
FROM assets
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: UpdateAssetProcessStatus :one
UPDATE assets
SET
    process_status = $2,
    modified_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL
RETURNING
    *;

-- name: GetAlbumAssets :many
SELECT id
FROM assets
WHERE
    album_id = $1
    and deleted_at IS NULL;

-- name: GetAlbum :one
SELECT * FROM albums WHERE id = $1 and deleted_at IS NULL;

-- name: GetAllAlbum :many
SELECT * FROM albums WHERE deleted_at IS NULL;

-- name: UpdateAlbumThumbnail :one
UPDATE albums
SET
    cover = $1,
    modified_at = NOW()
WHERE
    id = $2
    AND deleted_at IS NULL
RETURNING
    *;

-- name: GetAlbumsWithoutCover :many
SELECT * FROM albums WHERE cover = '' and deleted_at IS NULL;

-- name: GetRandomAlbumAsset :one
SELECT *
FROM assets
WHERE
    type <> 'audio'
    AND album_id = $1
    AND process_status = 'processed'
    AND deleted_at IS NULL
ORDER BY RANDOM()
LIMIT 1;

-- name: GetAlbumAssetForCover :one
SELECT *
FROM assets
WHERE
    type <> 'audio'
    AND album_id = $1
    AND process_status = 'processed'
    AND deleted_at IS NULL
ORDER BY filename ASC
LIMIT 1;

-- name: GetAlbumPortraitAssetForCover :one
SELECT *
FROM assets
WHERE
    type <> 'audio'
    AND album_id = $1
    AND process_status = 'processed'
    AND view_width < view_height
    AND deleted_at IS NULL
ORDER BY filename ASC
LIMIT 1;

-- name: GetAlbumLandscapeAssetForCover :one
SELECT *
FROM assets
WHERE
    type <> 'audio'
    AND album_id = $1
    AND process_status = 'processed'
    AND view_width > view_height
    AND deleted_at IS NULL
ORDER BY filename ASC
LIMIT 1;
