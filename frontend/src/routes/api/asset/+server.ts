import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";
import type { File } from "node:buffer";
import { json } from "@sveltejs/kit";
import { db } from "$lib/server/db";

import { createHash } from 'node:crypto';
import { createCacheAssetPath } from "$lib/server/cache";
import path from "node:path";

export const POST: RequestHandler = async ({ request }) => {
    const data = await request.formData();
    const file = data.get("file") as File;
    const albumId = data.get("albumId") as string;

    const album = await db.selectFrom('albums')
        .selectAll()
        .where('albums.id', '=', albumId)
        .where('albums.deleted_at', 'is', null)
        .limit(1)
        .executeTakeFirst()

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
        const asset = await db.insertInto("assets")
            .values({
                album_id: albumId,
                filename: file.name,
                checksum,
                type: "image",
                original: path.join("original", basename),
                size: buffer.length.toString()
            })
            .returningAll()
            .executeTakeFirstOrThrow();

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
    const existingAsset = await db.selectFrom('assets')
        .selectAll()
        .where('album_id', '=', albumId)
        .where('assets.checksum', '=', checksum)
        .where('assets.size', '=', size)
        .where('assets.deleted_at', 'is', null)
        .executeTakeFirst()

    return existingAsset !== null;
}