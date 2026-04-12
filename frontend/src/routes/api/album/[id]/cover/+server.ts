import { createCacheAlbumPath } from "$lib/server/cache";
import { db } from "$lib/server/db";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";
import * as mime from 'mime-types'

import notAvailableSvg from "$lib/assets/not-available-small.svg?raw"
import { notifyUpdateAlbumCover } from "$lib/server/grpc/worker";
import { json } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;
    const album = await db.selectFrom('albums')
        .selectAll()
        .where("id", "=", id)
        .where("albums.deleted_at", "is", null)
        .executeTakeFirst()

    if (!album || album.cover === "") {
        return new Response(notAvailableSvg, {
            headers: {
                "Content-Type": "image/svg+xml",
            }
        });
    }

    const path = createCacheAlbumPath(id, album.cover)
    const contentType = mime.lookup(path)
    if (!contentType) {
        return new Response(notAvailableSvg, {
            headers: {
                "Content-Type": "image/svg+xml",
            }
        });
    }
    const data = await fs.readFile(path);

    return new Response(data, {
        headers: {
            "Content-Type": contentType,
        }
    });
};

export const POST: RequestHandler = async ({ params, request }) => {
    const { id } = params;
    const input = await request.json()
    try {
        await notifyUpdateAlbumCover(id, input.asset_id)
    } catch {
        return json({ success: false })
    }
    return json({ success: true })
}