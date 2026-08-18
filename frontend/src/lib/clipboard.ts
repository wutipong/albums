export async function copyImageToClipboard(imageUrl: string | URL) {
	// 1. Load the JPEG image into an Image object
	const img = new Image();
	img.crossOrigin = 'anonymous'; // Prevents canvas staining for external URLs
	img.src = imageUrl.toString();

	await new Promise((resolve, reject) => {
		img.onload = resolve;
		img.onerror = reject;
	});

	// 2. Draw the image onto an HTML Canvas
	const canvas = document.createElement('canvas');
	canvas.width = img.width;
	canvas.height = img.height;

	const ctx = canvas.getContext('2d');
	ctx?.drawImage(img, 0, 0);

	// 3. Export the canvas as a PNG blob
	canvas.toBlob(async (blob) => {
		if (!blob) {
			console.error('Canvas conversion failed');
			return;
		}

		try {
			// 4. Write the PNG blob to the clipboard
			const item = new ClipboardItem({ 'image/png': blob });
			await navigator.clipboard.write([item]);
			console.log('JPEG successfully converted and copied as PNG!');
		} catch (err) {
			console.error('Clipboard write failed:', err);
		}
	}, 'image/png'); // Forces the output to be a PNG
}
