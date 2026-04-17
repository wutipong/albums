import type { RequestHandler } from "./$types";
import { json } from "@sveltejs/kit";
import { db } from "$lib/server/db";
import { s3 } from "$lib/server/s3";
import { CopyObjectCommand, DeleteObjectCommand } from "@aws-sdk/client-s3";
import { env } from "$env/dynamic/private";

export const POST: RequestHandler = async ({ request }) => {
    const req = await request.json()
    const id = req.id

    let asset = await db.selectFrom("assets")
        .selectAll()
        .where('id', '=', id)
        .where('deleted_at', 'is', null)
        .executeTakeFirst()

    if (!asset) {
        return json({ success: false, error: "Failed to get asset" }, { status: 500 });
    }

    const oldKey = asset.original
    asset.process_status = 'processing'
    asset.original = oldKey.replace('pending', 'public')
    const newKey = asset.original

    try {
        // 1. Copy the object to the permanent location
        const r1 = await s3.send(new CopyObjectCommand({
            Bucket: env.S3_BUCKET,
            Key: newKey,
            // Source must be "bucket-name/path/to/object"
            CopySource: `${env.S3_BUCKET}/${oldKey}`
        }));

        console.log(r1)

        // 2. Delete the original "pending" file
        const r2 = await s3.send(new DeleteObjectCommand({
            Bucket: env.S3_BUCKET,
            Key: oldKey
        }));

        console.log(r2)
    } catch (error) {
        return json({ success: false, error: "Failed to move asset." }, { status: 500 });
    }
    
    const resp = await db.updateTable('assets')
        .set(asset)
        .where('id', '=', asset.id)
        .executeTakeFirst()

    if (!resp) {
        return json({ success: false, error: "Failed to update asset" }, { status: 500 });
    }

    return json({ asset: asset, success: true });
};
