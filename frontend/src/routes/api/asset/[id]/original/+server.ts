import { db } from '$lib/server/db';
import { error, json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { GetObjectCommand } from '@aws-sdk/client-s3';
import { env } from '$env/dynamic/private';
import { s3 } from '$lib/server/s3';
import mime from 'mime-types';

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

	const command = new GetObjectCommand({
		Bucket: env.S3_BUCKET,
		Key: asset.original
	});

	const response = await s3.send(command);
	if (!response.Body) {
		return error(404, 'Asset not found in S3.');
	}

	const byteArray = await response.Body.transformToByteArray();
	const mimetype = mime.lookup(asset.filename) || 'application/octet-stream';

	return new Response(byteArray as BodyInit, {
		headers: {
			'Content-Type': mimetype,
			'Content-Length': response.ContentLength?.toString() || '',
			'Last-Modified': response.LastModified?.toUTCString() || '',
			ETag: response.ETag || ''
		}
	});
};
