import { generateImageUrl } from '@imgproxy/imgproxy-node';
import type { IPostgresInterval } from 'postgres-interval';
import { env } from '$env/dynamic/private';
import { s3 } from './s3';
import { Temporal } from 'temporal-polyfill';

export async function createResponseAssetList(
	assets: {
		id: string;
		album_id: string;
		created_at: Date;
		deleted_at: Date | null;
		filename: string;
		image_embedding: string | null;
		image_frames: number;
		modified_at: Date;
		original: string;
		preview: string;
		process_status: 'failed' | 'pending' | 'processed' | 'processing' | 'uploading';
		thumbnail: string;
		thumbnail_height: number;
		thumbnail_width: number;
		type: 'audio' | 'video' | 'animated' | 'image';
		video_duration: IPostgresInterval;
		view: string;
		view_height: number;
		view_width: number;
	}[]
) {
	const outAssets = [];
	for (const asset of assets) {
		const outAsset = {
			...asset,
			video_duration: 0,
			thumbnail_url: '',
			preview_url: '',
			view_url: '',
			original_url: ''
		};
		const video_duration = Temporal.Duration.from(asset.video_duration.toISOString());

		const bypass = asset.image_frames > 1 || asset.type == 'video';

		const thumbnail_url = asset.thumbnail === ''? '' : generateImageUrl({
			endpoint: env.IMGPROXY_URL,
			url: `s3://${env.S3_BUCKET}/${asset.thumbnail}`,
			options: {
				raw: bypass,
				resizing_type: 'auto',
				height: 200,
				enlarge: 1
			},
			salt: env.IMGPROXY_SALT,
			key: env.IMGPROXY_KEY
		});

		const preview_url = asset.preview === ''? '' : generateImageUrl({
			endpoint: env.IMGPROXY_URL,
			url: `s3://${env.S3_BUCKET}/${asset.preview}`,
			options: {
				raw: bypass,
				resizing_type: 'auto',
				height: 200,
				enlarge: 1
			},
			salt: env.IMGPROXY_SALT,
			key: env.IMGPROXY_KEY
		});

		let view_url = '';
		switch (asset.type) {
			case 'image':
				view_url = asset.view === ''? '' : generateImageUrl({
					endpoint: env.IMGPROXY_URL,
					url: `s3://${env.S3_BUCKET}/${asset.view}`,
					options: {
						raw: bypass,
						resizing_type: 'auto',
						height: 2000,
						enlarge: 1
					},
					salt: env.IMGPROXY_SALT,
					key: env.IMGPROXY_KEY
				});
				break;

			case 'video':
				view_url = asset.view === ''? '' : generateImageUrl({
					endpoint: env.IMGPROXY_URL,
					url: `s3://${env.S3_BUCKET}/${asset.view}`,
					options: {
						raw: true
					},
					salt: env.IMGPROXY_SALT,
					key: env.IMGPROXY_KEY
				});
				break;
		}

		const copy_url = asset.type === 'video' ? '' : `/api/asset/${asset.id}/original/`;

		const original_url = s3.presign(asset.original);
		const out = {
			...asset,
			video_duration,
			thumbnail_url,
			preview_url,
			view_url,
			original_url,
			copy_url
		};

		outAssets.push(out);
	}
	return outAssets;
}
