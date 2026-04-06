import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";
import type { File } from "node:buffer";
import sharp from "sharp";
import { json } from "@sveltejs/kit";
import { getDb } from "$lib/server/db/db";
import { getAlbum } from "$lib/server/db/albums_sql";
import { createAsset, getDuplicatedAlbumAsset } from "$lib/server/db/assets_sql";
import { createHash } from 'node:crypto';
import { createCacheAssetPath } from "$lib/cache";
import path from "node:path";

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
        const checksum = calculateChecksum(buffer);

        if (await checkDuplicate(albumId, checksum, buffer.length.toString())) {
            return json({ success: false, error: "Duplicate asset" }, { status: 400 });
        }

        const basename = path.basename(file.name);
        const asset = await createAsset(getDb(), {
            albumId,
            filename: file.name,
            checksum,
            type: "image",
            original: path.join("original", basename),
            size: buffer.length.toString()
        });
        if (!asset) {
            return json({ success: false, error: "Failed to create asset" }, { status: 500 });
        }

        await fs.mkdir(createCacheAssetPath(asset.id, "original"), { recursive: true });
        await fs.writeFile(createCacheAssetPath(asset.id, "original", basename), buffer);

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

async function checkDuplicate(albumId: string, checksum: string, size: string): Promise<boolean> {
    const existingAsset = await getDuplicatedAlbumAsset(getDb(), { albumId, checksum, size });
    return existingAsset !== null;
}