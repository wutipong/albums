import { env } from "$env/dynamic/private";
import { db } from "$lib/server/db";
import type { PageServerLoad } from "./$types";
import { generateImageUrl } from '@imgproxy/imgproxy-node'

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

        const thumbnail_url = generateImageUrl({
            endpoint: env.IMGPROXY_URL,
            url: `s3://${env.S3_BUCKET}/${asset.original}`,
            options: {
                resizing_type: "auto",
                height: 200,
                enlarge: 1,
            },
            salt: env.IMGPROXY_SALT,
            key: env.IMGPROXY_KEY,
        })

        const preview_url = generateImageUrl({
            endpoint: env.IMGPROXY_URL,
            url: `s3://${env.S3_BUCKET}/${asset.original}`,
            options: {
                resizing_type: "auto",
                height: 200,
                enlarge: 1,
            },
            salt: env.IMGPROXY_SALT,
            key: env.IMGPROXY_KEY,
        })

        const view_url = generateImageUrl({
            endpoint: env.IMGPROXY_URL,
            url: `s3://${env.S3_BUCKET}/${asset.original}`,
            options: {
                resizing_type: "auto",
                height: 2000,
                enlarge: 1,
            },
            salt: env.IMGPROXY_SALT,
            key: env.IMGPROXY_KEY,
        })

        const out = { ...asset, video_duration, thumbnail_url, preview_url, view_url }

        outAssets.push(out)
    }

    const album = await db.selectFrom('albums')
        .selectAll()
        .where('id', '=', id)
        .executeTakeFirst()

    return { ...album, assets: outAssets };
};
