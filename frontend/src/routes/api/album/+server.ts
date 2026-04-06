import type { RequestHandler } from '@sveltejs/kit';
import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db/db';

export const POST: RequestHandler = async ({ request }) => {
  const body = await request.json();
  const { name } = body;

  if (!name || typeof name !== 'string') {
    return json({ error: 'Missing or invalid name' }, { status: 400 });
  }

  const album = await db.insertInto("albums")
    .values({ name })
    .returningAll()
    .executeTakeFirstOrThrow()

  return json(album, { status: 201 });
};

export const GET: RequestHandler = async () => {
  const albums = await db.selectFrom('albums')
    .selectAll()
    .where('albums.deleted_at', 'is', null)
    .execute()
    
  return json({ albums });
}