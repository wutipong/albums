import type { Sql } from "postgres";

export const getAssetQuery = `-- name: GetAsset :one
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view FROM assets WHERE id = $1 AND deleted_at IS NULL LIMIT 1`;

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
    size: string;
    type: string;
    original: string;
    preview: string;
    thumbnail: string;
    view: string;
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
        deletedAt: row[6],
        size: row[7],
        type: row[8],
        original: row[9],
        preview: row[10],
        thumbnail: row[11],
        view: row[12]
    };
}

export const listAssetsByAlbumQuery = `-- name: ListAssetsByAlbum :many
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view FROM assets WHERE album_id = $1 AND deleted_at IS NULL ORDER BY filename ASC`;

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
    size: string;
    type: string;
    original: string;
    preview: string;
    thumbnail: string;
    view: string;
}

export async function listAssetsByAlbum(sql: Sql, args: ListAssetsByAlbumArgs): Promise<ListAssetsByAlbumRow[]> {
    return (await sql.unsafe(listAssetsByAlbumQuery, [args.albumId]).values()).map(row => ({
        id: row[0],
        albumId: row[1],
        filename: row[2],
        checksum: row[3],
        createdAt: row[4],
        modifiedAt: row[5],
        deletedAt: row[6],
        size: row[7],
        type: row[8],
        original: row[9],
        preview: row[10],
        thumbnail: row[11],
        view: row[12]
    }));
}

export const createAssetQuery = `-- name: CreateAsset :one
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
RETURNING id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view`;

export interface CreateAssetArgs {
    albumId: string;
    filename: string;
    checksum: string;
    type: string;
    original: string;
    size: string;
}

export interface CreateAssetRow {
    id: string;
    albumId: string;
    filename: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
    size: string;
    type: string;
    original: string;
    preview: string;
    thumbnail: string;
    view: string;
}

export async function createAsset(sql: Sql, args: CreateAssetArgs): Promise<CreateAssetRow | null> {
    const rows = await sql.unsafe(createAssetQuery, [args.albumId, args.filename, args.checksum, args.type, args.original, args.size]).values();
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
        deletedAt: row[6],
        size: row[7],
        type: row[8],
        original: row[9],
        preview: row[10],
        thumbnail: row[11],
        view: row[12]
    };
}

export const getDuplicatedAlbumAssetQuery = `-- name: GetDuplicatedAlbumAsset :one
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view 
    FROM assets 
    WHERE album_id = $1 AND checksum = $2 AND size = $3 AND deleted_at IS NULL 
    LIMIT 1`;

export interface GetDuplicatedAlbumAssetArgs {
    albumId: string;
    checksum: string;
    size: string;
}

export interface GetDuplicatedAlbumAssetRow {
    id: string;
    albumId: string;
    filename: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
    size: string;
    type: string;
    original: string;
    preview: string;
    thumbnail: string;
    view: string;
}

export async function getDuplicatedAlbumAsset(sql: Sql, args: GetDuplicatedAlbumAssetArgs): Promise<GetDuplicatedAlbumAssetRow | null> {
    const rows = await sql.unsafe(getDuplicatedAlbumAssetQuery, [args.albumId, args.checksum, args.size]).values();
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
        deletedAt: row[6],
        size: row[7],
        type: row[8],
        original: row[9],
        preview: row[10],
        thumbnail: row[11],
        view: row[12]
    };
}

export const getAssetsWithoutPreviewQuery = `-- name: GetAssetsWithoutPreview :many
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view 
    FROM assets 
    WHERE preview = '' AND deleted_at IS NULL`;

export interface GetAssetsWithoutPreviewRow {
    id: string;
    albumId: string;
    filename: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
    size: string;
    type: string;
    original: string;
    preview: string;
    thumbnail: string;
    view: string;
}

export async function getAssetsWithoutPreview(sql: Sql): Promise<GetAssetsWithoutPreviewRow[]> {
    return (await sql.unsafe(getAssetsWithoutPreviewQuery, []).values()).map(row => ({
        id: row[0],
        albumId: row[1],
        filename: row[2],
        checksum: row[3],
        createdAt: row[4],
        modifiedAt: row[5],
        deletedAt: row[6],
        size: row[7],
        type: row[8],
        original: row[9],
        preview: row[10],
        thumbnail: row[11],
        view: row[12]
    }));
}

export const getAssetsWithoutThumbnailQuery = `-- name: GetAssetsWithoutThumbnail :many
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view 
    FROM assets 
    WHERE thumbnail = '' AND deleted_at IS NULL`;

export interface GetAssetsWithoutThumbnailRow {
    id: string;
    albumId: string;
    filename: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
    size: string;
    type: string;
    original: string;
    preview: string;
    thumbnail: string;
    view: string;
}

export async function getAssetsWithoutThumbnail(sql: Sql): Promise<GetAssetsWithoutThumbnailRow[]> {
    return (await sql.unsafe(getAssetsWithoutThumbnailQuery, []).values()).map(row => ({
        id: row[0],
        albumId: row[1],
        filename: row[2],
        checksum: row[3],
        createdAt: row[4],
        modifiedAt: row[5],
        deletedAt: row[6],
        size: row[7],
        type: row[8],
        original: row[9],
        preview: row[10],
        thumbnail: row[11],
        view: row[12]
    }));
}

export const getAssetsWithoutViewQuery = `-- name: GetAssetsWithoutView :many
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view 
    FROM assets 
    WHERE view = '' AND deleted_at IS NULL`;

export interface GetAssetsWithoutViewRow {
    id: string;
    albumId: string;
    filename: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
    size: string;
    type: string;
    original: string;
    preview: string;
    thumbnail: string;
    view: string;
}

export async function getAssetsWithoutView(sql: Sql): Promise<GetAssetsWithoutViewRow[]> {
    return (await sql.unsafe(getAssetsWithoutViewQuery, []).values()).map(row => ({
        id: row[0],
        albumId: row[1],
        filename: row[2],
        checksum: row[3],
        createdAt: row[4],
        modifiedAt: row[5],
        deletedAt: row[6],
        size: row[7],
        type: row[8],
        original: row[9],
        preview: row[10],
        thumbnail: row[11],
        view: row[12]
    }));
}

export const updateAssetQuery = `-- name: UpdateAsset :one
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
RETURNING id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view`;

export interface UpdateAssetArgs {
    id: string;
    filename: string;
    checksum: string;
    type: string;
    original: string;
    preview: string;
    thumbnail: string;
    view: string;
}

export interface UpdateAssetRow {
    id: string;
    albumId: string;
    filename: string;
    checksum: string;
    createdAt: Date;
    modifiedAt: Date;
    deletedAt: Date | null;
    size: string;
    type: string;
    original: string;
    preview: string;
    thumbnail: string;
    view: string;
}

export async function updateAsset(sql: Sql, args: UpdateAssetArgs): Promise<UpdateAssetRow | null> {
    const rows = await sql.unsafe(updateAssetQuery, [args.id, args.filename, args.checksum, args.type, args.original, args.preview, args.thumbnail, args.view]).values();
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
        deletedAt: row[6],
        size: row[7],
        type: row[8],
        original: row[9],
        preview: row[10],
        thumbnail: row[11],
        view: row[12]
    };
}

