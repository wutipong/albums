import type { RequestHandler } from './$types';
import { updateAllAlbumThumbnail } from '$lib/server/grpc/worker';
import { error, json } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ url }) => {
	let missingOnly = url.searchParams.get('missingOnly') == 'true';
	try {
		await updateAllAlbumThumbnail(missingOnly);
	} catch {
		return error(422, 'unable to identify the worker service.');
	}
	return json({ success: true });
};
