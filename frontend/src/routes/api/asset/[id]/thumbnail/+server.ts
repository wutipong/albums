import { createCacheAssetPath } from "$lib/server/cache";
import { db } from "$lib/server/db";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";

import notAvailableSvg from "$lib/assets/not-available.svg?raw"

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;
    const asset = await db.selectFrom('assets')
        .selectAll()
        .where("id", "=", id)
        .where("assets.deleted_at", "is", null)
        .executeTakeFirst()

    if (!asset || asset.thumbnail === "") {
        return new Response(notAvailableSvg, {
            headers: {
                "Content-Type": "image/svg+xml",
            }
        });
    }

    const data = await fs.readFile(createCacheAssetPath(id, asset.thumbnail));

    return new Response(data, {
        headers: {
            "Content-Type": "image/webp"
        }
    });
};