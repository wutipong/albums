import { db } from "$lib/server/db";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, fetch }) => {
    const { id } = params;

    const assets = await db.selectFrom("assets")
        .selectAll()
        .where('album_id', '=', id)
        .where('deleted_at', 'is', null)
        .orderBy('filename')
        .execute()

    const outAssets = []
    for (const asset of assets) {
        const video_duration = asset.video_duration.seconds
        const out = { ...asset, video_duration, }

        outAssets.push(out)
    }

    const album = await db.selectFrom('albums')
        .selectAll()
        .where('id', '=', id)
        .executeTakeFirst()

    return { ...album, assets: outAssets };
};
