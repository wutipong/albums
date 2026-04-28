import type { RequestHandler } from "./$types";
import { notifyUpdateAlbumCover } from "$lib/server/grpc/worker";
import { json } from "@sveltejs/kit";

export const POST: RequestHandler = async ({ params, request }) => {
    const { id } = params;
    const input = await request.json()

    try {
        await notifyUpdateAlbumCover(id, input.asset_id)
    } catch {
        return json({ success: false })
    }
    return json({ success: true })
}