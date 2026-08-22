import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { s3, s3Public } from '$lib/server/s3';
import { getSignedUrl } from '@aws-sdk/s3-request-presigner';
import { PutObjectCommand } from '@aws-sdk/client-s3';
import * as mime from 'mime-types';

import { randomUUID } from 'node:crypto';
import { env } from '$env/dynamic/private';
import path from 'node:path';

export const POST: RequestHandler = async ({ request }) => {
	const req = await request.json();
	const albumId = req.album_id;
	const filename = req.filename;
	const checksum = req.checksum;
	const network = req.network ?? 'private';

	const contentType = mime.contentType(path.basename(filename));

	if (!contentType) {
		return json({ success: false, error: 'Failed to recognize filetype' }, { status: 400 });
	}

	const extension = mime.extension(mime.lookup(filename) || '');
	const key = `pending/${randomUUID()}.${extension}`;

	const type = contentType.substring(0, contentType.indexOf('/'));

	if (type != 'image' && type != 'video') {
		return json({ success: false, error: 'Unsupported asset type.' }, { status: 400 });
	}

	const album = await db
		.selectFrom('albums')
		.selectAll()
		.where('albums.id', '=', albumId)
		.where('albums.deleted_at', 'is', null)
		.limit(1)
		.executeTakeFirst();

	if (!album) {
		return json({ success: false, error: 'Album not found' }, { status: 404 });
	}

	const existed = await db
		.selectFrom('assets')
		.where('album_id', '=', albumId)
		.where('filename', '=', filename)
		.selectAll()
		.limit(1)
		.executeTakeFirst();
	if (existed && existed.process_status != 'uploading' && existed.process_status != 'failed') {
		return json({ success: false, error: 'duplicate asset' }, { status: 409 });
	}

	let asset = null;
	if (existed) {
		asset = existed;
		asset.original = key;
		asset.type = type;
		asset.process_status = 'uploading';
	} else {
		asset = await db
			.insertInto('assets')
			.values({
				album_id: albumId,
				filename: filename,
				type: type,
				process_status: 'uploading',
				original: key
			})
			.returningAll()
			.executeTakeFirstOrThrow();
	}

	if (!asset) {
		return json({ success: false, error: 'Failed to create asset' }, { status: 500 });
	}

	const command = new PutObjectCommand({
		Bucket: env.S3_BUCKET,
		Key: asset.original,
		ChecksumCRC32: checksum,
		ContentType: contentType
	});

	const url = await getSignedUrl(network === 'public' ? s3Public : s3, command, {
		expiresIn: 3600
	});

	return json({ id: asset.id, url, success: true });
};
