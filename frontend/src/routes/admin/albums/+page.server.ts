import { db } from '$lib/server/db';
import { sql } from 'kysely';
import type { PageServerLoad } from '../asset/$types';

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

	const failedCount = await db
		.selectFrom('assets')
		.select((eb) => [eb.fn.countAll().as('count')])
		.where('process_status', '=', 'failed')
		.executeTakeFirst();

	const album_count = await db
		.selectFrom('albums')
		.select((eb) => [eb.fn.countAll().as('count')])
		.executeTakeFirst();

	const queue = await db
		.selectFrom('neoq_jobs')
		.where('status', '=', 'new')
		.where(sql<string>`payload->>'command'`, '=', 'process-asset')
		.orderBy('id', 'asc')
		.limit(20)
		.selectAll()
		.execute();

	const queue_total = await db
		.selectFrom('neoq_jobs')
		.select((eb) => eb.fn.countAll().as('count'))
		.executeTakeFirst();

	const queue_processed = await db
		.selectFrom('neoq_jobs')
		.select((eb) => eb.fn.countAll().as('count'))
		.where('status', '<>', 'new')
		.executeTakeFirst();

	const queue_pending = await db
		.selectFrom('neoq_jobs')
		.select((eb) => eb.fn.countAll().as('count'))
		.where('status', '=', 'new')
		.executeTakeFirst();

	const queue_failed = await db
		.selectFrom('neoq_jobs')
		.select((eb) => eb.fn.countAll().as('count'))
		.where('status', '=', 'failed')
		.executeTakeFirst();

	let queueItems = [];
	for (const job of queue) {
		if (!job.payload) {
			continue;
		}

		const payload = job.payload as any;
		const asset_id = payload.id;
		const asset = await db
			.selectFrom('assets')
			.leftJoin('albums', 'assets.album_id', 'albums.id')
			.select('assets.id as asset_id')
			.select('type')
			.select('albums.name as album_name')
			.select('album_id')
			.select('assets.created_at as created_at')
			.where('assets.id', '=', asset_id)
			.executeTakeFirst();
		if (asset) {
			queueItems.push({ id: job.id, ...asset });
		}
	}
	return {
		count: count ? BigInt(count.count) : 0n,
		queueItems: queueItems,
		queueProcessed: queue_processed ? BigInt(queue_processed.count) : 0n,
		queuePending: queue_pending ? BigInt(queue_pending.count) : 0n,
		queueFailed: queue_failed ? BigInt(queue_failed.count) : 0n,
		queueTotal: queue_total ? BigInt(queue_total.count) : 0n,
		failedCount: failedCount ? BigInt(failedCount.count) : 0n,
		pendingCount: pendingCount ? BigInt(pendingCount.count) : 0n,
		albumCount: album_count ? BigInt(album_count.count) : 0n
	};
};
