import type { Sql } from "postgres";

export const getAssetQuery = `-- name: GetAsset :one
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at FROM assets WHERE id = $1 AND deleted_at IS NULL LIMIT 1`;

export interface GetAssetArgs {
    id: string;
}

export interface GetAssetRow {
    id: string;
    albumId: string;
    filename: string;
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
        checksum: row[3],
        createdAt: row[4],
        modifiedAt: row[5],
        deletedAt: row[6]
    };
}

export const listAssetsByAlbumQuery = `-- name: ListAssetsByAlbum :many
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at FROM assets WHERE album_id = $1 AND deleted_at IS NULL ORDER BY filename ASC`;

export interface ListAssetsByAlbumArgs {
    albumId: string;
}

export interface ListAssetsByAlbumRow {
    id: string;
    albumId: string;
    filename: string;
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
        checksum: row[3],
        createdAt: row[4],
        modifiedAt: row[5],
        deletedAt: row[6]
    }));
}

export const createAssetQuery = `-- name: CreateAsset :one
INSERT INTO assets (
  album_id,
  filename,
  checksum
) VALUES (
  $1, $2, $3
)
RETURNING id, album_id, filename, checksum, created_at, modified_at, deleted_at`;

export interface CreateAssetArgs {
    albumId: string;
    filename: string;
    checksum: string;
}

export interface CreateAssetRow {
    id: string;
    albumId: string;
    filename: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
}

export async function createAsset(sql: Sql, args: CreateAssetArgs): Promise<CreateAssetRow | null> {
    const rows = await sql.unsafe(createAssetQuery, [args.albumId, args.filename, args.checksum]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        id: row[0],
        albumId: row[1],
        filename: row[2],
        checksum: row[3],
        createdAt: row[4],
        modifiedAt: row[5],
        deletedAt: row[6]
    };
}

export const getAlbumAssetsByFilenameQuery = `-- name: GetAlbumAssetsByFilename :many
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at 
    FROM assets 
    WHERE album_id = $1 AND filename = $2 AND deleted_at IS NULL 
    ORDER BY id ASC`;

export interface GetAlbumAssetsByFilenameArgs {
    albumId: string;
    filename: string;
}

export interface GetAlbumAssetsByFilenameRow {
    id: string;
    albumId: string;
    filename: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
}

export async function getAlbumAssetsByFilename(sql: Sql, args: GetAlbumAssetsByFilenameArgs): Promise<GetAlbumAssetsByFilenameRow[]> {
    return (await sql.unsafe(getAlbumAssetsByFilenameQuery, [args.albumId, args.filename]).values()).map(row => ({
        id: row[0],
        albumId: row[1],
        filename: row[2],
        checksum: row[3],
        createdAt: row[4],
        modifiedAt: row[5],
        deletedAt: row[6]
    }));
}

export const getAlbumAssetByChecksumQuery = `-- name: GetAlbumAssetByChecksum :one
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at 
    FROM assets 
    WHERE album_id = $1 AND checksum = $2 AND deleted_at IS NULL 
    LIMIT 1`;

export interface GetAlbumAssetByChecksumArgs {
    albumId: string;
    checksum: string;
}

export interface GetAlbumAssetByChecksumRow {
    id: string;
    albumId: string;
    filename: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
}

export async function getAlbumAssetByChecksum(sql: Sql, args: GetAlbumAssetByChecksumArgs): Promise<GetAlbumAssetByChecksumRow | null> {
    const rows = await sql.unsafe(getAlbumAssetByChecksumQuery, [args.albumId, args.checksum]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        id: row[0],
        albumId: row[1],
        filename: row[2],
        checksum: row[3],
        createdAt: row[4],
        modifiedAt: row[5],
        deletedAt: row[6]
    };
}

