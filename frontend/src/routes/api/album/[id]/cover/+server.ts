import { createCacheAlbumPath } from "$lib/server/cache";
import { db } from "$lib/server/db";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";
import * as mime from 'mime-types'

import notAvailableSvg from "$lib/assets/not-available-small.svg?raw"

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