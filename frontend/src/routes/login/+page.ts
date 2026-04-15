import { authClient } from "$lib/auth-client";
import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
    return {
        ...await authClient.signIn.oauth2({
            providerId: "generic-oauth",
            callbackURL: "/album",
            errorCallbackURL: "/error",
            // disableRedirect: true,
        })
    }
};
