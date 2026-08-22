import { db } from '$lib/server/db';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const count = await db
		.selectFrom('assets')
		.select((eb) => [eb.fn.countAll().as('count')])
		.executeTakeFirst();

	const pendingCount = await db
		.selectFrom('assets')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('process_status', '=', 'pending')
		.executeTakeFirst();

	const uploadingCount = await db
		.selectFrom('assets')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('process_status', '=', 'uploading')
		.executeTakeFirst();

	const failedCount = await db
		.selectFrom('assets')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('process_status', '=', 'failed')
		.executeTakeFirst();

	const imageCount = await db
		.selectFrom('assets')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('type', '=', 'image')
		.executeTakeFirst();

	const embeddingCount = await db
		.selectFrom('assets')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('image_embedding', 'is not', null)
		.executeTakeFirst();

	const videoCount = await db
		.selectFrom('assets')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('type', '=', 'video')
		.executeTakeFirst();

	return {
		total: count ? BigInt(count.count) : 0n,
		pending: pendingCount ? BigInt(pendingCount.count) : 0n,
		uploading: uploadingCount ? BigInt(uploadingCount.count) : 0n,
		failed: failedCount ? BigInt(failedCount.count) : 0n,

		images: imageCount ? BigInt(imageCount.count) : 0n,
		embeddings: embeddingCount ? BigInt(embeddingCount.count) : 0n,
		video: videoCount ? BigInt(videoCount.count) : 0n
	};
};
