import { db } from '$lib/server/db';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params, locals }) => {
	const session = locals.session;
	const user = locals.user;

	return json({
		user: user,
		session: session
	});
};
