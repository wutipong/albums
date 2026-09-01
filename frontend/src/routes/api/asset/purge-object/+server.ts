import type { RequestHandler } from './$types';
import { notifyPurgeObject } from '$lib/server/grpc/worker';
import { json } from '@sveltejs/kit';

export const GET: RequestHandler = async ({}) => {
	try {
		await notifyPurgeObject();
	} catch {
		return json({ success: false });
	}
	return json({ success: true });
};
