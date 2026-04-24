import { sequence } from "@sveltejs/kit/hooks";
import { redirect, type Handle } from "@sveltejs/kit";
import { auth } from "$lib/server/auth";
import { svelteKitHandler } from "better-auth/svelte-kit";
import { error } from "node:console";
import { dev, building } from "$app/environment";
import { getMigrations } from "better-auth/db/migration";

const runMigrations = async () => {
    // Skip if we are in development mode OR currently building the app
    if (dev || building) return;

    try {
        const { runMigrations: execute } = await getMigrations(auth.options);
        await execute();
        console.log("Better Auth database migrations applied.");
    } catch (e) {
        console.error("Better Auth migration failed:", e);
        
        process.exit(1); 
    }
};

await runMigrations()

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
    const apiKey = event.request.headers.get("x-api-key")
    if (apiKey != null) {
        return handleSessionApiKey({ event, resolve })
    }

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

const handleSessionApiKey: Handle = async ({ event, resolve }) => {
    const apiKey = event.request.headers.get("x-api-key")
    if (!apiKey) {
        throw error("apikey is missing.")
    }

    const resp = await auth.api.verifyApiKey({
        body: {
            key: apiKey,
        },
    });

    if (resp.error) {
        throw resp.error.message
    }

    if (!resp.valid) {
        throw error("API key is invalid")
    }

    return resolve(event)
}

export const handle = sequence(handleBetterAuth, handleSession);
