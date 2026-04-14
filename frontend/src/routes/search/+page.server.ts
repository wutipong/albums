import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { encodeText } from "$lib/server/grpc/clip";
import { db } from "$lib/server/db";
import { cosineDistance } from "pgvector/kysely";

export const load: PageServerLoad = async ({ params, fetch, url }) => {
    const search = url.searchParams.get('search')

    if (!search) {
        return error(400, 'missing search parameter')
    }

    const resp = await encodeText(search)

    const floatArray = new Float32Array(
        resp.embedding.buffer,
        resp.embedding.byteOffset,
        resp.embedding.length / 4 // 4 bytes per float
    );

    const vectorArray = Array.from(floatArray);

    const assets = await db.selectFrom('assets')
        .selectAll()
        .where("deleted_at", 'is', null)
        .where("type", '=', 'image')
        .orderBy(cosineDistance('image_embedding', vectorArray))
        .limit(50)
        .execute()

    const outAssets = []
    for (const asset of assets) {
        const video_duration = asset.video_duration.seconds
        const out = { ...asset, video_duration,}

        outAssets.push(out)
    }

    return { search, assets: outAssets };
};