import { createCacheAssetPath } from "$lib/server/cache";
import { db } from "$lib/server/db";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";
import * as mime from "mime-types"
import path from "node:path";

import notAvailableSvg from "$lib/assets/not-available.svg?raw"

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;
    const asset = await db.selectFrom('assets')
        .selectAll()
        .where("id", "=", id)
        .where("assets.deleted_at", "is", null)
        .executeTakeFirst()

    if (!asset || asset.process_status !== "processed") {
        return new Response(notAvailableSvg, {
            headers: {
                "Content-Type": "image/svg+xml",
            }
        });
    }

    const filePath = createCacheAssetPath(id, asset.preview)
    const data = await fs.readFile(filePath);

    const contentType = mime.lookup(path.basename(filePath))
    if (!contentType) {
        return new Response("unable to determine asset content-type", { status: 404 });
    }

    return new Response(data, {
        headers: {
            "Content-Type": contentType
        }
    });
};