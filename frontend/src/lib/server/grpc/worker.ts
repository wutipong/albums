import { createChannel, createClient } from 'nice-grpc';
import { type WorkerServiceClient, WorkerServiceDefinition } from '$lib/server/grpc/proto/worker';
import { env } from '$env/dynamic/private';

const channel = createChannel(`http://${env.WORKER_ADDRESS ?? 'localhost:7173'}`);

const client: WorkerServiceClient = createClient(WorkerServiceDefinition, channel);

export async function notifyProcessAsset(id: string) {
	return await client.notifyProcessAsset({ id: id });
}

export async function notifyScanCache() {
	return await client.notifyScanCache({});
}

export async function notifyUpdateAlbumCover(albumId: string, assetId?: string) {
	return await client.updateAlbumThumbnail({
		id: albumId,
		assetId: assetId
	});
}

export async function updateAllImageEmbedding(onlyMissing: boolean = true) {
	return await client.updateAllImageEmbedding({
		onlyMissing: onlyMissing
	});
}

export async function updateAllAlbumThumbnail(onlyMissing: boolean = true) {
	return await client.updateAllAlbumThumbnail({
		onlyMissing: onlyMissing
	});
}

export async function notifyProcessAllAssets(onlyMissing: boolean = true) {
	return await client.notifyProcessAllAssets({ onlyMissing: onlyMissing });
}

export async function notifyPurgeObject(){
	return await client.purgeUnusedObject({})
}
