import type { RequestHandler } from '@sveltejs/kit';
import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { s3 } from '$lib/server/s3';
import { DeleteObjectsCommand } from '@aws-sdk/client-s3';

export const DELETE: RequestHandler = async ({ params }) => {
  const { id } = params;

  if (!id) {
    return json({ error: 'Album id is required.' }, { status: 400 });
  }

  const assets = await db.selectFrom('assets')
    .selectAll()
    .where('album_id', '=', id)
    .where('deleted_at', '=', null)
    .execute()

  const uniqueKeys = new Set();
  for (const asset of assets) {
    // Collect all potential keys into an array
    const keys = [asset.original, asset.preview, asset.thumbnail, asset.view];

    // Iterate over the keys and add any non-empty string to the Set
    for (const key of keys) {
      if (key && key !== '') {
        uniqueKeys.add(key);
      }
    }
  }

  // Convert the Set of unique keys into the desired array format
  const deletingObjs = Array.from(uniqueKeys).map(key => ({ Key: key as string }));

  try {
    await s3.send(new DeleteObjectsCommand({
      Bucket: "",
      Delete: {
        Objects: deletingObjs,
      }
    }))
  } catch (err) {
    return json({ error: 'Unable to delete album assets.' }, { status: 400 });
  }

  const deleted = await db.updateTable('assets')
    .set({ deleted_at: new Date() })
    .where('album_id', '=', id)
    .where('deleted_at', 'is', null)
    .returningAll()
    .execute()

  if (!deleted) {
    return json({ error: 'Unable to set deleted date on assets.' }, { status: 404 });
  }

  const album = await db.updateTable('albums')
    .set({ deleted_at: new Date() })
    .where('id', '=', id)
    .where('deleted_at', 'is', null)
    .returningAll()
    .executeTakeFirst();

  if (!album) {
    return json({ error: 'Album not found' }, { status: 404 });
  }

  return json({ success: true });
};
