import { env } from "$env/dynamic/private";
import { betterAuth } from "better-auth";
import { genericOAuth } from "better-auth/plugins"
import path from "node:path";

export const auth = betterAuth({
    plugins: [
        genericOAuth({
            config: [
                {
                    providerId: "generic-oauth",
                    clientId: env.OIDC_CLIENT_ID ?? "",
                    clientSecret: env.OIDC_SECRET,
                    discoveryUrl: env.OIDC_DISCOVERY_URL ?? "",
                    redirectURI: new URL("login/callback", env.BETTER_AUTH_URL).toString(),
                    // ... other config options
                },
                // Add more providers as needed
            ]
        })
    ]
})