import { getRequestEvent } from "$app/server";
import { env } from "$env/dynamic/private";
import { betterAuth } from "better-auth";
import { genericOAuth } from "better-auth/plugins"
import { sveltekitCookies } from "better-auth/svelte-kit";
import { apiKey } from "@better-auth/api-key"
import { Pool } from "pg";

export const auth = betterAuth({
    database: new Pool({
        connectionString: env.DATABASE_URL,
    }),
    plugins: [
        sveltekitCookies(getRequestEvent),
        apiKey({
            rateLimit: { enabled: false },
        }),
        genericOAuth({
            config: [
                {
                    providerId: env.OIDC_PROVIDER_ID ?? "",
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