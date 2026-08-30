import { db } from '$lib/server/db';
import type { PageServerLoad } from './$types';
import { createResponseAssetList } from '$lib/server/asset';

export const ssr = false;

export const load: PageServerLoad = async ({ params, locals }) => {
	const { id } = params;

	const assets = await db
		.selectFrom('assets')
		.selectAll()
		.where('album_id', '=', id)
		.where('deleted_at', 'is', null)
		.orderBy('filename')
		.execute();

	const outAssets = await createResponseAssetList(assets);
	const album = await db.selectFrom('albums').selectAll().where('id', '=', id).executeTakeFirst();

	return { ...album, assets: outAssets, user: locals.user, session: locals.session };
};
