import type { Sql } from "postgres";

export const getAlbumQuery = `-- name: GetAlbum :one
SELECT id, name, created_at, modified_at, deleted_at FROM albums
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1`;

export interface GetAlbumArgs {
    id: string;
}

export interface GetAlbumRow {
    id: string;
    name: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
}

export async function getAlbum(sql: Sql, args: GetAlbumArgs): Promise<GetAlbumRow | null> {
    const rows = await sql.unsafe(getAlbumQuery, [args.id]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        id: row[0],
        name: row[1],
        createdAt: row[2],
        modifiedAt: row[3],
        deletedAt: row[4]
    };
}

export const listAlbumsQuery = `-- name: ListAlbums :many
SELECT id, name, created_at, modified_at, deleted_at FROM albums
WHERE deleted_at IS NULL
ORDER BY name`;

export interface ListAlbumsRow {
    id: string;
    name: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
}

export async function listAlbums(sql: Sql): Promise<ListAlbumsRow[]> {
    return (await sql.unsafe(listAlbumsQuery, []).values()).map(row => ({
        id: row[0],
        name: row[1],
        createdAt: row[2],
        modifiedAt: row[3],
        deletedAt: row[4]
    }));
}

export const createAlbumQuery = `-- name: CreateAlbum :one
INSERT INTO albums (
  name
) VALUES (
  $1
)
RETURNING id, name, created_at, modified_at, deleted_at`;

export interface CreateAlbumArgs {
    name: string;
}

export interface CreateAlbumRow {
    id: string;
    name: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
}

export async function createAlbum(sql: Sql, args: CreateAlbumArgs): Promise<CreateAlbumRow | null> {
    const rows = await sql.unsafe(createAlbumQuery, [args.name]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        id: row[0],
        name: row[1],
        createdAt: row[2],
        modifiedAt: row[3],
        deletedAt: row[4]
    };
}

export const updateAlbumQuery = `-- name: UpdateAlbum :exec
UPDATE albums
  set name = $2,
      modified_at = NOW ()  
WHERE id = $1`;

export interface UpdateAlbumArgs {
    id: string;
    name: string;
}

export async function updateAlbum(sql: Sql, args: UpdateAlbumArgs): Promise<void> {
    await sql.unsafe(updateAlbumQuery, [args.id, args.name]);
}

export const deleteAlbumQuery = `-- name: DeleteAlbum :exec
UPDATE albums
  set name = $2,
      deleted_at = NOW ()  
WHERE id = $1`;

export interface DeleteAlbumArgs {
    id: string;
    name: string;
}

export async function deleteAlbum(sql: Sql, args: DeleteAlbumArgs): Promise<void> {
    await sql.unsafe(deleteAlbumQuery, [args.id, args.name]);
}

