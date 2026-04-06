import { createCacheAssetPath } from "$lib/cache";
import { getAsset } from "$lib/server/db/assets_sql";
import { getDb } from "$lib/server/db/db";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;
    const asset = await getAsset(getDb(), {id})

    const data = await fs.readFile(createCacheAssetPath(id, asset?.thumbnail??""));

    return new Response(data, {
        headers: {
            "Content-Type": "image/webp"
        }
    });
};