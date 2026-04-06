import type { RequestHandler } from "./$types";
import { json } from "@sveltejs/kit";
import { db } from "$lib/server/db/db";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;

    const album = await db.selectFrom('albums')
        .selectAll()
        .where('albums.id', '=', id)
        .where('albums.deleted_at', 'is', null)
        .limit(1)
        .executeTakeFirst()

    if (!album) {
        return new Response("Album not found", { status: 404 });
    }

    const assets = await db.selectFrom('assets')
        .selectAll()
        .where('album_id', '=', album.id)
        .where('deleted_at', 'is', null)
        .execute()

    const ids = assets.map(asset => asset.id);

    return json({
        ...album,
        assets: ids
    })
}