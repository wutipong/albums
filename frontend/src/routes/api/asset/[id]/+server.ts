import { db } from "$lib/server/db"
import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

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

    return json({...asset, available: true});
};