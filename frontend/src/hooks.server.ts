import { sequence } from "@sveltejs/kit/hooks";
import type { Handle } from "@sveltejs/kit";
import { auth } from "$lib/server/auth";
import { svelteKitHandler } from "better-auth/svelte-kit";
import { building } from '$app/environment';

const handleBetterAuth: Handle = async ({ event, resolve }) => {
    // path to your auth file
    const session = await auth.api.getSession({ headers: event.request.headers });

    if (session){ // Fetch current session from Better Auth 
        event.locals.session = session.session;
        event.locals.user = session.user;
    }

    return svelteKitHandler({ event, resolve, auth, building });
};

export const handle = sequence(handleBetterAuth);
