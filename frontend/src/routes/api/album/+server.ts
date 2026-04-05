import type { RequestHandler } from '@sveltejs/kit';
import { json } from '@sveltejs/kit';
import { createAlbum, listAlbums } from '$lib/db/query_sql';
import { getDb } from '$lib/db/db';

export const POST: RequestHandler = async ({ request }) => {
  const body = await request.json();
  const { name } = body;

  if (!name || typeof name !== 'string') {
    return json({ error: 'Missing or invalid name' }, { status: 400 });
  }

  const album = await createAlbum(getDb(), { name });
  return json(album, { status: 201 });
};

export const GET: RequestHandler = async () => {
  const albums = await listAlbums(getDb());
  return json({ albums });
}