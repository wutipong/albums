import type { RequestHandler } from './$types';
import { updateAllImageEmbedding } from '$lib/server/grpc/worker';
import { error, json } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ locals }) => {
	try {
		await updateAllImageEmbedding();
	} catch {
		return json({ success: false });
	}
	return json({ success: true });
};
