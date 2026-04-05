import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params, fetch }) => {
    const { id } = params;
    
    const album = await getAlbumInfo(fetch, id);

    return { ...album };
};

async function getAlbumInfo(fetch: { (input: RequestInfo | URL, init?: RequestInit): Promise<Response>; (input: string | URL | Request, init?: RequestInit): Promise<Response>; }, id: string) {
    const resp = await fetch(`/api/album/${id}`);
    const album = await resp.json();
    return album;
}
