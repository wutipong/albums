import { env } from '$env/dynamic/private';
import { db } from '$lib/server/db';
import type { PageServerLoad } from './$types';
import { generateImageUrl } from '@imgproxy/imgproxy-node';

export const load: PageServerLoad = async () => {
	let aspect: 'portrait' | 'landscape' = env.COVER_ASPECT == 'portrait' ? 'portrait' : 'landscape';

	const albums = await db
		.selectFrom('albums')
		.selectAll()
		.where('deleted_at', 'is', null)
		.orderBy('name')
		.execute();

	const outAlbums = [];
	for (const album of albums) {
		let cover_url = '';
		if (album.cover != '') {
			cover_url = generateImageUrl({
				endpoint: env.IMGPROXY_URL,
				url: `s3://${env.S3_BUCKET}/${album.cover}`,
				options: {
					resizing_type: 'fill',
					height: aspect == 'landscape' ? 200 : 450,
					width: 300,
					enlarge: 1
				},
				salt: env.IMGPROXY_SALT,
				key: env.IMGPROXY_KEY
			});
		}

		outAlbums.push({ ...album, cover_url });
	}

	return { aspect, albums: outAlbums };
};
