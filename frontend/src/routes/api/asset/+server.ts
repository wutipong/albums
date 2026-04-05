import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";
import type { File } from "node:buffer";
import sharp from "sharp";
import { json } from "@sveltejs/kit";
import { getDb } from "$lib/db/db";
import { getAlbum } from "$lib/db/albums_sql";
import { createAsset, getAlbumAssetByChecksum } from "$lib/db/assets_sql";
import { createHash } from 'node:crypto';
import { createCacheAssetPath } from "$lib/cache";

const THUMBNAIL_SIZE = 500;
const VIEW_SIZE = 2000;


export const POST: RequestHandler = async ({ request, params }) => {
    const data = await request.formData();
    const file = data.get("file") as File;
    const albumId = data.get("albumId") as string;

    const album = await getAlbum(getDb(), { id: albumId });
    if (!album) {
        return json({ success: false, error: "Album not found" }, { status: 404 });
    }

    if (!file) {
        return new Response("No file uploaded", { status: 400 });
    }

    const buffer = Buffer.from(await file.arrayBuffer());
    try {
        const s = sharp(buffer)

        const checksum = calculateChecksum(buffer);

        if (await checkDuplicate(checksum, albumId)) {
            return json({ success: false, error: "Duplicate asset" }, { status: 400 });
        }

        const asset = await createAsset(getDb(), {
            albumId,
            filename: file.name,
            checksum
        });
        if (!asset) {
            return json({ success: false, error: "Failed to create asset" }, { status: 500 });
        }

        await fs.mkdir(createCacheAssetPath(asset.id, ""), { recursive: true });

        s.resize(THUMBNAIL_SIZE, THUMBNAIL_SIZE, {
            fit: "inside"
        }).toFormat("webp").toFile(createCacheAssetPath(asset.id, "thumbnail.webp"));

        s.resize(VIEW_SIZE, VIEW_SIZE, {
            fit: "inside"
        }).toFormat("webp").toFile(createCacheAssetPath(asset.id, "view.webp"));

        return json({ asset, success: true });
    } catch (err) {
        console.error(err);
        return json({ success: false, error: "Failed to process image" }, { status: 500 });
    }
};

function calculateChecksum(buffer: Buffer): string {
    const hash = createHash('sha256');
    hash.update(buffer);

    return hash.digest('hex');
}

async function checkDuplicate(checksum: string, albumId: string): Promise<boolean> {
    const existingAsset = await getAlbumAssetByChecksum(getDb(), { albumId, checksum });
    return existingAsset !== null;
}