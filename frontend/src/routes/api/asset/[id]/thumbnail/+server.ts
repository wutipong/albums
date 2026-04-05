import { createCacheAssetPath } from "$lib/cache";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;

    const data = await fs.readFile(createCacheAssetPath(id, "thumbnail.webp"));

    return new Response(data, {
        headers: {
            "Content-Type": "image/webp"
        }
    });
};