import { QueryArrayConfig, QueryArrayResult } from "pg";

interface Client {
    query: (config: QueryArrayConfig) => Promise<QueryArrayResult>;
}

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

export async function getAlbum(client: Client, args: GetAlbumArgs): Promise<GetAlbumRow | null> {
    const result = await client.query({
        text: getAlbumQuery,
        values: [args.id],
        rowMode: "array"
    });
    if (result.rows.length !== 1) {
        return null;
    }
    const row = result.rows[0];
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

export async function listAlbums(client: Client): Promise<ListAlbumsRow[]> {
    const result = await client.query({
        text: listAlbumsQuery,
        values: [],
        rowMode: "array"
    });
    return result.rows.map(row => {
        return {
            id: row[0],
            name: row[1],
            createdAt: row[2],
            modifiedAt: row[3],
            deletedAt: row[4]
        };
    });
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

export async function createAlbum(client: Client, args: CreateAlbumArgs): Promise<CreateAlbumRow | null> {
    const result = await client.query({
        text: createAlbumQuery,
        values: [args.name],
        rowMode: "array"
    });
    if (result.rows.length !== 1) {
        return null;
    }
    const row = result.rows[0];
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

export async function updateAlbum(client: Client, args: UpdateAlbumArgs): Promise<void> {
    await client.query({
        text: updateAlbumQuery,
        values: [args.id, args.name],
        rowMode: "array"
    });
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

export async function deleteAlbum(client: Client, args: DeleteAlbumArgs): Promise<void> {
    await client.query({
        text: deleteAlbumQuery,
        values: [args.id, args.name],
        rowMode: "array"
    });
}

