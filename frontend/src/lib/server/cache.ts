import path from "node:path";
import { env } from '$env/dynamic/private'

export function createCacheAssetPath(assetId: string, ...args: string[]): string {
    try {
        const topLevelDir = assetId.substring(0, 2);
        const secondLevelDir = assetId.substring(2, 4);

        return path.join(env.CACHE_DIR, topLevelDir, secondLevelDir, assetId, ...args);
    } catch (error) {
        throw new Error("Failed to create cache asset path");
    }
}