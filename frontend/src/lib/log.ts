import { dev } from '$app/environment';
import pino from 'pino';

console.log(`Initializing logger. Environment: ${dev ? 'development' : 'production'}`);

const logger = pino({
    level: dev ? 'debug' : 'info',
    browser: {
        asObject: true
    },
    transport: dev
        ? {
                target: 'pino-pretty',
                options: {
                    colorize: true,
                    translateTime: 'SYS:standard',
                    ignore: 'pid,hostname'
                }
            }
        : undefined
});

export default logger;
