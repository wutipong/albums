import type { RequestHandler } from "./$types";
import { json } from "@sveltejs/kit";
import { db } from "$lib/server/db";
import { s3 } from "$lib/server/s3";
import { CopyObjectCommand, DeleteObjectCommand } from "@aws-sdk/client-s3";
import { env } from "$env/dynamic/private";
import * as mime from 'mime-types'
import { notifyProcessAsset } from "$lib/server/grpc/worker";

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
    asset.process_status = 'pending'
    asset.original = oldKey.replace('pending', 'public')
    const newKey = asset.original

    const mimetype = mime.lookup(asset.filename)
    if (!mimetype){
        return json({ success: false, error: "invalid content type" }, { status: 400 });
    }

    if(mimetype.startsWith('image/')){
        asset.type = 'image'
    } else if (mimetype.startsWith('video')){
        asset.type = 'video'
    }

    try {
        // 1. Copy the object to the permanent location
        const r1 = await s3.send(new CopyObjectCommand({
            Bucket: env.S3_BUCKET,
            Key: newKey,
            CopySource: `${env.S3_BUCKET}/${oldKey}`
        }));

        // 2. Delete the original "pending" file
        const r2 = await s3.send(new DeleteObjectCommand({
            Bucket: env.S3_BUCKET,
            Key: oldKey
        }));

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

    try{
        await notifyProcessAsset(asset.id)
    }catch (error){
        return json({ success: false, error: "Failed to notify asset processing" }, { status: 500 });
    }

    return json({ asset: asset, success: true });
};
