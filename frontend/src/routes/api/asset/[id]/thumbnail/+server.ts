import { createCacheAssetPath } from "$lib/cache";
import { db } from "$lib/server/db/db";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;
    const asset = await db.selectFrom('assets')
        .selectAll()
        .where("id", "=", id)
        .where("assets.deleted_at", "is", null)
        .executeTakeFirst()

    if (!asset || asset.thumbnail === "") {
        return new Response("Asset not found", { status: 404 });
    }

    const data = await fs.readFile(createCacheAssetPath(id, asset.thumbnail));

    return new Response(data, {
        headers: {
            "Content-Type": "image/webp"
        }
    });
};