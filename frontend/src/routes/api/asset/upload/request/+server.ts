import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";
import type { File } from "node:buffer";
import { json } from "@sveltejs/kit";
import { db } from "$lib/server/db";
import { s3 } from "$lib/server/s3";
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";
import { PutObjectCommand } from "@aws-sdk/client-s3";

import { randomUUID } from 'node:crypto';
import { env } from "$env/dynamic/private";

export const POST: RequestHandler = async ({ request }) => {
    const req = await request.json()
    const albumId = req.album_id;
    const filename = req.filename;
    const checksum = req.checksum;

    const album = await db.selectFrom('albums')
        .selectAll()
        .where('albums.id', '=', albumId)
        .where('albums.deleted_at', 'is', null)
        .limit(1)
        .executeTakeFirst()

    if (!album) {
        return json({ success: false, error: "Album not found" }, { status: 404 });
    }
    const key = `pending/${randomUUID()}`

    const asset = await db.insertInto("assets")
        .values({
            album_id: albumId,
            filename: filename,
            type: "image",
            process_status: "uploading",
            original: key,
        })
        .returningAll()
        .executeTakeFirstOrThrow();

    if (!asset) {
        return json({ success: false, error: "Failed to create asset" }, { status: 500 });
    }

    const command = new PutObjectCommand({
        Bucket: env.S3_BUCKET,
        Key: key,
        ChecksumCRC32: checksum,
    });
    const url = await getSignedUrl(s3, command, { expiresIn: 3600 });

    return json({ id: asset.id, url, success: true });
};
