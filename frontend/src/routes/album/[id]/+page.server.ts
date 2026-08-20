import { db } from '$lib/server/db';
import type { PageServerLoad } from './$types';

import { createResponseAssetList } from '$lib/server/asset';
import { auth } from '$lib/server/auth';

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

	const user = locals?.user

	return { ...album, assets: outAssets, user };
};
