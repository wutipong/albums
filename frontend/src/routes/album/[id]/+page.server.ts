import { env } from "$env/dynamic/private";
import { db } from "$lib/server/db";
import { s3, s3Public } from "$lib/server/s3";
import { GetObjectCommand } from "@aws-sdk/client-s3";
import type { PageServerLoad } from "./$types";
import { generateImageUrl } from '@imgproxy/imgproxy-node'
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";
import { URL } from "node:url";

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
        if (asset.process_status === 'processed') {
            const video_duration = asset.video_duration.seconds

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

            const preview_url = generateImageUrl({
                endpoint: env.IMGPROXY_URL,
                url: `s3://${env.S3_BUCKET}/${asset.preview}`,
                options: {
                    resizing_type: "auto",
                    height: 200,
                    enlarge: 1,
                },
                salt: env.IMGPROXY_SALT,
                key: env.IMGPROXY_KEY,
            })

            let view_url = '';
            switch (asset.type) {
                case "image":
                    view_url = generateImageUrl({
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
                    break;

                case "video":
                    view_url = await getSignedUrl(
                        s3Public,
                        new GetObjectCommand({
                            Bucket: env.S3_BUCKET,
                            Key: asset.view
                        })
                    )
                    break;
            }
            const out = { ...asset, video_duration, thumbnail_url, preview_url, view_url }
            outAssets.push(out)
        }
        else {
            const out = { ...asset, video_duration: 0, thumbnail_url: '', preview_url: '', view_url: '' }
            outAssets.push(out)
        }

    }

    const album = await db.selectFrom('albums')
        .selectAll()
        .where('id', '=', id)
        .executeTakeFirst()

    return { ...album, assets: outAssets };
};
