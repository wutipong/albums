import { sequence } from '@sveltejs/kit/hooks';
import { error, redirect, type Handle, type ServerInit } from '@sveltejs/kit';
import { auth } from '$lib/server/auth';
import { svelteKitHandler } from 'better-auth/svelte-kit';
import { dev, building } from '$app/environment';
import { getMigrations } from 'better-auth/db/migration';

export const init: ServerInit = async () => {
	if (!dev && !building) {
		try {
			const { runMigrations: execute } = await getMigrations(auth.options);
			await execute();
			console.log('Better Auth database migrations applied.');
		} catch (e) {
			console.error('Better Auth migration failed:', e);
			process.exit(1);
		}
	}
};

const handleBetterAuth: Handle = async ({ event, resolve }) => {
	const session = await auth.api.getSession({ headers: event.request.headers });

	if (session) {
		// Fetch current session from Better Auth
		event.locals.session = session.session;
		event.locals.user = session.user;
	}

	return svelteKitHandler({ event, resolve, auth, building });
};

const handleSession: Handle = async ({ event, resolve }) => {
	const apiKey = event.request.headers.get('x-api-key');
	if (apiKey != null) {
		return handleSessionApiKey({ event, resolve });
	}

	if (event.url.pathname == '/login') {
		return resolve(event);
	}

	const session = event.locals.session;

	if (session == null) {
		redirect(307, '/');
	}

	if (Date.now() > session.expiresAt) {
		redirect(307, '/');
	}

	return resolve(event);
};

const handleSessionApiKey: Handle = async ({ event, resolve }) => {
	const apiKey = event.request.headers.get('x-api-key');
	if (!apiKey) {
		throw error(401, 'apikey is missing.');
	}

	const resp = await auth.api.verifyApiKey({
		body: {
			key: apiKey
		}
	});

	if (resp.error) {
		error(401, resp.error.message);
	}

	if (!resp.valid) {
		error(401, 'API key is invalid');
	}

	if (event.url.pathname !== '/api' && !event.url.pathname.startsWith('/api/')) {
		error(401, 'non-API access prohibited.');
	}

	return resolve(event);
};

const handleAdminSession: Handle = async ({ event, resolve }) => {
	if (!event.url.pathname.startsWith('/admin')) {
		return resolve(event);
	}

	const session = await auth.api.getSession({ headers: event.request.headers });
	if (session?.user.role !== 'admin') {
		error(403, 'forbidden.');
	}

	return resolve(event);
};

export const handle = sequence(handleBetterAuth, handleAdminSession, handleSession);
