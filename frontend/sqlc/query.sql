-- name: GetAlbum :one
SELECT * FROM albums
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListAlbums :many
SELECT * FROM albums
WHERE deleted_at IS NULL
ORDER BY name;

-- name: CreateAlbum :one
INSERT INTO albums (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: UpdateAlbum :exec
UPDATE albums
  set name = $2,
      modified_at = NOW ()  
WHERE id = $1;

-- name: DeleteAlbum :exec
UPDATE albums
  set name = $2,
      deleted_at = NOW ()  
WHERE id = $1;