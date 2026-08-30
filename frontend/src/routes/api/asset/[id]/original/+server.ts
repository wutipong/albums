import { db } from '$lib/server/db';
import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { s3 } from '$lib/server/s3';
import mime from 'mime-types';
import path from 'node:path';

export const GET: RequestHandler = async ({ params }) => {
	const { id } = params;

	const asset = await db
		.selectFrom('assets')
		.selectAll()
		.where('id', '=', id)
		.where('assets.deleted_at', 'is', null)
		.executeTakeFirst();

	if (!asset) {
		return error(404, 'Asset not found.');
	}

	const file = s3.file(asset.original)
	const stat = await file.stat()
	const mimetype = mime.lookup(asset.filename) || 'application/octet-stream';

	return new Response(await file.arrayBuffer(), {
		headers: {
			'Content-Type': mimetype,
			'Content-Disposition': `attachment; filename="${path.basename(asset.filename)}"`,
			'Content-Length': stat.size.toString() || '',
			'Last-Modified': stat.lastModified.toISOString() || '',
			ETag: stat.etag || ''
		}
	});
};
