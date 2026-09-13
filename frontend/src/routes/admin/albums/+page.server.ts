import { db } from '$lib/server/db';
import type { PageServerLoad } from '../asset/$types';
import log from '$lib/log';

export const load: PageServerLoad = async () => {
	const album_count = await db
		.selectFrom('albums')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('deleted_at', 'is', null)
		.executeTakeFirst();

	log.debug({ album_count }, 'Album count query result.');

	const missing_cover = await db
		.selectFrom('albums')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('albums.cover', '=', '')
		.where('deleted_at', 'is', null)
		.executeTakeFirst();

	log.debug({ missing_cover }, 'Missing cover query result.');

	return {
		total: album_count ? BigInt(album_count.count) : 0n,
		missingCover: missing_cover ? BigInt(missing_cover.count) : 0n
	};
};
