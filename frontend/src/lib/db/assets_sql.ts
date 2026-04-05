import type { Sql } from "postgres";

export const getAssetQuery = `-- name: GetAsset :one
SELECT id, album_id, filename, size, checksum, created_at, modified_at, deleted_at FROM assets WHERE id = $1 AND deleted_at IS NULL LIMIT 1`;

export interface GetAssetArgs {
    id: string;
}

export interface GetAssetRow {
    id: string;
    albumId: string;
    filename: string;
    size: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
}

export async function getAsset(sql: Sql, args: GetAssetArgs): Promise<GetAssetRow | null> {
    const rows = await sql.unsafe(getAssetQuery, [args.id]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        id: row[0],
        albumId: row[1],
        filename: row[2],
        size: row[3],
        checksum: row[4],
        createdAt: row[5],
        modifiedAt: row[6],
        deletedAt: row[7]
    };
}

export const listAssetsByAlbumQuery = `-- name: ListAssetsByAlbum :many
SELECT id, album_id, filename, size, checksum, created_at, modified_at, deleted_at FROM assets WHERE album_id = $1 AND deleted_at IS NULL ORDER BY filename ASC`;

export interface ListAssetsByAlbumArgs {
    albumId: string;
}

export interface ListAssetsByAlbumRow {
    id: string;
    albumId: string;
    filename: string;
    size: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
}

export async function listAssetsByAlbum(sql: Sql, args: ListAssetsByAlbumArgs): Promise<ListAssetsByAlbumRow[]> {
    return (await sql.unsafe(listAssetsByAlbumQuery, [args.albumId]).values()).map(row => ({
        id: row[0],
        albumId: row[1],
        filename: row[2],
        size: row[3],
        checksum: row[4],
        createdAt: row[5],
        modifiedAt: row[6],
        deletedAt: row[7]
    }));
}

export const createAssetQuery = `-- name: CreateAsset :one
INSERT INTO assets (
  album_id,
  filename,
  size,
  checksum
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, album_id, filename, size, checksum, created_at, modified_at, deleted_at`;

export interface CreateAssetArgs {
    albumId: string;
    filename: string;
    size: string;
    checksum: string;
}

export interface CreateAssetRow {
    id: string;
    albumId: string;
    filename: string;
    size: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
}

export async function createAsset(sql: Sql, args: CreateAssetArgs): Promise<CreateAssetRow | null> {
    const rows = await sql.unsafe(createAssetQuery, [args.albumId, args.filename, args.size, args.checksum]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        id: row[0],
        albumId: row[1],
        filename: row[2],
        size: row[3],
        checksum: row[4],
        createdAt: row[5],
        modifiedAt: row[6],
        deletedAt: row[7]
    };
}

