import { createCacheAssetPath } from "$lib/server/cache";
import { db } from "$lib/server/db";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";

import notAvailableSvg from "$lib/assets/not-available.svg?raw"
import { json } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;
    const asset = await db.selectFrom('assets')
        .selectAll()
        .where("id", "=", id)
        .where("assets.deleted_at", "is", null)
        .executeTakeFirst()

    if (!asset || asset.process_status !== "processed") {
        return json({
            "available": false,
            "thumbnail_width": 267,
            "thumbnail_height": 200,
        })
    }

    return json({
        "available": true,
        "thumbnail_width": asset.thumbnail_width,
        "thumbnail_height": asset.thumbnail_height,
    })
};