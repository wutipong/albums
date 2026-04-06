import path from "node:path";

export function createCacheAssetPath(assetId: string, ...args: string[]): string {
    try{
        const topLevelDir = assetId.substring(0, 2);
        const secondLevelDir = assetId.substring(2, 4);

        return path.join("cache", topLevelDir, secondLevelDir, assetId, ...args);
    } catch (error) {
        throw new Error("Failed to create cache asset path");
    }
}