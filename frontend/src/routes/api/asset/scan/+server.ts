import type { RequestHandler } from "./$types";
import { notifyScanCache } from "$lib/server/grpc/worker";

export const GET: RequestHandler = async ({ params }) => {
	await notifyScanCache()
}
    