import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";
import { notifyProcessAsset } from "$lib/server/grpc/worker";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;

    try {
        notifyProcessAsset(id)
    } catch (err) {
        return json({ success: false, error: "Failed to notify asset processing" }, { status: 500 });
    }

    return json({ success: true });
};