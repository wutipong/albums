import type { RequestHandler } from "./$types";
import { notifyScanCache } from "$lib/server/grpc/worker";
import { json } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ }) => {
	try{
		await notifyScanCache()
	} catch {
		return json({ success: false})
	}
	return json({ success: true });
}
    