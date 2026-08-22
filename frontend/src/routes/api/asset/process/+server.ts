import type { RequestHandler } from './$types';
import { notifyProcessAllAssets, updateAllImageEmbedding } from '$lib/server/grpc/worker';
import { error, json } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ url }) => {
	let missingOnly = url.searchParams.get("missingOnly") === "true"
	try {
		await notifyProcessAllAssets(missingOnly);
	} catch {
		return json({ success: false });
	}
	return json({ success: true });
};
