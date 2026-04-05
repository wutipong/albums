import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params, fetch }) => {
	const res = await fetch('/api/album');
	const data = await res.json();
	return { albums: data.albums };
};