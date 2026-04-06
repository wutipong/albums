import { db } from "$lib/server/db/db"
import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;

    const asset = await db.selectFrom('assets')
        .selectAll()
        .where("id", "=", id)
        .where("assets.deleted_at", "is", null)
        .executeTakeFirst()

    if (!asset) {
        return new Response("Asset not found", { status: 404 });
    }

    return json(asset);
};