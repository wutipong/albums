import { getRequestEvent } from "$app/server";
import { env } from "$env/dynamic/private";
import { betterAuth } from "better-auth";
import { genericOAuth } from "better-auth/plugins"
import { sveltekitCookies } from "better-auth/svelte-kit";
import { Pool } from "pg";

export const auth = betterAuth({
    database: new Pool({
        connectionString: env.DATABASE_URL,
    }),
    plugins: [
        sveltekitCookies(getRequestEvent),
        genericOAuth({
            config: [
                {
                    providerId: "generic-oauth",
                    clientId: env.OIDC_CLIENT_ID ?? "",
                    clientSecret: env.OIDC_SECRET,
                    discoveryUrl: env.OIDC_DISCOVERY_URL ?? "",
                    // ... other config options
                },
                // Add more providers as needed
            ]
        })
    ]
})