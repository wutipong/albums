import { S3Client } from '@aws-sdk/client-s3';
import { env } from '$env/dynamic/private';

export const s3 = new S3Client({
	region: env.AWS_DEFAULT_REGION || 'dummy-region',
	endpoint: env.AWS_ENDPOINT_URL || 'http://localhost',
	credentials: {
		accessKeyId: env.AWS_ACCESS_KEY_ID || 'dummy-key',
		secretAccessKey: env.AWS_SECRET_ACCESS_KEY || 'dummy-secret'
	},
	forcePathStyle: true
});
