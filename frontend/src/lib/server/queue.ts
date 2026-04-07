import { env } from '$env/dynamic/private';
import { Queue } from 'bullmq';
import { Redis } from 'ioredis';

const connection = new Redis(env.REDIS_URL);
export const processAssetQueue = new Queue('process-asset', { connection });