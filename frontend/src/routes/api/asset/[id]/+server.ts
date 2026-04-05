import { getAsset } from "$lib/db/assets_sql";
import { getDb } from "$lib/db/db";
import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;

    const asset = await getAsset(getDb(), { id });

    if (!asset) {
        return new Response("Asset not found", { status: 404 });
    }
    
    return json(asset);
};