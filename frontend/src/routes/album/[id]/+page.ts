import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params, fetch }) => {
    const { id } = params;
    
    const resp = await fetch(`/api/album/${id}`);
    const album = await resp.json();

    return { ...album };
};