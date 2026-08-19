import { getRequestEvent } from '$app/server';
import { env } from '$env/dynamic/private';
import { betterAuth } from 'better-auth';
import { genericOAuth } from 'better-auth/plugins';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { apiKey } from '@better-auth/api-key';
import { Pool } from 'pg';

export const auth = betterAuth({
	database: new Pool({
		connectionString: env.DATABASE_URL
	}),
	plugins: [
		sveltekitCookies(getRequestEvent),
		apiKey({
			rateLimit: { enabled: false }
		}),
		genericOAuth({
			config: [
				{
					providerId: env.OIDC_PROVIDER_ID || 'placeholder-provider',
					clientId: env.OIDC_CLIENT_ID || 'placeholder-client',
					clientSecret: env.OIDC_SECRET || 'placeholder-secret',
					issuer: env.OIDC_ISSUER || 'https://placeholder-issuer.com',
					tokenUrl: env.OIDC_TOKEN || 'https://placeholder-issuer.com',
					authorizationUrl: env.OIDC_AUTHORIZE || 'https://placeholder-issuer.com',
					requireIssuerValidation: false,
					scopes: ['openid', 'email', 'profile']
				}
				// Add more providers as needed
			]
		})
	]
});
