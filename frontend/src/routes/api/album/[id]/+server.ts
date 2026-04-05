import path from "node:path";
import type { RequestHandler } from "./$types";
import fs from "node:fs/promises";
import {json} from "@sveltejs/kit";

import { getDb } from "$lib/db/db";
import { getAlbum } from "$lib/db/albums_sql";


export const GET: RequestHandler = async ({ params }) => {
    const { id } = params;

    const album = await getAlbum(getDb(), { id });
    const dir = await fs.opendir("./cache");
    
    const ids = [];

    for await (const dirent of dir) {
        if (dirent.isFile()) {
            ids.push(path.parse(dirent.name).name.split(".")[0]);
        }
    }

    return json({
        name: album?.name,
        assets: ids
    })
}