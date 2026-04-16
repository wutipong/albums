import { authClient } from "$lib/auth-client";
import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
   await authClient.signOut();

   redirect(308, '/login')
};
