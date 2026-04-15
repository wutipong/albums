import { authClient } from "$lib/auth-client";
import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
   return {...await authClient.signOut()};
};
