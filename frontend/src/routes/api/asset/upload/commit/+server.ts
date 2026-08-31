import type { RequestHandler } from './$types';
import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import * as mime from 'mime-types';
import { notifyProcessAsset } from '$lib/server/grpc/worker';

export const POST: RequestHandler = async ({ request }) => {
	const req = await request.json();
	const id = req.id;
	const success = req.success;

	let asset = await db
		.selectFrom('assets')
		.selectAll()
		.where('id', '=', id)
		.where('deleted_at', 'is', null)
		.executeTakeFirst();

	if (!asset) {
		return json({ status: 'asset record not found.' }, { status: 404 });
	}

	if (success) {
		asset.process_status = 'pending';

		const mimetype = mime.lookup(asset.filename);
		if (!mimetype) {
			return json({ status: 'invalid content type' }, { status: 400 });
		}

		if (mimetype.startsWith('image/')) {
			asset.type = 'image';
		} else if (mimetype.startsWith('video')) {
			asset.type = 'video';
		}
	} else {
		asset.process_status = 'failed';
		asset.original = '';
	}
	const resp = await db
		.updateTable('assets')
		.set(asset)
		.where('id', '=', asset.id)
		.executeTakeFirst();

	if (!resp) {
		return json({ status: 'Failed to update asset' }, { status: 503 });
	}

	if (success) {
		try {
			await notifyProcessAsset(asset.id);
		} catch (error) {
			return json(
				{ asset, status: 'asset is commited, but it is not queued to processing.' },
				{ status: 200 }
			);
		}
		return json({ asset: asset, status: 'asset is accepted' }, { status: 201 });
	} else {
		return json(
			{ status: 'asset status is updated, but the upload failed.' },
			{ status: 200 }
		);
	}
};
