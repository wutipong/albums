import { createChannel, createClient } from 'nice-grpc';
import {
    type WorkerServiceClient,
    WorkerServiceDefinition,
} from '$lib/server/grpc/proto/worker';
import { env } from '$env/dynamic/private';

const channel = createChannel(`http://${env.WORKER_ADDRESS ?? "localhost:7173"}`);

const client: WorkerServiceClient = createClient(
    WorkerServiceDefinition,
    channel,
);

export async function notifyProcessAsset(id: string){
    return await client.notifyProcessAsset({id: id})
}