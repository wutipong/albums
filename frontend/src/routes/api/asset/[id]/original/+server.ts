import { createCacheAssetPath } from "$lib/server/cache";
import { db } from "$lib/server/db"
import * as mime from "mime-types"
import fs from "node:fs/promises";
import path from "node:path";
import type { RequestHandler } from "./$types";

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

    const filePath = createCacheAssetPath(id, asset.original)
    const data = await fs.readFile(filePath);
    const contentType = mime.lookup(path.basename(filePath))
    if (!contentType) {
        return new Response("unable to determine asset content-type", { status: 404 });
    }

    return new Response(data, {
        headers: {
            "Content-Type": contentType,
        }
    });
};