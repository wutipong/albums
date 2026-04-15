import { env } from "$env/dynamic/private";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
    return {
        providerId: env.OIDC_PROVIDER_ID ?? "",
        callbackURL: '/album',
        errorCallbackURL: '/error'
        // disableRedirect: true,
    }
};
