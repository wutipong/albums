import path from "path";
import type { RequestHandler } from "./$types";
import fs from "fs/promises";
import {json} from "@sveltejs/kit";

export const GET: RequestHandler = async (/*{ params }*/) => {
    // const { id } = params;
    const dir = await fs.opendir("./cache");
    
    const ids = [];

    for await (const dirent of dir) {
        if (dirent.isFile()) {
            ids.push(path.parse(dirent.name).name.split(".")[0]);
        }
    }

    return json({
        assets: ids
    })
}