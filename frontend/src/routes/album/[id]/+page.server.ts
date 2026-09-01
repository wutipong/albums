import { db } from '$lib/server/db';
import type { PageServerLoad } from './$types';
import { createResponseAssetList } from '$lib/server/asset';
import { error } from '@sveltejs/kit';

export const ssr = false;

export const load: PageServerLoad = async ({ params, locals }) => {
	const { id } = params;

	const album = await db.selectFrom('albums').selectAll().where('id', '=', id).executeTakeFirst();

	if (album == undefined){
		error(404, 'Album not found.')
	}
	
	const assets = await db
		.selectFrom('assets')
		.selectAll()
		.where('album_id', '=', id)
		.where('deleted_at', 'is', null)
		.orderBy('filename')
		.execute();

	const outAssets = await createResponseAssetList(assets);

	return { ...album, assets: outAssets, user: locals.user, session: locals.session };
};
