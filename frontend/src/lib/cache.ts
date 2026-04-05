export function createCacheAssetPath(assetId: string, filename: string): string {
    try{
        return `./cache/${assetId.substring(0, 2)}/${assetId.substring(2, 4)}/${assetId}/${filename}`;
    } catch (error) {
        throw new Error("Failed to create cache asset path");
    }
}