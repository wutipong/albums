import { db } from '$lib/server/db';
import { sql } from 'kysely';
import type { PageServerLoad } from '../asset/$types';

export const load: PageServerLoad = async () => {
	const album_count = await db
		.selectFrom('albums')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('deleted_at', 'is', null)
		.executeTakeFirst();

	const missing_cover = await db
		.selectFrom('albums')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('albums.cover', '=', '')
		.where('deleted_at', 'is', null)
		.executeTakeFirst();

	return {
		total: album_count ? BigInt(album_count.count) : 0n,
		missingCover: missing_cover ? BigInt(missing_cover.count) : 0n
	};
};
