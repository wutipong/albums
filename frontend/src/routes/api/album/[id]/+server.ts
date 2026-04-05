import path from "node:path";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";
import {json} from "@sveltejs/kit";

import { getDb } from "$lib/db/db";
import { getAlbum } from "$lib/db/albums_sql";
import { listAssetsByAlbum } from "$lib/db/assets_sql";

export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;
    const album = await getAlbum(getDb(), { id });
    
    const assets = await listAssetsByAlbum(getDb(), { albumId: id });
    const ids = assets.map(asset => asset.id);

    return json({
        name: album?.name,
        assets: ids
    })
}