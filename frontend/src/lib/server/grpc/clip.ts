import { createChannel, createClient } from 'nice-grpc';
import {
    type EncodingServiceClient,
    EncodingServiceDefinition,
} from '$lib/server/grpc/proto/clip';
import { env } from '$env/dynamic/private';


const channel = createChannel(`http://${env.CLIP_ADDRESS ?? "localhost:8173"}`);
const client: EncodingServiceClient = createClient(
    EncodingServiceDefinition,
    channel,
);

export async function encodeText(input: string){
    return await client.encodeText({input: input})
}
