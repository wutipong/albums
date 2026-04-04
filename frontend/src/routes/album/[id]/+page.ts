import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params, fetch }) => {
    const { id } = params;
    
    const resp = await fetch(`/api/album/${id}`);
    const { assets } = await resp.json();

    return { assets };
};