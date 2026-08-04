import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { encodeText } from "$lib/server/grpc/clip";
import { db } from "$lib/server/db";
import { cosineDistance } from "pgvector/kysely";
import { env } from "$env/dynamic/private";
import { generateImageUrl } from "@imgproxy/imgproxy-node";
import { Temporal } from 'temporal-polyfill';
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";
import { s3Public } from "$lib/server/s3";
import { GetObjectCommand } from "@aws-sdk/client-s3";

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
        const video_duration = Temporal.Duration.from(asset.video_duration.toISOString())

        const thumbnail_url = generateImageUrl({
            endpoint: env.IMGPROXY_URL,
            url: `s3://${env.S3_BUCKET}/${asset.thumbnail}`,
            options: {
                resizing_type: "auto",
                height: 200,
                enlarge: 1,
            },
            salt: env.IMGPROXY_SALT,
            key: env.IMGPROXY_KEY,
        })

        const bypass = asset.image_frames > 1 || asset.type == 'video'
        const preview_url = generateImageUrl({
                endpoint: env.IMGPROXY_URL,
                url: `s3://${env.S3_BUCKET}/${asset.preview}`,
                options: {
                    raw: bypass,
                    resizing_type: "auto",
                    height: 200,
                    enlarge: 1,
                },
                salt: env.IMGPROXY_SALT,
                key: env.IMGPROXY_KEY,
            })


        const view_url = generateImageUrl({
            endpoint: env.IMGPROXY_URL,
            url: `s3://${env.S3_BUCKET}/${asset.view}`,
            options: {
                resizing_type: "auto",
                height: 2000,
                enlarge: 1,
            },
            salt: env.IMGPROXY_SALT,
            key: env.IMGPROXY_KEY,
        })

        const original_url = await getSignedUrl(
            s3Public,
            new GetObjectCommand({
                Bucket: env.S3_BUCKET,
                Key: asset.original
            })
        )

        const out = { ...asset, video_duration, thumbnail_url, view_url, preview_url, original_url }

        outAssets.push(out)
    }

    return { search, assets: outAssets };
};