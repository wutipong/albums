import { sequence } from "@sveltejs/kit/hooks";
import { redirect, type Handle } from "@sveltejs/kit";
import { auth } from "$lib/server/auth";
import { svelteKitHandler } from "better-auth/svelte-kit";
import { building } from '$app/environment';

const handleBetterAuth: Handle = async ({ event, resolve }) => {
    // path to your auth file
    const session = await auth.api.getSession({ headers: event.request.headers });

    if (session) { // Fetch current session from Better Auth 
        event.locals.session = session.session;
        event.locals.user = session.user;
    }

    return svelteKitHandler({ event, resolve, auth, building });
};

const handleSession: Handle = async ({ event, resolve }) => {
    const session = event.locals.session;

    if (event.url.pathname.startsWith('/login')) {
        return resolve(event)
    }

    if (session == null) {
        redirect(307, "/login")
    }

    if (Date.now() > session.expiresAt) {
        redirect(307, "/login")
    }

    return resolve(event);
}

export const handle = sequence(handleBetterAuth, handleSession);
