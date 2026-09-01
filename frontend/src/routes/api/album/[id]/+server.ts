import type { RequestHandler } from '@sveltejs/kit';
import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { s3 } from '$lib/server/s3';
import { env } from '$env/dynamic/private';
import { deleteAlbum } from '$lib/server/grpc/worker';

export const GET: RequestHandler = async ({ params }) => {
	const { id } = params;

	if (!id) {
		return json({ error: 'Album id is required.' }, { status: 400 });
	}

	const album = await db
		.selectFrom('albums')
		.where('id', '=', id)
		.where('deleted_at', 'is', null)
		.selectAll()
		.executeTakeFirst();

	if (!album) {
		return json({ error: 'Album not found' }, { status: 404 });
	}

	const assets = await db
		.selectFrom('assets')
		.selectAll()
		.where('album_id', '=', id)
		.where('deleted_at', '=', null)
		.execute();

	return json({ ...album, assets });
};

export const DELETE: RequestHandler = async ({ params }) => {
	const { id } = params;

	if (!id) {
		return json({ error: 'Album id is required.' }, { status: 400 });
	}

	await deleteAlbum(id);

	return json({ success: true });
};
