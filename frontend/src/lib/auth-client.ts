
import { createAuthClient } from "better-auth/svelte"
import { genericOAuthClient } from "better-auth/client/plugins"
import { apiKeyClient } from "@better-auth/api-key/client"

export const authClient = createAuthClient({
    plugins: [
        genericOAuthClient(),
        apiKeyClient(),
    ]
})