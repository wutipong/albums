import path from "node:path"

const imageExtensions = [
    ".jpg",
    ".jpeg",
    ".png",
    ".gif",
    ".svg",
    ".tiff",
    ".webp",
]

const animationExtensions = [
    ".gif",
    ".webp"
]

const videoExtensions = [
    ".mp4",
    ".webm"
]

export function isImage(p: string) {
    const ext = path.extname(p)

    return imageExtensions.includes(ext)
}

export function isAnimation(p: string) {
    const ext = path.extname(p)

    return animationExtensions.includes(ext)
}

export function isVideo(p: string) {
    const ext = path.extname(p)

    return videoExtensions.includes(ext)
}