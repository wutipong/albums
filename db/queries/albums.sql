
-- name: GetAlbum :one
SELECT * FROM albums WHERE id = $1 and deleted_at IS NULL;

-- name: GetAllAlbum :many
SELECT * FROM albums WHERE deleted_at IS NULL;

-- name: GetAlbumsWithoutCover :many
SELECT * FROM albums WHERE cover = '' and deleted_at IS NULL;

-- name: MarkAlbumDeleted :exec
UPDATE albums
SET deleted_at = now()
WHERE id = $1;
