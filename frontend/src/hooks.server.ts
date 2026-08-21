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
	if (event.url.pathname == '/login') {
		return resolve(event);
	}

	const session = event.locals.session;
	if (!session) {
		error(403, 'forbidden.');
	}

	return resolve(event);
};

const handleAdmin: Handle = async ({ event, resolve }) => {
	if (!event.url.pathname.startsWith('/admin')) {
		return resolve(event);
	}

	const user = event.locals.user;
	if (user.role !== 'admin') {
		error(403, 'forbidden.');
	}

	return resolve(event);
};

export const handle = sequence(handleBetterAuth, handleAdmin, handleSession);
