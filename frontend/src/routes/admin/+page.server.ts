import { db } from "$lib/server/db";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
    const count = await db.selectFrom("assets").select((eb) => [eb.fn.countAll().as("count")]).executeTakeFirst();
    const pendings = await db.selectFrom("assets").
        leftJoin("albums", "assets.album_id", "albums.id").
        select('assets.id as asset_id').
        select('albums.name as album_name').
        select('assets.created_at as created_at').
        where("process_status", "=", "pending").
        orderBy("assets.created_at", "asc").
        limit(20).execute();

    const pendingCount = await db.selectFrom("assets").select((eb) => [eb.fn.countAll().as("count")]).where("process_status", "=", "pending").executeTakeFirst();
    
    const album_count = await db.selectFrom("albums").select((eb) => [eb.fn.countAll().as("count")]).executeTakeFirst();

    return { 
        count: count ? BigInt(count.count) : 0n, 
        pendings, 
        pendingCount: pendingCount ? BigInt(pendingCount.count) : 0n,
        albumCount: album_count ? BigInt(album_count.count) : 0n,
     };
}