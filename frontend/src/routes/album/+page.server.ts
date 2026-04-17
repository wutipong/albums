import { db } from "$lib/server/db";
import type { PageLoad } from "./$types";
import { generateImageUrl } from '@imgproxy/imgproxy-node'

export const load: PageLoad = async ({ params, fetch }) => {
	const res = await fetch('/api/album');
	const data = await res.json();

	const albums = await db.selectFrom("albums")
		.selectAll()
		.where('deleted_at', 'is', null)
		.orderBy('name')
		.execute()

	return { albums };
};