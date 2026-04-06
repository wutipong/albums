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

export const getAlbumAssetWithoutPreviewQuery = `-- name: GetAlbumAssetWithoutPreview :many
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view 
    FROM assets 
    WHERE album_id = $1 AND preview = '' AND deleted_at IS NULL`;

export interface GetAlbumAssetWithoutPreviewArgs {
    albumId: string;
}

export interface GetAlbumAssetWithoutPreviewRow {
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

export async function getAlbumAssetWithoutPreview(sql: Sql, args: GetAlbumAssetWithoutPreviewArgs): Promise<GetAlbumAssetWithoutPreviewRow[]> {
    return (await sql.unsafe(getAlbumAssetWithoutPreviewQuery, [args.albumId]).values()).map(row => ({
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

export const getAlbumAssetWithoutThumbnailQuery = `-- name: GetAlbumAssetWithoutThumbnail :many
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view 
    FROM assets 
    WHERE album_id = $1 AND thumbnail = '' AND deleted_at IS NULL`;

export interface GetAlbumAssetWithoutThumbnailArgs {
    albumId: string;
}

export interface GetAlbumAssetWithoutThumbnailRow {
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

export async function getAlbumAssetWithoutThumbnail(sql: Sql, args: GetAlbumAssetWithoutThumbnailArgs): Promise<GetAlbumAssetWithoutThumbnailRow[]> {
    return (await sql.unsafe(getAlbumAssetWithoutThumbnailQuery, [args.albumId]).values()).map(row => ({
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

export const getAlbumAssetWithoutViewQuery = `-- name: GetAlbumAssetWithoutView :many
SELECT id, album_id, filename, checksum, created_at, modified_at, deleted_at, size, type, original, preview, thumbnail, view 
    FROM assets 
    WHERE album_id = $1 AND view = '' AND deleted_at IS NULL`;

export interface GetAlbumAssetWithoutViewArgs {
    albumId: string;
}

export interface GetAlbumAssetWithoutViewRow {
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

export async function getAlbumAssetWithoutView(sql: Sql, args: GetAlbumAssetWithoutViewArgs): Promise<GetAlbumAssetWithoutViewRow[]> {
    return (await sql.unsafe(getAlbumAssetWithoutViewQuery, [args.albumId]).values()).map(row => ({
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

